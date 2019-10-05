package controllers

import (
	"venueleasing/models"

	"github.com/astaxie/beego/orm"
)

type AuthAdminController struct {
	Base
}

// @router /auth/login [get,post]
func (c *AuthAdminController) Login() {
	if c.Ctx.Request.Method == "GET" {
		c.TplName = "login.html"
	} else {
		c.SetSession("admin_auth", 1)
		o := orm.NewOrm()
		var user models.User
		username := c.GetString("username")
		o.QueryTable("b_user").Filter("username", username).One(&user)
		if user.Password == str2md5(c.GetString("password")+user.Salt) {
			c.AjaxOk("login successes")
		} else {
			c.AjaxErr("error username or password")
		}
	}
}

// @router /auth/logout [get]
func (c *AuthAdminController) Logout() {
	c.DelSession("admin_auth")
	c.StopRun()
}
