package model

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"io"
	"time"
)

type ResponseBodyWriter struct {
	gin.ResponseWriter
	Body *bytes.Buffer
}

func (w *ResponseBodyWriter) Write(b []byte) (int, error) {
	w.Body.Write(b)
	return w.ResponseWriter.Write(b)
}

func LoggingMiddleWare(database *Database) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		startTime := time.Now()

		reqBody, _ := io.ReadAll(ctx.Request.Body)

		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(reqBody))

		w := &ResponseBodyWriter{Body: &bytes.Buffer{}, ResponseWriter: ctx.Writer}
		ctx.Writer = w

		ctx.Next()

		elapsedTime := time.Since(startTime)

		logEntry := &Log{
			ElapsedTime: elapsedTime.String(),
			Method:      ctx.Request.Method,
			Endpoint:    ctx.FullPath(),
			Query:       ctx.Request.URL.RawQuery,
			ReqBody:     string(reqBody),
			Code:        ctx.Writer.Status(),
			ResBody:     w.Body.String(),
		}

		var statement string = "INSERT INTO logs (elapsedTime, method, endpoint, query, reqBody, code, resBody) VALUES ($1, $2, $3, $4, $5, $6, $7);"
		database.database.Exec(statement, logEntry.ElapsedTime, logEntry.Method, logEntry.Endpoint, logEntry.Query, logEntry.ReqBody, logEntry.Code, logEntry.ResBody)
	}
}
