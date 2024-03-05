// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"github.com/udistrital/sga_proyecto_curricular_mid/controllers"
	"github.com/udistrital/utils_oas/errorhandler"

	"github.com/astaxie/beego"
)

func init() {

	beego.ErrorController(&errorhandler.ErrorHandlerController{})

	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/proyecto-academico",
			beego.NSInclude(
				&controllers.ProyectoAcademicoController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
