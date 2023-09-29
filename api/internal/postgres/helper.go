package postgres

import "reflect"

func pointerIfNotZero[T any](v T) *T {
	if reflect.ValueOf(v).IsZero() {
		return nil
	}
	return &v
}
