package http

import (
	"github.com/labstack/echo"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

type queryEntry struct {
	method       string
	queryKeyHash string
}

type queryMux struct {
	entries map[queryEntry]echo.HandlerFunc
}

func NewQueryMux() *queryMux {
	return &queryMux{
		entries: make(map[queryEntry]echo.HandlerFunc),
	}
}

func (r *queryMux) HEAD(path string, f echo.HandlerFunc) {
	r.addEntry(http.MethodHead, path, f)
}

func (r *queryMux) GET(path string, f echo.HandlerFunc) {
	r.addEntry(http.MethodGet, path, f)
}

func (r *queryMux) POST(path string, f echo.HandlerFunc) {
	r.addEntry(http.MethodPost, path, f)
}

func (r *queryMux) PUT(path string, f echo.HandlerFunc) {
	r.addEntry(http.MethodPut, path, f)
}

func (r *queryMux) DELETE(path string, f echo.HandlerFunc) {
	r.addEntry(http.MethodDelete, path, f)
}

func (r *queryMux) addEntry(method string, path string, f echo.HandlerFunc) {
	values, err := url.ParseQuery(path)
	if err != nil {
		panic(err)
	}
	queryKeyHash := getKeyHash(values)

	entry := queryEntry{
		method:       method,
		queryKeyHash: queryKeyHash,
	}

	r.entries[entry] = f
}

func (r *queryMux) Handle() echo.HandlerFunc {
	return func(c echo.Context) error {
		method := c.Request().Method
		queryKeysHash := getKeyHash(c.Request().URL.Query())

		pattern := queryEntry{
			method:       method,
			queryKeyHash: queryKeysHash,
		}

		if handler, ok := r.entries[pattern]; ok {
			return handler(c)
		}

		return c.NoContent(http.StatusNotFound)
	}
}

func getKeyHash(values url.Values) string {
	var keys []string
	for key := range values {
		keys = append(keys, key)
	}

	sort.Strings(keys)
	return strings.Join(keys, ",")
}
