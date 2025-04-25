package middleware

import (
	"bytes"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Read request body
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// Create response writer
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		// Process request
		c.Next()

		// Stop timer
		duration := time.Since(start)

		// Log request and response
		logrus.WithFields(logrus.Fields{
			"client_ip":    c.ClientIP(),
			"method":       c.Request.Method,
			"path":         c.Request.URL.Path,
			"status_code":  c.Writer.Status(),
			"duration":     duration,
			"user_agent":   c.Request.UserAgent(),
			"request_body": string(requestBody),
			"response":     blw.body.String(),
		}).Info("HTTP Request")
	}
} 