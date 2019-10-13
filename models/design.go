package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type Design struct {
	DesignId   int    `orm:"pk;auto"`
	Key        string `orm:"unique"`
	Name       string
	Desc       string
	Status     int       `orm:"default(1)"`
	UpdateTime time.Time `orm:"auto_now;type(datetime)"`
}

func init() {
	orm.RegisterModelWithPrefix("b_", new(Design))
}
