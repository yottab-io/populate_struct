package populate_struct

import "errors"

var (
	InvalidPath = errors.New("invalid path: encountered a non-map type before reaching the end of the path")
	FieldNotFound = errors.New("field not found")
)
