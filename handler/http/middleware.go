package http

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"net/http"
	"os"
	"strings"
)

var cors = middleware.CORSWithConfig(middleware.CORSConfig{
	Skipper:      middleware.DefaultSkipper,
	AllowOrigins: []string{"*"},
	AllowMethods: []string{http.MethodOptions, http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
})

const (
	domainParamRegion = "region"
	domainParamBucket = "bucket"
)

var domainParam = func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		host := c.Request().Host
		parts := strings.Split(host, ".")

		bucketName := ""
		regionName := ""
		switch len(parts) {
		case 3:
			// {region}.example.com
			regionName = parts[0]
		case 4:
			// {bucket_name}.{region}.example.com
			bucketName = parts[0]
			regionName = parts[1]
		}

		if regionName == "" {
			regionName = c.QueryParam("region")
		}
		if bucketName == "" {
			bucketName = c.QueryParam("bucket")
		}

		c.Set(domainParamRegion, regionName)
		c.Set(domainParamBucket, bucketName)
		return next(c)
	}
}

func getRegionName(c echo.Context) string {
	regionName, _ := c.Get(domainParamRegion).(string)
	return regionName
}

func getBucketName(c echo.Context) string {
	bucketName, _ := c.Get(domainParamBucket).(string)
	return bucketName
}

func newAccessLog(logFile *os.File) echo.MiddlewareFunc {
	return middleware.LoggerWithConfig(middleware.LoggerConfig{Output: logFile})
}
