package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/udistrital/sga_mid_proyecto_curricular/helpers"
	"github.com/udistrital/sga_mid_proyecto_curricular/services"
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
// @Failure 400 the request contains incorrect syntax
// @router / [post]
func (c *CrearProyectoAcademicoController) PostProyecto() {

	var Proyecto_academico map[string]interface{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &Proyecto_academico); err == nil {
		if !services.ManejoPeticionesProyecto(&Proyecto_academico, &alerta, &alertas) {
			c.Data["json"] = alerta
			c.ServeJSON()
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
// @Failure 400 the request contains incorrect syntax
// @router /coordinador [post]
func (c *CrearProyectoAcademicoController) PostCoordinadorById() {
	var CoordinadorNuevo map[string]interface{}
	var resultado map[string]interface{}
	// var alerta models.Alert

	// alertas := []interface{}{"Response:"}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &CoordinadorNuevo); err == nil {
		if resultado["Type"] != "error" {
			var CoordinadorAntiguos []map[string]interface{}
			idStr := fmt.Sprintf("%v", CoordinadorNuevo["ProyectoAcademicoInstitucionId"].(map[string]interface{})["Id"])

			errcordinador := request.GetJson("http://"+beego.AppConfig.String("ProyectoAcademicoService")+"/proyecto_academico_rol_tercero_dependencia/?query=ProyectoAcademicoInstitucionId.Id:"+idStr, &CoordinadorAntiguos)
			if errcordinador == nil {
				if CoordinadorAntiguos[0]["Id"] != nil {
					services.ManejoCoordinadorAntiguo(&alertas, &alerta, CoordinadorAntiguos)
					services.RegistrarCoordinador(&alertas, &alerta, CoordinadorNuevo)
					c.Data["json"] = alerta
					c.ServeJSON()
				} else {
					if err := json.Unmarshal(c.Ctx.Input.RequestBody, &CoordinadorNuevo); err == nil {
						services.RegistrarCoordinador(&alertas, &alerta, CoordinadorNuevo)
					} else {
						helpers.ManejoError(&alerta, &alertas, "", err)
					}
				}
			} else {
				helpers.ManejoError(&alerta, &alertas, "", errcordinador)
				c.Data["json"] = alerta
			}
		} else {
			errorMessage := resultado["Body"].(string)
			if errorMessage == "<QuerySeter> no row found" {
				c.Data["json"] = nil
			} else {
				helpers.ManejoError(&alerta, &alertas, fmt.Sprintf("%v", resultado["Body"]))
				c.Data["json"] = alerta
			}
		}
	} else {
		helpers.ManejoError(&alerta, &alertas, "", err)
	}
	//alerta.Body = alertas
	c.Data["json"] = alerta
	c.ServeJSON()
}
