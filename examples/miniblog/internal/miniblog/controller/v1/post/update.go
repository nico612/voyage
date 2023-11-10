package post

import (
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/nico612/go-project/examples/miniblog/internal/pkg/core"
	"github.com/nico612/go-project/examples/miniblog/internal/pkg/errno"
	"github.com/nico612/go-project/examples/miniblog/internal/pkg/known"
	"github.com/nico612/go-project/examples/miniblog/internal/pkg/log"
	v1 "github.com/nico612/go-project/examples/miniblog/pkg/api/miniblog/v1"
)

// Update 更新博客.
func (ctrl *PostController) Update(c *gin.Context) {
	log.C(c).Infow("Update post function called")

	var r v1.UpdatePostRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteResponse(c, errno.ErrBind, nil)

		return
	}

	if _, err := govalidator.ValidateStruct(r); err != nil {
		core.WriteResponse(c, errno.ErrInvalidParameter.SetMessage(err.Error()), nil)

		return
	}

	if err := ctrl.b.Posts().Update(c, c.GetString(known.XUsernameKey), c.Param("postID"), &r); err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	core.WriteResponse(c, nil, nil)
}
