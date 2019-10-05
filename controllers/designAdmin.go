package controllers

import (
	"crypto/md5"
	"fmt"
	"io"
	"venueleasing/models"

	"github.com/astaxie/beego/orm"
)

func str2md5(str string) string {
	w := md5.New()
	//write str to w
	io.WriteString(w, str)
	//convert w's hash to []byte by w.Sum(nil)
	return fmt.Sprintf("%x", w.Sum(nil))
}

type DesignAdminController struct {
	BaseAdmin
}

//exhibite design list
// @router /design/show [get]
func (c *DesignAdminController) Show() {
	o := orm.NewOrm()
	var design []*models.Design
	o.QueryTable("b_design").All(&design)
	c.Data["Design"] = design
	c.TplName = "DesignList.html"
}

//add design
// @router /design/add [get,post]
func (c *DesignAdminController) Add() {
	if c.Ctx.Request.Method == "GET" {
		c.TplName = "DesignAdd.html"
	} else {
		o := orm.NewOrm()
		var design models.Design
		design.Name = c.GetString("name")
		design.Desc = c.GetString("desc")
		design.Status = 1
		design.Key = str2md5(design.Name)
		_, err := o.Insert(&design)
		if err != nil {
			c.AjaxErr("add failed")
		} else {
			c.AjaxOk("add successed")
		}
	}
}

//remove design
// @router /design/delete [post]
func (c *DesignAdminController) Delete() {
	o := orm.NewOrm()
	id, err := c.GetInt("design_id")
	if err == nil {
		if num, err := o.Delete(&models.Design{DesignId: id}); err == nil && num == 1 {
			c.AjaxOk("removed")
		}
	}
	c.AjaxErr("failed")
}

//modify design
// @router /design/update/?:id [get,post]
func (c *DesignAdminController) Update() {
	if c.Ctx.Request.Method == "GET" {
		o := orm.NewOrm()
		var design models.Design
		id := c.Ctx.Input.Param(":id")
		o.QueryTable("b_design").Filter("design_id", id).One(&design)
		c.Data["DesignInfo"] = design
		c.TplName = "DesignUpdate.html"
	} else {
		o := orm.NewOrm()
		id, _ := c.GetInt("design_id")
		design := models.Design{DesignId: id}
		if o.Read(&design) == nil {
			if temp := c.GetString("name"); temp != "" {
				design.Name = temp
				design.Key = str2md5(temp)
			}
			if temp := c.GetString("desc"); temp != "" {
				design.Desc = temp
			}
			if temp, _ := c.GetInt("status"); temp != 0 {
				design.Status = temp
			}
			if num, err := o.Update(&design); err == nil && num == 1 {
				c.AjaxOk("successed")
			}
		}
		c.AjaxErr("failed")
	}
}
