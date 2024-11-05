package logger

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"time"
)

type LogCallback = func(log *RequestLog)

type Logger struct {
	Callback LogCallback
}

func NewLogger(callback LogCallback) *Logger {
	return &Logger{
		Callback: callback,
	}
}
func (l *Logger) Log(ctx *fiber.Ctx) error {
	startTime := time.Now().UTC()
	chainErr := ctx.Next()
	if chainErr != nil {
		if err := ctx.App().ErrorHandler(ctx, chainErr); err != nil {
			_ = ctx.SendStatus(fiber.StatusInternalServerError)
		}
	}
	finishTime := time.Now().UTC()
	newRequest := NewLog(ctx)
	newRequest.Latency = finishTime.Sub(startTime).String()
	newRequest.CreatedAt = startTime
	err := json.Unmarshal(ctx.Body(), &newRequest.Body)
	if err != nil {
		newRequest.Body = make(map[string]string)
	}
	l.Callback(newRequest)
	printLog(newRequest)
	return nil
}

func printLog(log *RequestLog) {
	fmt.Printf("UTC:%s => %3s %s%s  %3d  | %7v \n",
		log.CreatedAt.Format(time.TimeOnly),
		log.Method,
		log.IP,
		log.Path,
		log.Status,
		log.Latency,
	)
}
