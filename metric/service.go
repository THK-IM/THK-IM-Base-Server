package metric

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

var defaultMetricPath = "/metrics"

var reqCnt = &Metric{
	ID:          "reqCnt",
	Name:        "requests_total",
	Description: "How many HTTP requests processed, partitioned by status code and HTTP method.",
	Type:        "counter_vec",
	Args:        []string{"code", "method", "handler", "host", "url"},
}

var reqDur = &Metric{
	ID:          "reqDur",
	Name:        "request_duration_seconds",
	Description: "The HTTP request latencies in seconds.",
	Type:        "histogram_vec",
	Args:        []string{"code", "method", "url"},
}

var resSz = &Metric{
	ID:          "resSz",
	Name:        "response_size_bytes",
	Description: "The HTTP response sizes in bytes.",
	Type:        "summary",
}

var reqSz = &Metric{
	ID:          "reqSz",
	Name:        "request_size_bytes",
	Description: "The HTTP request sizes in bytes.",
	Type:        "summary",
}

var httpMetrics = []*Metric{
	reqCnt,
	reqDur,
	resSz,
	reqSz,
}

type RequestCounterURLLabelMappingFn func(c *gin.Context) string

type Service struct {
	NodeId                  int64
	ServerName              string
	MetricsPath             string
	MetricsList             []*Metric
	logger                  *logrus.Entry
	reqCnt                  *prometheus.CounterVec
	reqDur                  *prometheus.HistogramVec
	reqSz, resSz            prometheus.Summary
	PushGateway             PrometheusPushGateway
	ReqCntURLLabelMappingFn RequestCounterURLLabelMappingFn
}

type PrometheusPushGateway struct {
	PushIntervalSeconds time.Duration
	PushGatewayURL      string
	MetricsURL          string
	Job                 string
}

func NewService(serverName string, nodeId int64, logger *logrus.Entry) *Service {
	s := &Service{
		NodeId:      nodeId,
		ServerName:  serverName,
		MetricsPath: defaultMetricPath,
		ReqCntURLLabelMappingFn: func(c *gin.Context) string {
			url := c.Request.URL.Path
			for _, param := range c.Params {
				url = strings.Replace(url, param.Value, ":"+param.Key, 1)
			}
			return url
		},
		logger: logger,
	}
	return s
}

func (s *Service) InitMetrics(extraMetrics ...*Metric) {
	var metricsList []*Metric
	for _, metric := range extraMetrics {
		metricsList = append(metricsList, metric)
	}
	for _, metric := range httpMetrics {
		metricsList = append(metricsList, metric)
	}
	s.MetricsList = metricsList
	s.registerMetrics(s.ServerName)
}

func (s *Service) SetPushGateway(jobName, pushGatewayURL, metricsURL string, pushInterval time.Duration) {
	s.PushGateway.Job = jobName
	s.PushGateway.PushGatewayURL = pushGatewayURL
	s.PushGateway.MetricsURL = metricsURL
	s.PushGateway.PushIntervalSeconds = pushInterval
	s.startPushTicker()
}

func (s *Service) getMetrics() ([]byte, error) {
	response, err := http.Get(s.PushGateway.MetricsURL)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = response.Body.Close()
	}()
	return io.ReadAll(response.Body)
}

func (s *Service) getPushGatewayURL() string {
	if s.PushGateway.Job == "" {
		s.PushGateway.Job = "gin"
	}
	h := fmt.Sprintf("node-%d", s.NodeId)
	return s.PushGateway.PushGatewayURL + "/metrics/job/" + s.PushGateway.Job + "/instance/" + h
}

func (s *Service) sendMetricsToPushGateway(metrics []byte) {
	if metrics == nil {
		s.logger.Info("metric bytes is nil ")
		return
	}
	req, err := http.NewRequest("POST", s.getPushGatewayURL(), bytes.NewBuffer(metrics))
	client := &http.Client{}
	if _, err = client.Do(req); err != nil {
		s.logger.Error("Error sending to push gateway, ", err)
	}
}

func (s *Service) startPushTicker() {
	ticker := time.NewTicker(s.PushGateway.PushIntervalSeconds)
	go func() {
		for range ticker.C {
			m, err := s.getMetrics()
			if err != nil {
				fmt.Println(err)
			} else {
				s.sendMetricsToPushGateway(m)
			}
		}
	}()
}

// NewMetric associates prometheus.Collector based on Metric.Type
func (s *Service) newMetric(m *Metric, serverName string) prometheus.Collector {
	var metric prometheus.Collector
	switch m.Type {
	case "counter_vec":
		metric = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Subsystem: serverName,
				Name:      m.Name,
				Help:      m.Description,
			},
			m.Args,
		)
	case "counter":
		metric = prometheus.NewCounter(
			prometheus.CounterOpts{
				Subsystem: serverName,
				Name:      m.Name,
				Help:      m.Description,
			},
		)
	case "gauge_vec":
		metric = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Subsystem: serverName,
				Name:      m.Name,
				Help:      m.Description,
			},
			m.Args,
		)
	case "gauge":
		metric = prometheus.NewGauge(
			prometheus.GaugeOpts{
				Subsystem: serverName,
				Name:      m.Name,
				Help:      m.Description,
			},
		)
	case "histogram_vec":
		metric = prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Subsystem: serverName,
				Name:      m.Name,
				Help:      m.Description,
			},
			m.Args,
		)
	case "histogram":
		metric = prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Subsystem: serverName,
				Name:      m.Name,
				Help:      m.Description,
			},
		)
	case "summary_vec":
		metric = prometheus.NewSummaryVec(
			prometheus.SummaryOpts{
				Subsystem: serverName,
				Name:      m.Name,
				Help:      m.Description,
			},
			m.Args,
		)
	case "summary":
		metric = prometheus.NewSummary(
			prometheus.SummaryOpts{
				Subsystem: serverName,
				Name:      m.Name,
				Help:      m.Description,
			},
		)
	}
	return metric
}

func (s *Service) registerMetrics(serverName string) {
	for _, metricDef := range s.MetricsList {
		metric := s.newMetric(metricDef, serverName)
		if err := prometheus.Register(metric); err != nil {
			s.logger.Errorf("%s could not be registered in MetricServer, error: %v", metricDef.Name, err)
		}
		switch metricDef {
		case reqCnt:
			s.reqCnt = metric.(*prometheus.CounterVec)
		case reqDur:
			s.reqDur = metric.(*prometheus.HistogramVec)
		case resSz:
			s.resSz = metric.(prometheus.Summary)
		case reqSz:
			s.reqSz = metric.(prometheus.Summary)
		}
		metricDef.MetricCollector = metric
	}
}

func (s *Service) Use(g *gin.Engine) {
	g.Use(s.handlerFunc()) // middleware 收集http指标

	// http获取指标路由接口
	h := promhttp.Handler()
	g.GET(s.MetricsPath, func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	})

}

func (s *Service) handlerFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == s.MetricsPath {
			c.Next()
			return
		}

		start := time.Now()
		reqSize := computeApproximateRequestSize(c.Request)

		c.Next()

		status := strconv.Itoa(c.Writer.Status())
		elapsed := float64(time.Since(start)) / float64(time.Second)
		resSize := float64(c.Writer.Size())

		url := s.ReqCntURLLabelMappingFn(c)

		s.reqDur.WithLabelValues(status, c.Request.Method, url).Observe(elapsed)
		s.reqCnt.WithLabelValues(status, c.Request.Method, c.HandlerName(), c.Request.Host, url).Inc()
		s.reqSz.Observe(float64(reqSize))
		s.resSz.Observe(resSize)
	}
}

func computeApproximateRequestSize(r *http.Request) int {
	s := 0
	if r.URL != nil {
		s = len(r.URL.Path)
	}

	s += len(r.Method)
	s += len(r.Proto)
	for name, values := range r.Header {
		s += len(name)
		for _, value := range values {
			s += len(value)
		}
	}
	s += len(r.Host)

	// N.B. r.Form and r.MultipartForm are assumed to be included in r.URL.

	if r.ContentLength != -1 {
		s += int(r.ContentLength)
	}
	return s
}
