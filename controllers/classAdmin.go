package controllers

import (
	"strings"
	"venueleasing/models"

	"github.com/astaxie/beego/orm"
)

func seralizeTime(str string) []string {
	return strings.Split(str, ",")
}

type ClassAdminController struct {
	BaseAdmin
}

//exhibite venue list
// @router /room/show [get,post]
func (c *ClassAdminController) Show() {
	if c.Ctx.Request.Method == "GET" {
		o := orm.NewOrm()
		var rooms []*models.Room
		o.QueryTable("b_room").All(&rooms)
		c.Data["Room"] = rooms
		c.TplName = "RoomList.html"
	}
}

//add venue
// @router /room/add [get,post]
func (c *ClassAdminController) Add() {
	if c.Ctx.Request.Method == "GET" {
		c.TplName = "RoomAdd.html"
	} else {
		o := orm.NewOrm()
		var room models.Room
		room.Name = c.GetString("name")
		room.StartTime = c.GetString("start_time")
		room.EndTime = c.GetString("end_time")
		room.EndTime = c.GetString("end_time")
		room.PosNum, _ = c.GetInt("pos_num")
		room.Desc = c.GetString("desc")
		room.Status = 1
		_, err := o.Insert(&room)
		if err != nil {
			c.AjaxErr("add fail")
		} else {
			c.AjaxOk("add success")
		}

	}
}

//remove venue
// @router /room/delete [post]
func (c *ClassAdminController) Delete() {
	o := orm.NewOrm()
	id, err := c.GetInt("room_id")
	if err == nil {
		if num, err := o.Delete(&models.Room{RoomId: id}); err == nil && num == 1 {
			c.AjaxOk("remove success")
		}
	}
	c.AjaxErr("remove failed")
}

//modify venue details
// @router /room/update/?:id [get,post]
func (c *ClassAdminController) Update() {
	if c.Ctx.Request.Method == "GET" {
		o := orm.NewOrm()
		var room models.Room
		id := c.Ctx.Input.Param(":id")
		o.QueryTable("b_room").Filter("room_id", id).One(&room)
		c.Data["RoomInfo"] = room
		c.TplName = "RoomUpdate.html"
	} else {
		o := orm.NewOrm()
		id, _ := c.GetInt("room_id")
		room := models.Room{RoomId: id}
		if o.Read(&room) == nil {
			if temp := c.GetString("name"); temp != "" {
				room.Name = temp
			}
			if temp := c.GetString("start_time"); temp != "" {
				room.StartTime = temp
			}
			if temp := c.GetString("end_time"); temp != "" {
				room.EndTime = temp
			}
			if temp, _ := c.GetInt("pos_num"); temp != 0 {
				room.PosNum = temp
			}
			if temp := c.GetString("desc"); temp != "" {
				room.Desc = temp
			}
			if temp, _ := c.GetInt("status"); temp != 0 {
				room.Status = temp
			}
			if num, err := o.Update(&room); err == nil && num == 1 {
				c.AjaxOk("modify successed")
			}
		}
		c.AjaxErr("modify failed")
	}
}
