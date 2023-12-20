// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"sga_mid_proyecto_curricular/controllers"

	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/consulta_proyecto_academico",
			beego.NSInclude(
				&controllers.ConsultaProyectoAcademicoController{},
			),
		),
		beego.NSNamespace("/proyecto_academico",
			beego.NSInclude(
				&controllers.CrearProyectoAcademicoController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
