package http

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	log "github.com/sirupsen/logrus"
	"os"
)

func ServerHTTP(httpPort int, accessLogFile *os.File, route func(e *echo.Echo)) {
	accessLogMiddleware := newAccessLog(accessLogFile)

	e := echo.New()
	e.Use(accessLogMiddleware, middleware.Recover())
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		if _, ok := err.(*echo.HTTPError); !ok {
			log.WithError(err).Errorf("handle http error: method=%s, url=%s", c.Request().Method, c.Request().URL)
		}

		if err != nil {
			if err := jsonError(c, err); err != nil {
				log.WithError(err).Errorf("handle http error: write json error, method=%s, url=%s", c.Request().Method, c.Request().URL)
			}
		}
	}

	route(e)
	if err := e.Start(fmt.Sprintf(":%d", httpPort)); err != nil {
		log.WithError(err).Fatalf("start http server error: port=%d", httpPort)
	}
}
