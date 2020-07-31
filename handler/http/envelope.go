package http

import (
	"encoding/json"
	"github.com/cloud-mesh/object-storage/model"
	"github.com/labstack/echo"
	"net/http"
)

type envelope struct {
	Code int
	Msg  string
	Data map[string]interface{}
}

func (e envelope) MarshalJSON() ([]byte, error) {
	var envelopMap map[string]interface{}
	if e.Data != nil {
		envelopMap = e.Data
	} else {
		envelopMap = make(map[string]interface{})
	}

	envelopMap["errcode"] = e.Code
	envelopMap["msg"] = e.Msg

	return json.Marshal(envelopMap)
}

func jsonOK(ctx echo.Context, data map[string]interface{}) error {
	return ctx.JSON(http.StatusOK, envelope{
		Code: model.ErrCodeOK,
		Data: data,
	})
}

func jsonError(ctx echo.Context, err error) error {
	if err == nil {
		return jsonOK(ctx, nil)
	}

	httpCode := http.StatusOK
	if err, ok := err.(interface {
		HTTPCode() int
	}); ok {
		httpCode = err.HTTPCode()
	}

	return ctx.JSON(httpCode, envelope{
		Code: model.GetCode(err),
		Msg:  err.Error(),
	})
}
