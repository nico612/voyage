package user

import (
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/nico612/go-project/examples/miniblog/internal/pkg/core"
	"github.com/nico612/go-project/examples/miniblog/internal/pkg/log"
	v1 "github.com/nico612/go-project/examples/miniblog/pkg/api/miniblog/v1"
	"github.com/nico612/go-project/pkg/auth"
)

func (ctr *UserController) Create(c *gin.Context) {

	var (
		r   v1.CreateUserRequest
		err error
	)

	if err = c.ShouldBindJSON(&r); err != nil {
		log.C(c).Infow("Create user function called")
		core.WriteResponse(c, err, nil)
		return
	}

	// 参数验证
	if _, err = govalidator.ValidateStruct(r); err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	// 密码加密
	if r.Password, err = auth.Encrypt(r.Password); err != nil {
		log.C(c).Infow("auth encrypt err", "error", err)
		core.WriteResponse(c, err, nil)
		return
	}

	if err := ctr.b.Users().Create(c, &r); err != nil {
		log.C(c).Errorw("create user", "error", err.Error())
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, nil, nil)
}
