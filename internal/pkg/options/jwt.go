package options

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"time"
)

type JwtOptions struct {
	Realm      string        `json:"realm"       mapstructure:"realm"` // jwt 标识
	Key        string        `json:"key"         mapstructure:"key"`
	Timeout    time.Duration `json:"timeout"     mapstructure:"timeout"`
	MaxRefresh time.Duration `json:"max-refresh" mapstructure:"max-refresh"`
}

func NewJwtOptions() *JwtOptions {
	return &JwtOptions{
		Realm:      "jwt",
		Key:        "",
		Timeout:    24 * time.Hour,
		MaxRefresh: 24 * time.Hour,
	}
}

func (s *JwtOptions) Validate() []error {
	var errs []error

	if !govalidator.StringLength(s.Key, "6", "32") {
		errs = append(errs, fmt.Errorf("--secret-key must larger than 5 and little than 33"))
	}

	return errs
}
