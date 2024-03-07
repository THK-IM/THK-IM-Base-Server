package httptest

import (
	"fmt"
	"math"
	"net/http"
	"time"
)

var defaultHttpClient = http.Client{
	Transport: &http.Transport{
		MaxIdleConns:        20,
		MaxIdleConnsPerHost: 60,
		IdleConnTimeout:     20 * time.Second,
	},
	Timeout: time.Second,
}

type (
	Result struct {
		serialNo   int   // 请求流水号
		statusCode int   // 状态码
		bodySize   int64 // 响应body大小
		duration   int64 // 耗时 单位ms
		err        error // 错误
	}
	Task struct {
		count         int // 总测试次数
		concurrent    int // 测试并发数
		startTime     time.Time
		endTime       time.Time
		results       []*Result
		resultChan    chan *Result
		taskIndexChan chan int
		requestFunc   RequestFunc
		resultFunc    ResultFunc
	}
	RequestFunc func(index, channelIndex int, client http.Client) *Result
	ResultFunc  func(task *Task)
)

func (h Result) SerialNo() int {
	return h.serialNo
}

func (h Result) StatusCode() int {
	return h.statusCode
}

func (h Result) BodySize() int64 {
	return h.bodySize
}

func (h Result) Duration() int64 {
	return h.duration
}

func NewHttpTestResult(serialNo, statusCode int, bodySize, duration int64, err error) *Result {
	return &Result{
		serialNo:   serialNo,
		statusCode: statusCode,
		bodySize:   bodySize,
		duration:   duration,
		err:        err,
	}
}

func (task *Task) StartTime() time.Time {
	return task.startTime
}

func (task *Task) EndTime() time.Time {
	return task.endTime
}

func (task *Task) Results() []*Result {
	return task.results
}

func (task *Task) Start() {
	// 准备接受测试指标协程启动
	task.prepare()

	task.startTime = time.Now()
	// 开始往channel发送数据进行压测
	go func() {
		for i := 0; i < task.count; i++ {
			task.taskIndexChan <- i
		}
	}()
}

func (task *Task) prepare() {
	// 创建并发协程
	for i := 0; i < task.concurrent; i++ {
		task.startGoroutines(i)
	}
	go func() {
		for {
			select {
			case httpResult, opened := <-task.resultChan:
				if !opened {
					return
				}
				task.results = append(task.results, httpResult)
				if len(task.results) >= task.count {
					task.endTime = time.Now()
					close(task.taskIndexChan)
					close(task.resultChan)
					if task.resultFunc != nil {
						task.resultFunc(task)
					}
				}
			}
		}
	}()
}

func (task *Task) startGoroutines(routineIndex int) {
	go func() {
		for {
			select {
			case index, opened := <-task.taskIndexChan:
				if !opened {
					return
				}
				if task.requestFunc != nil {
					httpResult := task.requestFunc(index, routineIndex, defaultHttpClient)
					task.resultChan <- httpResult
				}
			}
		}
	}()
}

func (task *Task) PrintResults() {
	var maxDuration int64 = 0
	var minDuration int64 = math.MaxInt64
	var totalBodySize int64 = 0
	durationMap := make(map[int64]int64, 0)
	statusCodeMap := make(map[int]int64, 0)
	for _, result := range task.results {
		duration := result.Duration()
		if duration <= 10 {
			durationMap[10]++
		} else if duration <= 50 {
			durationMap[50]++
		} else if duration <= 100 {
			durationMap[50]++
		} else if duration <= 500 {
			durationMap[500]++
		} else if duration <= 1000 {
			durationMap[1000]++
		} else if duration <= 5000 {
			durationMap[5000]++
		} else if duration <= 10000 {
			durationMap[10000]++
		}
		statusCodeMap[result.StatusCode()]++
		if maxDuration < duration {
			maxDuration = duration
		}
		if minDuration >= duration {
			minDuration = duration
		}
		totalBodySize += result.BodySize()
	}
	cost := task.endTime.UnixMilli() - task.startTime.UnixMilli()
	println(fmt.Sprintf("Count: %d, cost: %d ms, max_duartion: %d ms, min_duration: %d ms, total_body_size: %d",
		len(task.results), cost, maxDuration, minDuration, totalBodySize))
	for k, v := range durationMap {
		println(fmt.Sprintf("Duration less than %d ms, count: %d", k, v))
	}
	for k, v := range statusCodeMap {
		println(fmt.Sprintf("StatusCode: %d , count: %d", k, v))
	}
}

func NewHttpTestTask(count, concurrent int, request RequestFunc, result ResultFunc) *Task {
	return &Task{
		count:         count,
		concurrent:    concurrent,
		results:       make([]*Result, 0),
		resultChan:    make(chan *Result),
		taskIndexChan: make(chan int),
		requestFunc:   request,
		resultFunc:    result,
	}
}
