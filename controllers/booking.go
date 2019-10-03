package controllers

import (
	"github.com/astaxie/beego/orm"
)

type BookingController struct {
	Base
}

func (c *BookingController) Show() {
	o:=orm.NewOrm()
	var rooms []
}
