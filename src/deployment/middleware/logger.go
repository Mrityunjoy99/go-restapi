package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"time"

	"github.com/Mrityunjoy99/sample-go/src/tools/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type logEntry struct {
	Timestamp string      `json:"timestamp"`
	Method    string      `json:"method"`
	Path      string      `json:"path"`
	Status    int         `json:"status"`
	Duration  string      `json:"duration"`
	ClientIP  string      `json:"client_ip"`
	RequestID string      `json:"request_id,omitempty"`
	Request   interface{} `json:"request,omitempty"`
	Response  interface{} `json:"response,omitempty"`
}

func LoggerMiddleware(logger logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Read and log request body
		var requestBody interface{}
		if c.Request.Body != nil {
			bodyBytes, _ := io.ReadAll(c.Request.Body)
			// Restore the io.ReadCloser to its original state
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			// Try to unmarshal the request body as JSON
			if len(bodyBytes) > 0 {
				var bodyMap map[string]interface{}
				if err := json.Unmarshal(bodyBytes, &bodyMap); err == nil {
					requestBody = bodyMap
				} else {
					requestBody = string(bodyBytes)
				}
			}
		}

		// Capture response body
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		// Process request
		c.Next()

		// Log after request is processed
		duration := time.Since(start)

		// Parse response body if it's JSON
		var responseBody interface{}
		if blw.body.Len() > 0 {
			var respMap map[string]interface{}
			if err := json.Unmarshal(blw.body.Bytes(), &respMap); err == nil {
				responseBody = respMap
			} else {
				responseBody = blw.body.String()
			}
		}

		logger.Info("Incoming Request Log", zap.String("method", c.Request.Method), zap.String("path", c.Request.URL.Path), zap.Int("status", c.Writer.Status()), zap.String("duration", duration.String()), zap.String("client_ip", c.ClientIP()), zap.Any("request", requestBody), zap.Any("response", responseBody))
	}
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
