package engine

import (
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"

	jsoniter "github.com/json-iterator/go"
)

func handleError(w http.ResponseWriter, r *http.Request, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(fmt.Sprint(err)))
}

func makeAPIHandler[InType, OutType any](handler func(*InType) (*OutType, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var param InType
		json := jsoniter.ConfigCompatibleWithStandardLibrary
		if r.Method == http.MethodGet {
			if err := DecodeForm(r, &param); err != nil {
				handleError(w, r, err)
				return
			}
		} else if r.Method == http.MethodPost {
			data, err := io.ReadAll(r.Body)
			if err != nil {
				handleError(w, r, err)
				return
			}
			err = json.Unmarshal(data, &param)
			if err != nil {
				handleError(w, r, err)
				return
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		outCode := 0
		outMsg := ""
		out, err := handler(&param)
		if err != nil {
			outCode = 233
			outMsg = err.Error()
		}
		outMap := map[string]any{
			"code": outCode,
			"msg":  outMsg,
			"data": out,
		}

		outData, err := json.Marshal(outMap)
		if err != nil {
			handleError(w, r, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(outData)
	}
}

func DecodeForm(r *http.Request, ptr any) error {
	rValue := reflect.ValueOf(ptr).Elem()
	rType := rValue.Type()

	for i := 0; i < rValue.NumField(); i++ {
		if !rValue.Field(i).CanSet() {
			continue
		}
		if tag := rType.Field(i).Tag.Get("form"); tag != "" {
			formValue := r.FormValue(tag)
			if formValue == "" {
				continue
			}

			switch rType.Field(i).Type.Kind() {
			case reflect.String:
				rValue.Field(i).SetString(formValue)
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				value, err := strconv.ParseInt(formValue, 10, 64)
				if err != nil {
					return err
				}
				rValue.Field(i).SetInt(value)
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				value, err := strconv.ParseUint(formValue, 10, 64)
				if err != nil {
					return err
				}
				rValue.Field(i).SetUint(value)
			default:
				return fmt.Errorf("invalid field type for %s: %T", tag, rValue.Field(i).Interface())
			}
		}
	}
	return nil
}
