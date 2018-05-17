// Package middleware provdes a logging middleware for the github.com/labstack/echo framework which sends information about the requests to Graylog
package middleware

import (
	"fmt"
	"time"

	"code.solebrity.com/backend-engineering/solebrityv2/gelf"
	"github.com/labstack/echo"
)

var graylogger *gelf.Logger

// LoggerMiddleware creates an echo.MiddlwareFunc to wrap handler
func LoggerMiddleware(host string, port int, applicationName string) echo.MiddlewareFunc {
	graylogger = gelf.NewLogger(host, port, applicationName)

	return func(nextHandler echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// if config.Skipper(c) {
			// return nextHandler(c)
			// }

			req := c.Request()
			res := c.Response()
			start := time.Now()
			var err error
			if err = nextHandler(c); err != nil {
				c.Error(err)
			}
			stop := time.Now()

			t := stop.Sub(start).Nanoseconds() / 1000000.0
			fmt.Println(stop.Sub(start).Nanoseconds())

			shortM := fmt.Sprintf("%d %s %s (%s) %d", res.Status(), req.Method(), req.URI(), req.RealIP(), t)
			longM := shortM
			if err != nil {
				longM += "\n" + err.Error()
			}

			metadata := map[string]interface{}{
				"status":        res.Status(),
				"ip":            req.RealIP(),
				"path":          req.URI(),
				"method":        req.Method(),
				"response_time": t}

			if res.Status() >= 500 {
				graylogger.Critical(shortM, longM, metadata)
			} else if res.Status() >= 400 {
				graylogger.Error(shortM, longM, metadata)
			} else if res.Status() >= 300 {
				graylogger.Notice(shortM, longM, metadata)
			} else {
				graylogger.Info(shortM, longM, metadata)
			}

			return nil
		}
	}
}
