package sevk

// Bool returns a pointer to v for use with optional boolean fields.
func Bool(v bool) *bool { return &v }

// Int returns a pointer to v for use with optional integer fields.
func Int(v int) *int { return &v }

// String returns a pointer to v for use with optional string fields.
func String(v string) *string { return &v }
