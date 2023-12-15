package validator

import "github.com/nico612/voyage/pkg/validator"

func Initialize() {
	_ = validator.RegisterRule("PageVerify",
		validator.Rules{
			"Page":     {validator.NotEmpty()},
			"PageSize": {validator.NotEmpty()},
		},
	)
	_ = validator.RegisterRule("IdVerify",
		validator.Rules{
			"Id": {validator.NotEmpty()},
		},
	)
	_ = validator.RegisterRule("AuthorityIdVerify",
		validator.Rules{
			"AuthorityId": {validator.NotEmpty()},
		},
	)
}
