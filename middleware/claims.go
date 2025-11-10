package middleware

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/thk-im/thk-im-base-server/crypto"
	"github.com/thk-im/thk-im-base-server/dto"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type aesWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *aesWriter) Write(b []byte) (int, error) {
	return w.body.Write(b)
}

func (w *aesWriter) WriteString(s string) (int, error) {
	return w.body.WriteString(s)
}

const ClaimsKey = dto.ClaimsKey

func Claims(crypto crypto.Crypto) gin.HandlerFunc {
	return func(context *gin.Context) {
		claims := dto.ThkClaims{}
		traceID := context.Request.Header.Get(dto.TraceID)
		if strings.EqualFold(traceID, "") {
			traceID = uuid.NewString()
		}
		claims.PutValue(dto.TraceID, traceID)

		parentSpanID := context.Request.Header.Get(dto.SpanID)
		spanID := ""
		if parentSpanID == "" {
			parentSpanID = "0"
			spanID = "1"
		} else {
			i, err := strconv.Atoi(parentSpanID)
			if err == nil {
				spanID = fmt.Sprintf("%d", i+1)
			} else {
				spanID = "1"
			}
		}
		claims.PutValue(dto.ParentSpanID, parentSpanID)
		claims.PutValue(dto.SpanID, spanID)

		clientIP := context.Request.Header.Get(dto.OriginIP)
		claims.PutValue(dto.OriginIP, clientIP)

		device := getValueFromContext(context, dto.Device)
		claims.PutValue(dto.Device, device)

		deviceId := getValueFromContext(context, dto.DeviceId)
		claims.PutValue(dto.DeviceId, deviceId)

		timeZone := getValueFromContext(context, dto.TimeZone)
		claims.PutValue(dto.TimeZone, timeZone)

		platform := getValueFromContext(context, dto.Platform)
		claims.PutValue(dto.Platform, platform)

		channel := getValueFromContext(context, dto.Channel)
		claims.PutValue(dto.Channel, channel)

		version := getValueFromContext(context, dto.Version)
		claims.PutValue(dto.Version, version)

		language := getValueFromContext(context, dto.Language)
		claims.PutValue(dto.Language, language)

		token := getValueFromContext(context, dto.JwtToken)
		token = strings.ReplaceAll(token, "Bearer ", "")
		token = strings.ReplaceAll(token, " ", "")
		claims.PutValue(dto.JwtToken, token)

		context.Set(ClaimsKey, claims)

		if parentSpanID == "0" && crypto != nil && context.Request.RequestURI != "/metrics" {
			// 需要解密
			oldWriter := context.Writer
			blw := &aesWriter{body: bytes.NewBufferString(""), ResponseWriter: context.Writer}
			rawData, err := context.GetRawData()
			if err != nil {
				context.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			if len(rawData) > 0 {
				deData, errDecrypt := crypto.DecryptUriBody(context.Request.RequestURI, string(rawData))
				if errDecrypt != nil {
					context.AbortWithStatus(http.StatusBadRequest)
					return
				}
				context.Request.Body = io.NopCloser(bytes.NewBuffer(deData))
			}
			context.Writer = blw

			context.Next()

			// 需要加密
			context.Writer = oldWriter
			responseBytes := blw.body.Bytes()
			crData, errCrypt := crypto.EncryptUriBody(context.Request.RequestURI, responseBytes)
			if errCrypt != nil {
				context.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			_, _ = context.Writer.WriteString(crData)
		} else {
			context.Next()
		}

	}
}

func getValueFromContext(ctx *gin.Context, key string) string {
	value := ctx.Query(key)
	if value == "" {
		value = ctx.Request.Header.Get(key)
		if value == "" {
			value, _ = ctx.Cookie(key)
		}
	}
	return value
}
