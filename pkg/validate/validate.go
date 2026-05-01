// TODO: Make `validate` package HTTP-router agnostic. Because now it kinda depends on `echo`

package validate

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

func Bind(c echo.Context, dst interface{}) error {
	if reflect.TypeOf(dst).Kind() != reflect.Ptr {
		return errors.New("validate: invalid dstination: expected pointer")
	}

	// NOTE: Validate used only on parsing JSON data, so it makes sense to set an `application/json` header here.
	// And sometimes I just don't want to set it manualy in curl
	c.Request().Header.Set("Content-Type", "application/json")
	if err := c.Bind(dst); err != nil {
		return fmt.Errorf("json.Unmarshall: %v", err)
	}

	if reflect.TypeOf(dst).Elem().Kind() == reflect.Map {
		return nil
	}

	return Struct(dst)
}

func Struct(dst interface{}) error {
	val := reflect.ValueOf(dst)

	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("expected a pointer to a struct")
	}

	var missing []string
	if err := walk(val.Elem(), "", &missing); err != nil {
		return err
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required fields: %s", strings.Join(missing, ", "))
	}
	return nil
}

func walk(val reflect.Value, prefix string, missing *[]string) error {
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		name := fieldName(fieldType)
		path := name
		if prefix != "" {
			path = prefix + "." + name
		}

		required := fieldType.Tag.Get("required") == "true"

		if def, ok := fieldType.Tag.Lookup("default"); ok && !required && isEmptyValue(field) && field.CanSet() {
			if err := setDefault(field, def); err != nil {
				return fmt.Errorf("field `%s`: %w", path, err)
			}
		}

		if required && isEmptyValue(field) {
			*missing = append(*missing, path)
		}

		if field.Kind() == reflect.Struct && fieldType.Type != reflect.TypeOf(time.Time{}) {
			var nested reflect.Value
			if field.CanAddr() {
				nested = field
			} else {
				tmp := reflect.New(field.Type()).Elem()
				tmp.Set(field)
				nested = tmp
			}
			if err := walk(nested, path, missing); err != nil {
				return err
			}
		}
	}
	return nil
}

func fieldName(f reflect.StructField) string {
	for _, tag := range []string{"json", "yaml"} {
		if v := f.Tag.Get(tag); v != "" {
			return strings.Split(v, ",")[0]
		}
	}
	return f.Name
}

func setDefault(field reflect.Value, def string) error {
	if field.Type() == reflect.TypeOf(time.Duration(0)) {
		d, err := time.ParseDuration(def)
		if err != nil {
			return err
		}
		field.SetInt(int64(d))
		return nil
	}

	switch field.Kind() {
	case reflect.String:
		field.SetString(def)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n, err := strconv.ParseInt(def, 10, 64)
		if err != nil {
			return err
		}
		field.SetInt(n)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		n, err := strconv.ParseUint(def, 10, 64)
		if err != nil {
			return err
		}
		field.SetUint(n)
	case reflect.Float32, reflect.Float64:
		n, err := strconv.ParseFloat(def, 64)
		if err != nil {
			return err
		}
		field.SetFloat(n)
	case reflect.Bool:
		b, err := strconv.ParseBool(def)
		if err != nil {
			return err
		}
		field.SetBool(b)
	default:
		return fmt.Errorf("unsupported default for kind %s", field.Kind())
	}
	return nil
}

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String:
		return v.String() == ""
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() <= 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() <= 0
	case reflect.Slice, reflect.Array:
		return v.Len() == 0
	case reflect.Map:
		return v.Len() == 0
	case reflect.Ptr, reflect.Interface:
		return v.IsNil()
	default:
		return false
	}
}
