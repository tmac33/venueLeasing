package controllers

import (
	"venueleasing/models"

	"github.com/astaxie/beego/orm"
)

type BookingController struct {
	Base
}

func (c *BookingController) Show() {
	o := orm.NewOrm()
	var rooms []*models.Room
	o.QueryTable("b_room").Filter("status", 1).All(&rooms)
	c.Data["RoomList"] = rooms
	var designs []*models.Design
	o.QueryTable("b_design").Filter("status", 1).All(&designs)
	c.Data["DesignList"] = designs
	c.TplName = "NewBooking.html"
}
