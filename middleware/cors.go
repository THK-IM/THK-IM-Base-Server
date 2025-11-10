package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
	"strings"
)

// NewCors 创建 CORS 中间件
/*
r.Use(NewCors([]string{
	"https://abc.com",
	"https://web.abc.com",
	"*.google.com",
	"http://localhost:*", // 支持任意端口
}))
*/
func NewCors(allowed []string) gin.HandlerFunc {
	var patterns []*regexp.Regexp

	for _, rule := range allowed {
		rule = strings.TrimSpace(rule)
		if strings.HasPrefix(rule, "*.") {
			domain := strings.TrimPrefix(rule, "*.")
			patternStr := `^https?://([a-zA-Z0-9-]+\.)?` + regexp.QuoteMeta(domain) + `(:\d+)?$`
			patterns = append(patterns, regexp.MustCompile(patternStr))
		} else if strings.Contains(rule, "*") {
			patternStr := "^" + strings.ReplaceAll(regexp.QuoteMeta(rule), `\*`, ".*") + "$"
			patterns = append(patterns, regexp.MustCompile(patternStr))
		} else {
			patternStr := "^" + regexp.QuoteMeta(rule) + "$"
			patterns = append(patterns, regexp.MustCompile(patternStr))
		}
	}

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin == "" {
			c.Next()
			return
		}

		allowedOrigin := ""
		for _, p := range patterns {
			if p.MatchString(origin) {
				allowedOrigin = origin
				break
			}
		}

		if allowedOrigin != "" {
			c.Header("Access-Control-Allow-Origin", allowedOrigin)
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			c.Header("Access-Control-Expose-Headers", "*")
			c.Header("Access-Control-Allow-Headers", "*") // <-- 允许所有请求头
		}

		if c.Request.Method == http.MethodOptions {
			// OPTIONS 请求直接返回 200，带上 headers
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	}
}
