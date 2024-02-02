package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/udistrital/sga_mid_proyecto_curricular/helpers"
	"github.com/udistrital/sga_mid_proyecto_curricular/models"
	"github.com/udistrital/sga_mid_proyecto_curricular/services"
)

type ConsultaProyectoAcademicoController struct {
	beego.Controller
}

// URLMapping ...
func (c *ConsultaProyectoAcademicoController) URLMapping() {
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("GetOnePorId", c.GetOnePorId)
	c.Mapping("Put", c.PutInhabilitarProyecto)
	c.Mapping("GetOneRegistroPorId", c.GetOneRegistroPorId)
}

// GetAll ...
// @Title GetAll
// @Description get ConsultaProyectoAcademico
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.ConsultaProyectoAcademico
// @Failure 403
// @router / [get]
func (c *ConsultaProyectoAcademicoController) GetAll() {
	var resultado map[string]interface{}
	var alerta models.Alert
	alertas := append([]interface{}{"Response:"})

	if resultado["Type"] != "error" {
		c.Data["json"] = services.PeticionProyectos(&alerta, &alertas)
	} else {
		if resultado["Body"] == "<QuerySeter> no row found" {
			c.Data["json"] = nil
		} else {
			helpers.ManejoError(&alerta, &alertas, fmt.Sprintf("%v", resultado["Body"]))
			c.Data["json"] = alerta
		}
	}
	c.ServeJSON()
}

// GetOnePorId ...
// @Title GetOnePorId
// @Description get ConsultaProyectoAcademico by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.ConsultaProyectoAcademico
// @Failure 403 :id is empty
// @router /:id [get]
func (c *ConsultaProyectoAcademicoController) GetOnePorId() {
	var resultado map[string]interface{}
	var alerta models.Alert
	alertas := append([]interface{}{"Response:"})
	idStr := c.Ctx.Input.Param(":id")

	if resultado["Type"] != "error" {
		c.Data["json"] = services.PeticionProyectosGetOneId(idStr, &alerta, &alertas)
	} else {
		if resultado["Body"] == "<QuerySeter> no row found" {
			c.Data["json"] = nil
		} else {
			helpers.ManejoError(&alerta, &alertas, fmt.Sprintf("%v", resultado["Body"]))
			c.Data["json"] = alerta
		}
	}
	c.ServeJSON()
}

// PutInhabilitarProyecto ...
// @Title PutInhabilitarProyecto
// @Description Inhabilitar Proyecto
// @Param	id		path 	string	true		"el id del proyecto a inhabilitar"
// @Param   body        body    {}  true        "body Inhabilitar Proyecto content"
// @Success 200 {}
// @Failure 403 :id is empty
// @router /inhabilitar_proyecto/:id [put]
func (c *ConsultaProyectoAcademicoController) PutInhabilitarProyecto() {
	idStr := c.Ctx.Input.Param(":id")
	var ProyectoAcademico map[string]interface{}
	var alerta models.Alert
	alertas := append([]interface{}{"Response:"})
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &ProyectoAcademico); err == nil {
		services.InhabilitarProyecto(&alerta, &alertas, idStr, ProyectoAcademico)
	} else {
		helpers.ManejoError(&alerta, &alertas, "", err)
	}
	c.Data["json"] = alerta
	c.ServeJSON()
}

// GetOneRegistroPorId ...
// @Title GetOneRegistroPorId
// @Description get ConsultaRegistro by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.ConsultaProyectoAcademico
// @Failure 403 :id is empty
// @router /get_registro/:id [get]
func (c *ConsultaProyectoAcademicoController) GetOneRegistroPorId() {
	var resultado map[string]interface{}
	var alerta models.Alert
	alertas := append([]interface{}{"Response:"})
	idStr := c.Ctx.Input.Param(":id")

	if resultado["Type"] != "error" {
		c.Data["json"] = services.PeticionRegistrosGetRegistroId(idStr, &alerta, &alertas)
	} else {
		if resultado["Body"] == "<QuerySeter> no row found" {
			c.Data["json"] = nil
		} else {
			helpers.ManejoError(&alerta, &alertas, fmt.Sprintf("%v", resultado["Body"]))
			c.Data["json"] = alerta
		}
	}
	c.ServeJSON()
}
