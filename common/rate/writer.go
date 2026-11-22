package rate

import (
	"github.com/juju/ratelimit"
)

type Writer struct {
	limiter *ratelimit.Bucket
}
