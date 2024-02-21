package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/udistrital/sga_mid_proyecto_curricular/helpers"
	"github.com/udistrital/sga_mid_proyecto_curricular/models"
	"github.com/udistrital/sga_mid_proyecto_curricular/services"
	"github.com/udistrital/utils_oas/errorhandler"
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
	defer errorhandler.HandlePanic(&c.Controller)

	var resultado map[string]interface{}
	var alerta models.Alert
	alertas := append([]interface{}{"Response:"})

	if resultado["Type"] != "error" {
		if respuesta, exito := services.PeticionProyectos(&alerta, &alertas); !exito {
			c.Ctx.Output.SetStatus(404)
			c.Data["json"] = respuesta
		} else {
			c.Ctx.Output.SetStatus(200)
			c.Data["json"] = respuesta
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
	defer errorhandler.HandlePanic(&c.Controller)

	var resultado map[string]interface{}
	var alerta models.Alert
	alertas := append([]interface{}{"Response:"})
	idStr := c.Ctx.Input.Param(":id")

	if resultado["Type"] != "error" {
		if respuesta, exito := services.PeticionProyectosGetOneId(idStr, &alerta, &alertas); !exito {
			c.Ctx.Output.SetStatus(404)
			c.Data["json"] = respuesta
		} else {
			c.Ctx.Output.SetStatus(200)
			c.Data["json"] = respuesta
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
	defer errorhandler.HandlePanic(&c.Controller)

	idStr := c.Ctx.Input.Param(":id")
	var ProyectoAcademico map[string]interface{}
	var alerta models.Alert
	alertas := append([]interface{}{"Response:"})
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &ProyectoAcademico); err == nil {
		if exito := services.InhabilitarProyecto(&alerta, &alertas, idStr, ProyectoAcademico); !exito {
			c.Ctx.Output.SetStatus(404)
		} else {
			c.Ctx.Output.SetStatus(200)
		}
	} else {
		c.Ctx.Output.SetStatus(400)
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
	defer errorhandler.HandlePanic(&c.Controller)

	var resultado map[string]interface{}
	var alerta models.Alert
	alertas := append([]interface{}{"Response:"})
	idStr := c.Ctx.Input.Param(":id")

	if resultado["Type"] != "error" {
		if respuesta, exito := services.PeticionRegistrosGetRegistroId(idStr, &alerta, &alertas); !exito {
			c.Ctx.Output.SetStatus(404)
			c.Data["json"] = respuesta
		} else {
			c.Ctx.Output.SetStatus(200)
			c.Data["json"] = respuesta
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
	c.ServeJSON()
}
