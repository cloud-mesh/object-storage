package http

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"net/http"
	"os"
)


var cors = middleware.CORSWithConfig(middleware.CORSConfig{
	Skipper:      middleware.DefaultSkipper,
	AllowOrigins: []string{"*"},
	AllowMethods: []string{http.MethodOptions, http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
})

func newAccessLog(logFile *os.File) echo.MiddlewareFunc {
	return middleware.LoggerWithConfig(middleware.LoggerConfig{Output: logFile})
}
