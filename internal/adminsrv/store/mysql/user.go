package mysql

import "gorm.io/gorm"

type User struct {
	*gorm.DB
}

func newUsers(ds *datasource) *User {
	return &User{ds.db}
}

func (u *User) Login() error {

	return nil
}
