package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/udistrital/sga_mid_proyecto_curricular/helpers"
	"github.com/udistrital/sga_mid_proyecto_curricular/models"
	"github.com/udistrital/sga_mid_proyecto_curricular/services"
	"github.com/udistrital/utils_oas/errorhandler"
	"github.com/udistrital/utils_oas/request"
)

type CrearProyectoAcademicoController struct {
	beego.Controller
}

// URLMapping ...
func (c *CrearProyectoAcademicoController) URLMapping() {
	c.Mapping("PostProyecto", c.PostProyecto)
	c.Mapping("PostCoordinadorById", c.PostCoordinadorById)
}

// PostProyecto ...
// @Title PostProyecto
// @Description Crear Proyecto
// @Param   body        body    {}  true        "body Agregar Proyecto content"
// @Success 200 {}
// @Failure 403 body is empty
// @router / [post]
func (c *CrearProyectoAcademicoController) PostProyecto() {
	defer errorhandler.HandlePanic(&c.Controller)

	var Proyecto_academico map[string]interface{}
	var alerta models.Alert
	alertas := append([]interface{}{"Response:"})
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &Proyecto_academico); err == nil {
		if !services.ManejoPeticionesProyecto(&Proyecto_academico, &alerta, &alertas) {
			c.Ctx.Output.SetStatus(404)
			c.Data["json"] = alerta
			c.ServeJSON()
		} else {
			c.Ctx.Output.SetStatus(200)
		}
	} else {
		helpers.ManejoError(&alerta, &alertas, "", err)
	}

	c.Data["json"] = alerta
	c.ServeJSON()
}

// PostCoordinadorById ...
// @Title PostCoordinadorById
// @Description Post a de un cordinador de un proyecto existente, cambia estado activo a false a los coordinadores anteriores y crea el nuevo
// @Param	id		path 	string	true		"The key for staticblock"
// @Param   body        body    {}  true        "body Agregar Registro content"
// @Success 200 {object} models.ConsultaProyectoAcademico
// @Failure 403 :id is empty
// @router /coordinador [post]
func (c *CrearProyectoAcademicoController) PostCoordinadorById() {
	defer errorhandler.HandlePanic(&c.Controller)

	var CoordinadorNuevo map[string]interface{}
	var resultado map[string]interface{}
	var alerta models.Alert

	alertas := []interface{}{"Response:"}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &CoordinadorNuevo); err == nil {
		if resultado["Type"] != "error" {
			var CoordinadorAntiguos []map[string]interface{}
			idStr := fmt.Sprintf("%v", CoordinadorNuevo["ProyectoAcademicoInstitucionId"].(map[string]interface{})["Id"])

			errcordinador := request.GetJson("http://"+beego.AppConfig.String("ProyectoAcademicoService")+"/proyecto_academico_rol_tercero_dependencia/?query=ProyectoAcademicoInstitucionId.Id:"+idStr, &CoordinadorAntiguos)
			if errcordinador == nil {
				if CoordinadorAntiguos[0]["Id"] != nil {
					if exito := services.ManejoCoordinadorAntiguo(&alertas, &alerta, CoordinadorAntiguos); !exito {
						c.Ctx.Output.SetStatus(400)
					}
					if exito := services.RegistrarCoordinador(&alertas, &alerta, CoordinadorNuevo); !exito {
						c.Ctx.Output.SetStatus(400)
					}

					c.Ctx.Output.SetStatus(200)
					c.Data["json"] = alerta
					c.ServeJSON()
				} else {
					if err := json.Unmarshal(c.Ctx.Input.RequestBody, &CoordinadorNuevo); err == nil {
						if exito := services.RegistrarCoordinador(&alertas, &alerta, CoordinadorNuevo); !exito {
							c.Ctx.Output.SetStatus(404)
						} else {
							c.Ctx.Output.SetStatus(200)
						}
					} else {
						c.Ctx.Output.SetStatus(400)
						helpers.ManejoError(&alerta, &alertas, "", err)
					}
				}
			} else {
				c.Ctx.Output.SetStatus(400)
				helpers.ManejoError(&alerta, &alertas, "", errcordinador)
				c.Data["json"] = alerta
			}
		} else {
			if resultado["Body"] == "<QuerySeter> no row found" {
				c.Ctx.Output.SetStatus(404)
				c.Data["json"] = nil
			} else {
				c.Ctx.Output.SetStatus(400)
				helpers.ManejoError(&alerta, &alertas, fmt.Sprintf("%v", resultado["Body"]))
				c.Data["json"] = alerta
			}
		}
	} else {
		c.Ctx.Output.SetStatus(400)
		helpers.ManejoError(&alerta, &alertas, "", err)
	}
	//alerta.Body = alertas
	c.Data["json"] = alerta
	c.ServeJSON()
}
