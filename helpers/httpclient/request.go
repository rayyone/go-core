package httpclient

import (
	"github.com/rayyone/go-core/helpers/maps"
)

// BodyParams Query Paramsv
type BodyParams map[string]interface{}

// Sanitize Remove key with nil value
func (b BodyParams) Sanitize() BodyParams {
	_ = maps.Sanitize(b)

	return b
}
