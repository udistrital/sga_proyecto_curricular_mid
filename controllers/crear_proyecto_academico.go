package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/time_bogota"
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

		Proyecto_academicoPost := make(map[string]interface{})
		Proyecto_academicoPost = map[string]interface{}{
			"ProyectoAcademicoInstitucion": Proyecto_academico["ProyectoAcademicoInstitucion"],
			"Enfasis":                      Proyecto_academico["Enfasis"],
			"Registro":                     Proyecto_academico["Registro"],
			"Titulaciones":                 Proyecto_academico["Titulaciones"],
		}

		Proyecto_academico_oikosPost := Proyecto_academico["Oikos"]

		var resultadoOikos map[string]interface{}
		var resultadoProyecto map[string]interface{}

		errOikos := request.SendJson("http://"+beego.AppConfig.String("OikosService")+"/dependencia_padre/tr_dependencia_padre", "POST", &resultadoOikos, Proyecto_academico_oikosPost)
		if resultadoOikos["Type"] == "error" || errOikos != nil || resultadoOikos["Status"] == "404" || resultadoOikos["Message"] != nil {

			errorMessage := "Error in GetOnePorId: "
			if errOikos != nil {
				errorMessage += "Oikos: " + errOikos.Error() + "; "
			} else {
				errorMessage += "Resultado Oikos: " + resultadoOikos["Message"].(string) + "; "
			}

			logs.Error(errorMessage)
			c.Data["message"] = "Error service PostProyecto: " + errorMessage
			c.Abort("404")
		} else {
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Request successful", "Data": Proyecto_academico}
			idDependenciaProyecto := resultadoOikos["HijaId"].(map[string]interface{})["Id"]
			Proyecto_academicoPost["ProyectoAcademicoInstitucion"].(map[string]interface{})["DependenciaId"] = idDependenciaProyecto
		}

		errProyecto := request.SendJson("http://"+beego.AppConfig.String("ProyectoAcademicoService")+"/tr_proyecto_academico", "POST", &resultadoProyecto, Proyecto_academicoPost)
		if resultadoProyecto["Type"] == "error" || errProyecto != nil || resultadoProyecto["Status"] == "404" || resultadoProyecto["Message"] != nil {
			logs.Error(errProyecto)
			c.Data["message"] = "Error service PostProyecto: " + errProyecto
			c.Abort("404")
		} else {
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Request successful", "Data": Proyecto_academico}
		}

	} else {
		logs.Error(err)
		c.Data["message"] = "Error service PostProyecto: " + err.Error()
		c.Abort("404")
	}

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

					for _, cordinadorFecha := range CoordinadorAntiguos {
						if cordinadorFecha["Activo"] == true {
							cordinadorFecha["Activo"] = false
							coordinador_cambiado := cordinadorFecha
							coordinador_cambiado["FechaFinalizacion"] = time_bogota.Tiempo_bogota()
							Id_coordinador_cambiado := cordinadorFecha["Id"]
							idcoordinador := Id_coordinador_cambiado.(float64)
							var resultado map[string]interface{}
							errcoordinadorcambiado := request.SendJson("http://"+beego.AppConfig.String("ProyectoAcademicoService")+"/proyecto_academico_rol_tercero_dependencia/"+strconv.FormatFloat(idcoordinador, 'f', -1, 64), "PUT", &resultado, &coordinador_cambiado)
							if resultado["Type"] == "error" || errcoordinadorcambiado != nil || resultado["Status"] == "404" || resultado["Message"] != nil {
								logs.Error(resultado)
								c.Data["message"] = "Error service PostCoordinadorById: " + resultado["Message"].(string)
								c.Abort("404")
							}
						} else {
							fmt.Println("Todos los registros estan nulos")
						}

					}

					var resultadoCoordinadorNuevo map[string]interface{}
					CoordinadorNuevo["FechaFinalizacion"] = "0001-01-01T00:00:00-05:00"
					errRegistro := request.SendJson("http://"+beego.AppConfig.String("ProyectoAcademicoService")+"/proyecto_academico_rol_tercero_dependencia", "POST", &resultadoCoordinadorNuevo, CoordinadorNuevo)
					if resultadoCoordinadorNuevo["Type"] == "error" || errRegistro != nil || resultadoCoordinadorNuevo["Status"] == "404" || resultadoCoordinadorNuevo["Message"] != nil {
						logs.Error(resultadoCoordinadorNuevo)
						c.Data["message"] = "Error service PostCoordinadorById: " + resultadoCoordinadorNuevo["Message"].(string)
						c.Abort("404")
					} else {
						c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Request successful", "Data": CoordinadorNuevo}
					}

					c.ServeJSON()
				} else {
					if err := json.Unmarshal(c.Ctx.Input.RequestBody, &CoordinadorNuevo); err == nil {
						var resultadoCoordinadorNuevo map[string]interface{}
						CoordinadorNuevo["FechaFinalizacion"] = "0001-01-01T00:00:00-05:00"

						errRegistro := request.SendJson("http://"+beego.AppConfig.String("ProyectoAcademicoService")+"/proyecto_academico_rol_tercero_dependencia", "POST", &resultadoCoordinadorNuevo, CoordinadorNuevo)
						if resultadoCoordinadorNuevo["Type"] == "error" || errRegistro != nil || resultadoCoordinadorNuevo["Status"] == "404" || resultadoCoordinadorNuevo["Message"] != nil {
							logs.Error(resultadoCoordinadorNuevo)
							c.Data["message"] = "Error service PostCoordinadorById: " + resultadoCoordinadorNuevo["Message"].(string)
							c.Abort("404")

						} else {
							c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Request successful", "Data": CoordinadorNuevo}
						}

					} else {
						logs.Error(err)
						c.Data["message"] = "Error service PostCoordinadorById: " + err.Error()
						c.Abort("404")
					}

				}
			} else {
				logs.Error(errcordinador)
				c.Data["message"] = "Error service PostCoordinadorById: " + errcordinador.Error()
				c.Abort("404")
			}
		} else {
			errorMessage := resultado["Body"].(string)
			if errorMessage == "<QuerySeter> no row found" {
				c.Data["json"] = nil
			} else {
				logs.Error(resultado)
				c.Data["message"] = "Error service PostCoordinadorById: " + errorMessage
				c.Abort("404")
			}
		}
	} else {
		logs.Error(err)
		c.Data["message"] = "Error service PostCoordinadorById: " + err.Error()
		c.Abort("404")
	}

	c.ServeJSON()
}
