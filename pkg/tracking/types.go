package tracking

import "github.com/wehmoen-dev/ronin-buddy/pkg/config"

type Sentry struct {
	enabled bool
	config  *config.ActionConfig
}
