package httpcodec

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/types"
	"io"
	"net/http"
	"reflect"

	"github.com/RussellLuo/kun/pkg/werror"
	"github.com/RussellLuo/kun/pkg/werror/gcode"
)

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type FailureResponse struct {
	Error Error `json:"error"`
}

var customTypes = map[string]reflect.Type{}

type JSON struct{}

func (j JSON) SetCustomType(name string, ct reflect.Type) error {
	if _, ok := customTypes[name]; ok {
		return gcode.ErrInvalidArgument
	}
	customTypes[name] = ct
	return nil
}

func (j JSON) DecodeRequestParam(name string, values []string, out interface{}) error {
	if err := defaultBasicParam.Decode(values, out); err != nil {
		if err == ErrUnsupportedType {
			panic(fmt.Errorf("DecodeRequestParam not implemented for %q (of type %T)", name, out))
		}
		return werror.Wrap(gcode.ErrInvalidArgument, err)
	}
	return nil
}

func (j JSON) DecodeRequestParams(name string, values map[string][]string, out interface{}) error {
	if err := defaultStructParams.Decode(values, out); err != nil {
		if err == ErrUnsupportedType {
			panic(fmt.Errorf("DecodeRequestParams not implemented for %q (of type %T)", name, out))
		}
		return werror.Wrap(gcode.ErrInvalidArgument, err)
	}
	return nil
}

func (j JSON) DecodeRequestBody(r *http.Request, out interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(out); err != nil {
		return werror.Wrap(gcode.ErrInvalidArgument, err)
	}
	return nil
}

func (j JSON) DecodeRequestBodyFromType(r *http.Request, out interface{}, modifierFunc func() (io.Reader, error)) error {
	x, err := modifierFunc()
	if err != nil {
		return werror.Wrap(gcode.ErrInvalidArgument, err)
	}

	if err := json.NewDecoder(x).Decode(out); err != nil {
		return werror.Wrap(gcode.ErrInvalidArgument, err)
	}
	return nil
}

func (j JSON) EncodeSuccessResponse(w http.ResponseWriter, statusCode int, body interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(body)
}

func (j JSON) EncodeFailureResponse(w http.ResponseWriter, err error) error {
	statusCode := gcode.HTTPStatusCode(err)
	code, message := gcode.ToCodeMessage(err)
	return j.EncodeSuccessResponse(w, statusCode, FailureResponse{
		Error: Error{
			Code:    code,
			Message: message,
		},
	})
}

func (j JSON) EncodeRequestParam(name string, value interface{}) []string {
	return defaultBasicParam.Encode(value)
}

func (j JSON) EncodeRequestParams(name string, value interface{}) map[string][]string {
	return defaultStructParams.Encode(value)
}

func (j JSON) EncodeRequestBody(body interface{}) (io.Reader, map[string]string, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return nil, nil, err
	}
	headers := map[string]string{
		"Content-Type": "application/json; charset=utf-8",
	}
	return bytes.NewBuffer(data), headers, nil
}

func (j JSON) DecodeSuccessResponse(body io.ReadCloser, out interface{}) error {
	return json.NewDecoder(body).Decode(out)
}

func (j JSON) DecodeFailureResponse(body io.ReadCloser, out *error) error {
	var resp FailureResponse
	if err := json.NewDecoder(body).Decode(&resp); err != nil {
		return err
	}

	*out = gcode.FromCodeMessage(resp.Error.Code, resp.Error.Message)
	return nil
}

func (j JSON) RegisterTypes(types map[string]types.Type) error {
	return gcode.ErrNotImplemented
}

func (j JSON) RetrieveType(name string) (types.Type, error) {
	return nil, gcode.ErrNotImplemented
}