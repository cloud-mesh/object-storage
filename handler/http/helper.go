package http

import (
	"github.com/labstack/echo"
	"strconv"
	"time"
)

const (
	defaultPageSize = 20
	timeLayout      = "2006-01-02 15:04:05"
)

func getInt(val string, defaultVal int) int {
	num, err := strconv.Atoi(val)
	if err != nil {
		return defaultVal
	}
	return num
}

func getTime(val string) *time.Time {
	t, err := time.Parse(timeLayout, val)
	if err != nil {
		return nil
	}

	return &t
}

func getPaging(c echo.Context) (page int, pageSize int) {
	pageVal := c.QueryParam("page")
	pageSizeVal := c.QueryParam("page_size")

	page, _ = strconv.Atoi(pageVal)
	pageSize, _ = strconv.Atoi(pageSizeVal)
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = defaultPageSize
	}
	return
}
