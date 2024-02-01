package controllers

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/astaxie/beego"
	"github.com/udistrital/sga_mid_proyecto_curricular/models"
	"github.com/udistrital/sga_mid_proyecto_curricular/services"
	"github.com/udistrital/utils_oas/request"
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
		var proyectos []map[string]interface{}
		errproyecto := request.GetJson("http://"+beego.AppConfig.String("ProyectoAcademicoService")+"/tr_proyecto_academico/", &proyectos)

		if errproyecto == nil {
			services.ManejoProyectos(&proyectos)

			c.Data["json"] = proyectos
		} else {
			services.ManejoError(&alerta, &alertas, "", errproyecto)
			c.Data["json"] = alerta
		}
	} else {
		if resultado["Body"] == "<QuerySeter> no row found" {
			c.Data["json"] = nil
		} else {
			services.ManejoError(&alerta, &alertas, fmt.Sprintf("%v", resultado["Body"]))
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
		// var idOikos float64
		var idUnidad float64
		var proyectos []map[string]interface{}
		// var dependencias []map[string]interface{}
		var unidades []map[string]interface{}

		//errproyecto := request.GetJson("http://"+beego.AppConfig.String("ProyectoAcademicoService")+"/tr_proyecto_academico/"+idStr, &proyectos)
		fmt.Println("URL", "http://"+beego.AppConfig.String("ProyectoAcademicoService")+"/tr_proyecto_academico/"+idStr)
		fmt.Println("VARIABLE", beego.AppConfig.String("ProyectoAcademicoService"))
		errproyecto := request.GetJson("http://"+beego.AppConfig.String("ProyectoAcademicoService")+"/tr_proyecto_academico/"+idStr, &proyectos)
		// errdependencia := request.GetJson("http://"+beego.AppConfig.String("OikosService")+"/dependencia_tipo_dependencia/?query=TipoDependenciaId:2", &dependencias)
		errunidad := request.GetJson("http://"+beego.AppConfig.String("CoreService")+"/unidad_tiempo/", &unidades)

		if proyectos[0]["ProyectoAcademico"] != nil {

			// if errproyecto == nil && errdependencia == nil && errunidad == nil {
			if errproyecto == nil && errunidad == nil {

				for _, proyecto := range proyectos {
					registros := proyecto["Registro"].([]interface{})
					proyectobase := proyecto["ProyectoAcademico"].(map[string]interface{})
					proyecto["FechaVenimientoAcreditacion"] = nil
					proyecto["FechaVenimientoCalidad"] = nil
					proyecto["TieneRegistroAltaCalidad"] = false
					proyecto["NumeroActoAdministrativoAltaCalidad"] = nil
					proyecto["AnoActoAdministrativoIdAltaCalidad"] = nil
					proyecto["FechaCreacionActoAdministrativoAltaCalidad"] = nil
					proyecto["VigenciaActoAdministrativoAltaCalidad"] = nil
					proyecto["EnlaceActoAdministrativoAltaCalidad"] = nil

					/*
						for _, dependencia := range dependencias {
							proyectotem := dependencia["DependenciaId"].(map[string]interface{})
							idOikos = proyectotem["Id"].(float64)
							if proyectobase["DependenciaId"].(float64) == idOikos {
								proyecto["NombreDependencia"] = proyectotem["Nombre"]
								proyecto["IdDependenciaFacultad"] = proyectotem["Id"]
								proyecto["TelefonoDependencia"] = proyectotem["TelefonoDependencia"]
							}
						}
					*/

					// Información de la facultad
					var dependenciaFacultad map[string]interface{}
					errdependenciaFacultad := request.GetJson("http://"+beego.AppConfig.String("OikosService")+"/dependencia/"+fmt.Sprintf("%.f", proyectobase["FacultadId"].(float64)), &dependenciaFacultad)
					// if errdependencia["Type"] == "error" || errdependencia != nil || dependencia["Status"] == "404" || dependencia["Message"] != nil {
					if errdependenciaFacultad == nil {
						proyecto["NombreFacultad"] = dependenciaFacultad["Nombre"]
						proyecto["IdDependenciaFacultad"] = dependenciaFacultad["Id"]
					}

					// Información de la dependencia del proyecto
					var dependencia map[string]interface{}
					errdependencia := request.GetJson("http://"+beego.AppConfig.String("OikosService")+"/dependencia/"+fmt.Sprintf("%.f", proyectobase["DependenciaId"].(float64)), &dependencia)
					// if errdependencia["Type"] == "error" || errdependencia != nil || dependencia["Status"] == "404" || dependencia["Message"] != nil {
					if errdependencia == nil {
						proyecto["TelefonoDependencia"] = dependencia["TelefonoDependencia"]
					}

					if proyectobase["Oferta"] == true {
						proyecto["OfertaLetra"] = "Si"
					} else if proyectobase["Oferta"] == false {
						proyecto["OfertaLetra"] = "No"
					}
					if proyectobase["CiclosPropedeuticos"] == true {
						proyecto["CiclosLetra"] = "Si"
					} else if proyectobase["CiclosPropedeuticos"] == false {
						proyecto["CiclosLetra"] = "NO"
					}

					for _, unidad := range unidades {
						unidadTem := unidad
						idUnidad = unidadTem["Id"].(float64)
						if proyectobase["UnidadTiempoId"].(float64) == idUnidad {
							proyecto["NombreUnidad"] = unidadTem["Nombre"]
						}

					}

					for _, registrotemp := range registros {
						registro := registrotemp.(map[string]interface{})

						tiporegistro := registro["TipoRegistroId"].(map[string]interface{})

						if tiporegistro["Id"].(float64) == 1 && registro["Activo"] == true {
							proyecto["FechaVenimientoAcreditacion"] = registro["VencimientoActoAdministrativo"]
							proyecto["FechaVenimientoCalidad"] = "00/00/0000"
						} else if tiporegistro["Id"].(float64) == 2 && registro["Activo"] == true {

							proyecto["FechaVenimientoCalidad"] = registro["VencimientoActoAdministrativo"]
							proyecto["TieneRegistroAltaCalidad"] = true
							proyecto["NumeroActoAdministrativoAltaCalidad"] = registro["NumeroActoAdministrativo"]
							proyecto["AnoActoAdministrativoIdAltaCalidad"] = registro["AnoActoAdministrativoId"]
							proyecto["FechaCreacionActoAdministrativoAltaCalidad"] = registro["FechaCreacionActoAdministrativo"]
							proyecto["VigenciaActoAdministrativoAltaCalidad"] = registro["VigenciaActoAdministrativo"]
							proyecto["EnlaceActoAdministrativoAltaCalidad"] = registro["EnlaceActo"]
						}
					}

				}

				c.Data["json"] = proyectos

			} else {
				alertas = append(alertas, errproyecto.Error())
				alerta.Code = "400"
				alerta.Type = "error"
				alerta.Body = alertas
				c.Data["json"] = alerta
			}

		} else {
			c.Data["json"] = proyectos
		}
	} else {
		if resultado["Body"] == "<QuerySeter> no row found" {
			c.Data["json"] = nil
		} else {
			alertas = append(alertas, resultado["Body"])
			alerta.Code = "400"
			alerta.Type = "error"
			alerta.Body = alertas
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

		var resultadoProyecto map[string]interface{}
		errProyecto := request.SendJson("http://"+beego.AppConfig.String("ProyectoAcademicoService")+"/proyecto_academico_institucion/"+idStr, "PUT", &resultadoProyecto, ProyectoAcademico)
		if resultadoProyecto["Type"] == "error" || errProyecto != nil || resultadoProyecto["Status"] == "404" || resultadoProyecto["Message"] != nil {
			alertas = append(alertas, resultadoProyecto)
			alerta.Type = "error"
			alerta.Code = "400"
		} else {
			alertas = append(alertas, ProyectoAcademico)
		}

	} else {
		alerta.Type = "error"
		alerta.Code = "400"
		alertas = append(alertas, err.Error())
	}
	alerta.Body = alertas
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
		var registros []map[string]interface{}

		errproyecto := request.GetJson("http://"+beego.AppConfig.String("ProyectoAcademicoService")+"/registro_calificado_acreditacion/?query=ProyectoAcademicoInstitucionId.Id:"+idStr, &registros)

		if errproyecto == nil {
			if registros[0]["Id"] != nil {
				for _, registro := range registros {
					vigenciatemporal := registro["VigenciaActoAdministrativo"].(string)
					vigenciatemporal = strings.Replace(vigenciatemporal, "A", " A", 1)
					registro["VigenciaActoAdministrativo"] = vigenciatemporal
					if registro["Activo"] == true {
						registro["ActivoLetra"] = "Si"

					} else if registro["Activo"] == false {
						registro["ActivoLetra"] = "No"
					}
				}
			}

			c.Data["json"] = registros

		} else {
			alertas = append(alertas, errproyecto.Error())
			alerta.Code = "400"
			alerta.Type = "error"
			alerta.Body = alertas
			c.Data["json"] = alerta
		}

	} else {
		if resultado["Body"] == "<QuerySeter> no row found" {
			c.Data["json"] = nil
		} else {
			alertas = append(alertas, resultado["Body"])
			alerta.Code = "400"
			alerta.Type = "error"
			alerta.Body = alertas
			c.Data["json"] = alerta
		}
	}
	c.ServeJSON()
}
