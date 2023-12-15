package base

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/nico612/voyage/api/adminsrv/v1"
	"github.com/nico612/voyage/internal/adminsrv/auth"
	"github.com/nico612/voyage/internal/adminsrv/models"
	"github.com/nico612/voyage/internal/pkg/code"
	"github.com/nico612/voyage/internal/pkg/response"
	"github.com/nico612/voyage/pkg/errors"
	"github.com/nico612/voyage/pkg/log"
	"github.com/redis/go-redis/v9"
)

// Login 登录
// @Tags      Base
// @Summary   用户登录
// @Security  BearerTokenAuth
// @accept    application/json
// @Produce   application/json
// @Success   200  {object}  response.Response{data=v1.LoginResp,message=string}  "生成验证码,返回包括随机数id,base64,验证码长度,是否开启验证码"
// @Router    /base/captcha [post]
func (ctr *BaseController) Login(c *gin.Context) {

	var loginReq v1.LoginReq
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		log.Errorf("parse login parameters: %s", err.Error())
		response.Failed(c, errors.WrapC(err, code.ErrBind, "请求参数错误"))
		return
	}

	// 验证验证码
	if ctr.cfg.GenericServerRunOptions.Mode != "debug" {
		if !ctr.captcha.Verify(loginReq.CaptchaId, loginReq.Captcha) {
			response.Failed(c, errors.WithCode(code.ErrInvalidCaptcha, "无效的图形验证码"))
			return
		}
	}

	// 验证用户名和密码
	user, err := ctr.service.SysUsers().GetUserWithUsername(c, loginReq.Username)
	if err != nil {
		response.Failed(c, errors.WithCode(code.ErrUserNotFound, "login not found user error: %s", err.Error()))
		return
	}

	// 生成token
	token, expireAt, err := ctr.GenerateToken(user)
	if err != nil {
		log.Errorf("generator token error: %s", err.Error())
		response.Failed(c, errors.WithCode(code.ErrUnknown, "generator token error: %s", err.Error()))
		return
	}

	// 多点登录拦截
	if ctr.cfg.GenericServerRunOptions.UseMultipoint {
		if err = ctr.multipointLoginLimit(c, token, user); err != nil {
			err = errors.WrapC(err, code.ErrUnknown, "use multipoint limit error: %s", err.Error())
			log.Errorf("use multipoint error : %v", err)
			response.Failed(c, err)
			return
		}
	}

	resp := v1.LoginResp{
		User:      user,
		Token:     token,
		ExpiresAt: expireAt,
	}

	response.Success(c, resp)
}

// GenerateToken 签发token
func (ctr *BaseController) GenerateToken(user *models.SysUser) (string, int64, error) {
	jwtOpt := ctr.cfg.JwtOptions

	// 构建 token 第二段 payload 数据
	cliams := auth.NewClaims(&auth.Options{
		Key:        jwtOpt.Key,
		Timeout:    jwtOpt.Timeout,
		MaxRefresh: jwtOpt.MaxRefresh,
	})
	cliams.UserID = user.ID
	cliams.Username = user.Username
	cliams.AuthorityId = user.AuthorityId

	// 生成token
	jwtAuth := auth.NewJwtAuth([]byte(jwtOpt.Key))
	token, err := jwtAuth.GeneratorToken(cliams)
	if err != nil {
		return "", 0, err
	}

	return token, cliams.ExpiresAt.Unix(), nil
}

// 多端登录拦截
func (ctr *BaseController) multipointLoginLimit(ctx *gin.Context, newToken string, user *models.SysUser) error {

	// 取出之前 redis 中储存的token
	oldToken, err := ctr.service.JwtBlack().GetRedisJWT(ctx, user.ID)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			// 将现在的token储存到redis中, 并返回
			return ctr.service.JwtBlack().SetRedisJWT(ctx, user.ID, newToken, ctr.cfg.JwtOptions.Timeout)
		}
		return err
	}

	// 将之前的token加入黑名单
	blackJwt := &models.JwtBlackList{Jwt: oldToken}
	return ctr.service.JwtBlack().JoinJwtBlack(ctx, blackJwt)
}
