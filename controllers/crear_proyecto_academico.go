package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/udistrital/sga_mid_proyecto_curricular/helpers"
	"github.com/udistrital/sga_mid_proyecto_curricular/models"
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
// @Failure 403 body is empty
// @router / [post]
func (c *CrearProyectoAcademicoController) PostProyecto() {

	var Proyecto_academico map[string]interface{}
	var alerta models.Alert
	alertas := append([]interface{}{"Response:"})
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
			helpers.ManejoError(&alerta, &alertas, fmt.Sprintf("%v", resultadoOikos), errOikos)
			c.Data["json"] = alerta
			c.ServeJSON()
		} else {
			services.AsignarProyectoAcademico(&alertas, &Proyecto_academico, resultadoOikos, &Proyecto_academicoPost)
			/*alertas = append(alertas, Proyecto_academico)
			idDependenciaProyecto := resultadoOikos["HijaId"].(map[string]interface{})["Id"]
			Proyecto_academicoPost["ProyectoAcademicoInstitucion"].(map[string]interface{})["DependenciaId"] = idDependenciaProyecto*/
		}

		errProyecto := request.SendJson("http://"+beego.AppConfig.String("ProyectoAcademicoService")+"/tr_proyecto_academico", "POST", &resultadoProyecto, Proyecto_academicoPost)
		if resultadoProyecto["Type"] == "error" || errProyecto != nil || resultadoProyecto["Status"] == "404" || resultadoProyecto["Message"] != nil {
			helpers.ManejoError(&alerta, &alertas, "", errProyecto)
			c.Data["json"] = alerta
			c.ServeJSON()
		} else {
			alertas = append(alertas, Proyecto_academico)
		}

	} else {
		helpers.ManejoError(&alerta, &alertas, "", err)
	}

	//alerta.Body = alertas
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
					services.ManejoCoordinadorAntiguo(&alertas, &alerta, CoordinadorAntiguos)
					services.RegistrarCoordinador(&alertas, &alerta, CoordinadorNuevo)
					/*for _, cordinadorFecha := range CoordinadorAntiguos {
						if cordinadorFecha["Activo"] == true {
							cordinadorFecha["Activo"] = false
							coordinador_cambiado := cordinadorFecha
							coordinador_cambiado["FechaFinalizacion"] = time_bogota.Tiempo_bogota()
							Id_coordinador_cambiado := cordinadorFecha["Id"]
							idcoordinador := Id_coordinador_cambiado.(float64)
							var resultado map[string]interface{}
							errcoordinadorcambiado := request.SendJson("http://"+beego.AppConfig.String("ProyectoAcademicoService")+"/proyecto_academico_rol_tercero_dependencia/"+strconv.FormatFloat(idcoordinador, 'f', -1, 64), "PUT", &resultado, &coordinador_cambiado)
							if resultado["Type"] == "error" || errcoordinadorcambiado != nil || resultado["Status"] == "404" || resultado["Message"] != nil {
								alertas = append(alertas, resultado)
								alerta.Type = "error"
								alerta.Code = "400"
								alerta.Body = alertas
							}
						} else {
							fmt.Println("Todos los registros estan nulos")
						}

					}*/

					/*var resultadoCoordinadorNuevo map[string]interface{}
					CoordinadorNuevo["FechaFinalizacion"] = "0001-01-01T00:00:00-05:00"
					errRegistro := request.SendJson("http://"+beego.AppConfig.String("ProyectoAcademicoService")+"/proyecto_academico_rol_tercero_dependencia", "POST", &resultadoCoordinadorNuevo, CoordinadorNuevo)
					if resultadoCoordinadorNuevo["Type"] == "error" || errRegistro != nil || resultadoCoordinadorNuevo["Status"] == "404" || resultadoCoordinadorNuevo["Message"] != nil {
						alertas = append(alertas, resultadoCoordinadorNuevo)
						alerta.Type = "error"
						alerta.Code = "400"
						alerta.Body = alertas
					} else {
						alertas = append(alertas, CoordinadorNuevo)
					}*/
					//alerta.Body = alertas
					c.Data["json"] = alerta
					c.ServeJSON()
				} else {
					if err := json.Unmarshal(c.Ctx.Input.RequestBody, &CoordinadorNuevo); err == nil {
						services.RegistrarCoordinador(&alertas, &alerta, CoordinadorNuevo)
						/*var resultadoCoordinadorNuevo map[string]interface{}
						CoordinadorNuevo["FechaFinalizacion"] = "0001-01-01T00:00:00-05:00"

						errRegistro := request.SendJson("http://"+beego.AppConfig.String("ProyectoAcademicoService")+"/proyecto_academico_rol_tercero_dependencia", "POST", &resultadoCoordinadorNuevo, CoordinadorNuevo)
						if resultadoCoordinadorNuevo["Type"] == "error" || errRegistro != nil || resultadoCoordinadorNuevo["Status"] == "404" || resultadoCoordinadorNuevo["Message"] != nil {
							alertas = append(alertas, resultadoCoordinadorNuevo)
							alerta.Type = "error"
							alerta.Code = "400"
							alerta.Body = alertas
						} else {
							alertas = append(alertas, CoordinadorNuevo)
						}*/
					} else {
						helpers.ManejoError(&alerta, &alertas, "", err)
					}
				}
			} else {
				helpers.ManejoError(&alerta, &alertas, "", errcordinador)
				c.Data["json"] = alerta
			}
		} else {
			if resultado["Body"] == "<QuerySeter> no row found" {
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
