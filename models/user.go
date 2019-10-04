package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type User struct {
	UserId     int       `orm:"pk;quto"`
	UserName   string    `orm:"unique"`
	NickName   string    `orm:"size(100)"`
	Password   string    `orm:"size(100)"`
	Salt       string    `orm:"size(10)"`
	UpdateTime time.Time `orm:"quto_now;type(datetime)"`
}

func init() {
	orm.RegisterModelWithPrefix("b_", new(User))
}
