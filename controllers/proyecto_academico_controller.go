package controllers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/sga_proyecto_curricular_mid/services"
	"github.com/udistrital/utils_oas/errorhandler"
)

type ProyectoAcademicoController struct {
	beego.Controller
}

// URLMapping ...
func (c *ProyectoAcademicoController) URLMapping() {
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("GetOnePorId", c.GetOnePorId)
	c.Mapping("Put", c.PutInhabilitarProyecto)
	c.Mapping("GetOneRegistroPorId", c.GetOneRegistroPorId)
	c.Mapping("PostProyecto", c.PostProyecto)
	c.Mapping("PostCoordinadorById", c.PostCoordinadorById)
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
// @Failure 404 not found resource
// @router / [get]
func (c *ProyectoAcademicoController) GetAll() {
	defer errorhandler.HandlePanic(&c.Controller)

	respuesta := services.PeticionProyectos()

	c.Ctx.Output.SetStatus(respuesta.Status)
	c.Data["json"] = respuesta
	c.ServeJSON()
}

// GetOnePorId ...
// @Title GetOnePorId
// @Description get ConsultaProyectoAcademico by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.ConsultaProyectoAcademico
// @Failure 404 not found resource
// @router /:id [get]
func (c *ProyectoAcademicoController) GetOnePorId() {
	defer errorhandler.HandlePanic(&c.Controller)

	idStr := c.Ctx.Input.Param(":id")

	respuesta := services.PeticionProyectosGetOneId(idStr)
	c.Ctx.Output.SetStatus(respuesta.Status)
	c.Data["json"] = respuesta
	c.ServeJSON()
}

// PutInhabilitarProyecto ...
// @Title PutInhabilitarProyecto
// @Description Inhabilitar Proyecto
// @Param	id		path 	string	true		"el id del proyecto a inhabilitar"
// @Param   body        body    {}  true        "body Inhabilitar Proyecto content"
// @Success 200 {}
// @Failure 400 the request contains incorrect syntax
// @router /:id/inhabilitar [put]
func (c *ProyectoAcademicoController) PutInhabilitarProyecto() {
	defer errorhandler.HandlePanic(&c.Controller)

	idStr := c.Ctx.Input.Param(":id")
	data := c.Ctx.Input.RequestBody

	respuesta := services.InhabilitarProyecto(idStr, data)

	c.Ctx.Output.SetStatus(respuesta.Status)
	c.Data["json"] = respuesta
	c.ServeJSON()
}

// GetOneRegistroPorId ...
// @Title GetOneRegistroPorId
// @Description get ConsultaRegistro by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.ConsultaProyectoAcademico
// @Failure 403 :id is empty
// @router /registro/:id [get]
func (c *ProyectoAcademicoController) GetOneRegistroPorId() {
	defer errorhandler.HandlePanic(&c.Controller)

	idStr := c.Ctx.Input.Param(":id")

	repuesta := services.PeticionRegistrosGetRegistroId(idStr)
	c.Ctx.Output.SetStatus(repuesta.Status)
	c.Data["json"] = repuesta
	c.ServeJSON()
}

// PostProyecto ...
// @Title PostProyecto
// @Description Crear Proyecto
// @Param   body        body    {}  true        "body Agregar Proyecto content"
// @Success 200 {}
// @Failure 403 body is empty
// @router / [post]
func (c *ProyectoAcademicoController) PostProyecto() {
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
func (c *ProyectoAcademicoController) PostCoordinadorById() {
	defer errorhandler.HandlePanic(&c.Controller)

	data := c.Ctx.Input.RequestBody

	respuesta := services.CoordinaById(data)

	c.Ctx.Output.SetStatus(respuesta.Status)
	c.Data["json"] = respuesta
	c.ServeJSON()
}
