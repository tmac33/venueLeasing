package routers

import (
	"venueleasing/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.BookingController{}, "*:Show")
}
