package controllers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/sga_proyecto_curricular_mid/services"
	"github.com/udistrital/utils_oas/errorhandler"
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

	data := c.Ctx.Input.RequestBody

	respuesta := services.ManejoPeticionesProyecto(data)

	c.Ctx.Output.SetStatus(respuesta.Status)
	c.Data["json"] = respuesta
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

	data := c.Ctx.Input.RequestBody

	respuesta := services.CoordinaById(data)

	c.Ctx.Output.SetStatus(respuesta.Status)
	c.Data["json"] = respuesta
	c.ServeJSON()
}
