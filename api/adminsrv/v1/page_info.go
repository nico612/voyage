package v1

import "github.com/nico612/voyage/pkg/validator"

type PageInfo struct {
	Page     int    `json:"page" form:"page"`         // 页码
	PageSize int    `json:"pageSize" form:"pageSize"` // 每页大小
	Total    int64  `json:"total"`                    //总数 作为参数不用填写
	Offset   int    `json:"-"`
	Limit    int    `json:"-"`
	Keyword  string `json:"keyword" form:"keyword"` // 关键字
}

func (p *PageInfo) Validate() error {
	if p.Page < 1 {
		p.Page = 1
	}

	if p.PageSize == 0 {
		p.PageSize = 10
	}

	if p.PageSize > 100 {
		p.PageSize = 100
	}

	p.Offset = p.PageSize * (p.Page - 1)
	p.Limit = p.PageSize

	return nil
}

type PageResult struct {
	List interface{} `json:"list"`
	PageInfo
}

// GetById Find by id structure
type GetById struct {
	ID uint `json:"id" form:"id"` // 主键ID
}

func (r GetById) Validate() error {
	idVerify := validator.Rules{"ID": []string{validator.NotEmpty()}}
	return validator.Verify(r, idVerify)
}

type IdsReq struct {
	Ids []int `json:"ids" form:"ids"`
}

// GetAuthorityId Get role by id structure
type GetAuthorityId struct {
	AuthorityId uint `json:"authorityId" form:"authorityId"` // 角色ID
}

type Empty struct{}
