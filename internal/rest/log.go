package rest

import (
	"bytes"
	"encoding/json"
	"time"

	"stavki/internal/log"

	"github.com/gin-gonic/gin"
)

type (
	bodyLogWriter struct {
		gin.ResponseWriter
		body *bytes.Buffer
	}
)

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func logMiddleware(logger log.FieldLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		c.Next()

		status := c.Writer.Status()
		entry := logger.WithFields(log.Fields{
			"method":     c.Request.Method,
			"url":        c.Request.Host + c.Request.URL.String(),
			"status":     status,
			"latency":    time.Since(start),
			"ip":         c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
		})

		if userID, ok := c.Get("userID"); ok {
			entry = entry.WithField("user_id", userID)
		}

		var errResp struct {
			Error string `json:"error"`
		}
		if err := json.Unmarshal(blw.body.Bytes(), &errResp); err == nil && errResp.Error != "" {
			logger = logger.WithField("error_msg", errResp.Error)
		}

		if len(c.Errors) > 0 {
			entry = entry.WithField("errors", c.Errors.String())
		}
		if status >= 500 {
			entry.Error("server error")
		} else if status >= 400 {
			entry.Warn("client error")
		} else {
			entry.Info("")
		}

	}
}
