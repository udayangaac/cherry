package cherry

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"reflect"
	"strconv"
)

// Content types (v1)
// application/x-www-form-urlencoded
// application/json
const (
	QueryParamTag = "param"
	BodyTag       = "body"
	PathVariable  = "path"
	Json          = "json"
)

type Reader interface {
	Read(ctx context.Context, req *http.Request, i interface{}) (err error)
}

type reader struct{}

func NewReader() Reader {
	return &reader{}
}

func (r *reader) Read(ctx context.Context, req *http.Request, i interface{}) (err error) {
	err = r.read(ctx, req, i)
	return
}

func (r *reader) read(_ context.Context, req *http.Request, i interface{}) (err error) {
	fmt.Println("Request URL", req.URL.Path)
	rv := reflect.ValueOf(i)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errors.New("given value isn't a pointer")
	}
	rt := reflect.TypeOf(i).Elem()
	queryParams := req.URL.Query()
	vars := mux.Vars(req)
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		if val, ok := f.Tag.Lookup(QueryParamTag); ok {
			if err = r.setValue(rv.Elem().Field(i), queryParams.Get(val)); err != nil {
				return
			}
		}
		if val, ok := f.Tag.Lookup(BodyTag); ok {
			if Json == val {
				if err = json.NewDecoder(req.Body).Decode(rv.Elem().Field(i).Interface()); err != nil {
					return
				}
			}
		}
		if val, ok := f.Tag.Lookup(PathVariable); ok {
			if pathVal, isExist := vars[val]; isExist {
				if err = r.setValue(rv.Elem().Field(i), pathVal); err != nil {
					return
				}
			} else {
				return errors.New(fmt.Sprintf("no path variable: %s", val))
			}
		}
	}
	return
}

func (r *reader) setValue(value reflect.Value, readValue string) (err error) {
	// Return for empty values
	if len(readValue) == 0 {
		return
	}
	k := value.Kind()
	switch k {
	case reflect.Bool:
		var convertedVal bool
		convertedVal, err = strconv.ParseBool(readValue)
		value.SetBool(convertedVal)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var convertedVal int64
		convertedVal, err = strconv.ParseInt(readValue, 10, 64)
		value.SetInt(convertedVal)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		var convertedVal uint64
		convertedVal, err = strconv.ParseUint(readValue, 10, 64)
		value.SetUint(convertedVal)

	case reflect.Float32, reflect.Float64:
		var convertedVal float64
		convertedVal, err = strconv.ParseFloat(readValue, 64)
		value.SetFloat(convertedVal)

	case reflect.Complex64, reflect.Complex128:
		var convertedVal complex128
		convertedVal, err = strconv.ParseComplex(readValue, 64)
		value.SetComplex(convertedVal)
	case reflect.Array:
	case reflect.String:
		value.SetString(readValue)
	case reflect.Ptr:
		err = json.Unmarshal([]byte(readValue), value.Interface())
	}
	return
}
