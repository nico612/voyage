package models

type JwtBlackList struct {
	BaseModel
	Jwt string `gorm:"type:text;comment:jwtToken"`
}
