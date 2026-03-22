package quota

import "errors"

var (
	ErrNilRule = errors.New("quota: rule is nil")
)
