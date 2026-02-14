package httpx

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

// Parse populates a request struct from the HTTP request.
// It reads JSON body fields, path parameters (via chi), and query/form parameters.
func Parse(r *http.Request, v interface{}) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return nil
	}
	val = val.Elem()
	typ := val.Type()

	// Parse JSON body for POST/PUT/PATCH with content
	if r.Body != nil && r.ContentLength != 0 &&
		(r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch) {
		if err := json.NewDecoder(r.Body).Decode(v); err != nil && err.Error() != "EOF" {
			return err
		}
	}

	// Parse path params and form/query params from struct tags
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldVal := val.Field(i)
		if !fieldVal.CanSet() {
			continue
		}

		// Path params: `path:"id"`
		if pathTag := field.Tag.Get("path"); pathTag != "" {
			if param := chi.URLParam(r, pathTag); param != "" {
				setFieldValue(fieldVal, param)
			}
		}

		// Query/form params: `form:"page"`
		if formTag := field.Tag.Get("form"); formTag != "" {
			parts := strings.Split(formTag, ",")
			name := parts[0]
			if qv := r.URL.Query().Get(name); qv != "" {
				setFieldValue(fieldVal, qv)
			}
		}
	}

	return nil
}

func setFieldValue(field reflect.Value, value string) {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if n, err := strconv.ParseInt(value, 10, 64); err == nil {
			field.SetInt(n)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if n, err := strconv.ParseUint(value, 10, 64); err == nil {
			field.SetUint(n)
		}
	case reflect.Float32, reflect.Float64:
		if f, err := strconv.ParseFloat(value, 64); err == nil {
			field.SetFloat(f)
		}
	case reflect.Bool:
		if b, err := strconv.ParseBool(value); err == nil {
			field.SetBool(b)
		}
	}
}

// OkJson writes a JSON success response with status 200.
func OkJson(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(v)
}

// ErrorResponse writes a JSON error response.
func ErrorResponse(w http.ResponseWriter, err error) {
	code := http.StatusInternalServerError
	msg := err.Error()

	// Map common error messages to appropriate status codes
	lower := strings.ToLower(msg)
	switch {
	case strings.Contains(lower, "not found"):
		code = http.StatusNotFound
	case strings.Contains(lower, "unauthorized") || strings.Contains(lower, "invalid token"):
		code = http.StatusUnauthorized
	case strings.Contains(lower, "forbidden") || strings.Contains(lower, "access denied"):
		code = http.StatusForbidden
	case strings.Contains(lower, "invalid") || strings.Contains(lower, "required") ||
		strings.Contains(lower, "missing") || strings.Contains(lower, "bad request"):
		code = http.StatusBadRequest
	case strings.Contains(lower, "conflict") || strings.Contains(lower, "already exists") ||
		strings.Contains(lower, "duplicate"):
		code = http.StatusConflict
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
