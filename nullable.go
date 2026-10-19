package sevk

import (
	"bytes"
	"encoding/json"
)

// Nullable represents an optional field that can carry either a value or an
// explicit JSON null. A nil *Nullable omits the field when used with omitempty.
type Nullable[T any] struct {
	value  T
	isNull bool
}

// NewNullable returns a Nullable carrying v.
func NewNullable[T any](v T) *Nullable[T] {
	return &Nullable[T]{value: v}
}

// Null returns a Nullable that marshals as JSON null.
func Null[T any]() *Nullable[T] {
	return &Nullable[T]{isNull: true}
}

// IsNull reports whether the field marshals as JSON null.
func (n *Nullable[T]) IsNull() bool {
	return n != nil && n.isNull
}

// Value returns the carried value and whether one is set.
func (n *Nullable[T]) Value() (T, bool) {
	if n == nil || n.isNull {
		var zero T
		return zero, false
	}
	return n.value, true
}

// MarshalJSON implements json.Marshaler.
func (n Nullable[T]) MarshalJSON() ([]byte, error) {
	if n.isNull {
		return []byte("null"), nil
	}
	return json.Marshal(n.value)
}

// UnmarshalJSON implements json.Unmarshaler.
func (n *Nullable[T]) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte("null")) {
		*n = Nullable[T]{isNull: true}
		return nil
	}

	var value T
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	*n = Nullable[T]{value: value}
	return nil
}
