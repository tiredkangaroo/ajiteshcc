package server

import (
	"reflect"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
)

func handler[T any](fn func(c echo.Context, req T) error) echo.HandlerFunc {
	rt := reflect.TypeFor[T]()
	var requiredFields []int
	for i := range rt.NumField() {
		field := rt.Field(i)
		requiredStr := field.Tag.Get("required")
		if requiredStr == "false" {
			continue
		}
		requiredFields = append(requiredFields, i)
	}
	return func(c echo.Context) error {
		var req T
		if err := c.Bind(&req); err != nil {
			return err
		}
		rv := reflect.ValueOf(req)
		for _, fieldNumber := range requiredFields {
			if rv.Field(fieldNumber).IsZero() {
				return c.String(400, "bad request: missing required field "+rt.Field(fieldNumber).Name)
			}
		}
		return fn(c, req)
	}
}

func pgText(s string) pgtype.Text {
	return pgtype.Text{
		String: s,
		Valid:  true,
	}
}
