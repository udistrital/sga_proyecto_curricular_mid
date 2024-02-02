package services

import (
	"fmt"
	"strings"

	"github.com/astaxie/beego"
	"github.com/udistrital/sga_mid_proyecto_curricular/helpers"
	"github.com/udistrital/sga_mid_proyecto_curricular/models"
	"github.com/udistrital/utils_oas/request"
)

// FUNCIONES QUE SE USAN EN GETALL

func PeticionProyectos(alerta *models.Alert, alertas *[]interface{}) interface{} {
	var proyectos []map[string]interface{}
	errproyecto := request.GetJson("http://"+beego.AppConfig.String("ProyectoAcademicoService")+"/tr_proyecto_academico/", &proyectos)

	if errproyecto == nil {
		ManejoProyectosGetAll(&proyectos)
		return proyectos
	} else {
		helpers.ManejoError(alerta, alertas, "", errproyecto)
		return *alerta
	}
}

func ManejoProyectosGetAll(proyectos *[]map[string]interface{}) {
	for _, proyecto := range *proyectos {
		registros := proyecto["Registro"].([]interface{})
		proyectobase := proyecto["ProyectoAcademico"].(map[string]interface{})
		proyecto["FechaVenimientoAcreditacion"] = nil
		proyecto["FechaVenimientoCalidad"] = nil

		// Información de la facultad
		var dependencia map[string]interface{}
		errdependencia := request.GetJson("http://"+beego.AppConfig.String("OikosService")+"/dependencia/"+fmt.Sprintf("%.f", proyectobase["FacultadId"].(float64)), &dependencia)
		if errdependencia == nil {
			proyecto["NombreFacultad"] = dependencia["Nombre"]
		}

		if proyectobase["Oferta"] == true {
			proyecto["OfertaLetra"] = "Si"
		} else if proyectobase["Oferta"] == false {
			proyecto["OfertaLetra"] = "No"
		}

		ManejoRegistrosGetAll(registros, proyecto)
	}
}

func ManejoRegistrosGetAll(registros []interface{}, proyecto map[string]interface{}) {
	for _, registrotemp := range registros {
		registro := registrotemp.(map[string]interface{})
		tiporegistro := registro["TipoRegistroId"].(map[string]interface{})

		if tiporegistro["Id"].(float64) == 1 && registro["Activo"] == true {
			proyecto["FechaVenimientoAcreditacion"] = registro["VencimientoActoAdministrativo"]
			proyecto["FechaVenimientoCalidad"] = nil
			if tiporegistro["Id"].(float64) == 2 && registro["Activo"] == true {
				proyecto["FechaVenimientoCalidad"] = registro["VencimientoActoAdministrativo"]
			}
		} else if tiporegistro["Id"].(float64) == 2 && registro["Activo"] == true {
			proyecto["FechaVenimientoCalidad"] = registro["VencimientoActoAdministrativo"]
		}
	}
}

// FUNCIONES QUE SE USAN EN GETONEPORID

func ManejoRegistrosGetOneId(registros []interface{}, proyecto map[string]interface{}) {
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

func ManejoUnidadesGetOneId(unidades []map[string]interface{}, idUnidad float64, proyectobase map[string]interface{}, proyecto *map[string]interface{}) {
	for _, unidad := range unidades {
		unidadTem := unidad
		idUnidad = unidadTem["Id"].(float64)
		if proyectobase["UnidadTiempoId"].(float64) == idUnidad {
			(*proyecto)["NombreUnidad"] = unidadTem["Nombre"]
		}

	}
}

func AsignarInfoProyectoGetOneId(proyecto *map[string]interface{}, proyectobase *map[string]interface{}) {
	(*proyecto)["FechaVenimientoAcreditacion"] = nil
	(*proyecto)["FechaVenimientoCalidad"] = nil
	(*proyecto)["TieneRegistroAltaCalidad"] = false
	(*proyecto)["NumeroActoAdministrativoAltaCalidad"] = nil
	(*proyecto)["AnoActoAdministrativoIdAltaCalidad"] = nil
	(*proyecto)["FechaCreacionActoAdministrativoAltaCalidad"] = nil
	(*proyecto)["VigenciaActoAdministrativoAltaCalidad"] = nil
	(*proyecto)["EnlaceActoAdministrativoAltaCalidad"] = nil

	// Información de la facultad
	var dependenciaFacultad map[string]interface{}
	errdependenciaFacultad := request.GetJson("http://"+beego.AppConfig.String("OikosService")+"/dependencia/"+fmt.Sprintf("%.f", (*proyectobase)["FacultadId"].(float64)), &dependenciaFacultad)
	// if errdependencia["Type"] == "error" || errdependencia != nil || dependencia["Status"] == "404" || dependencia["Message"] != nil {
	if errdependenciaFacultad == nil {
		(*proyecto)["NombreFacultad"] = dependenciaFacultad["Nombre"]
		(*proyecto)["IdDependenciaFacultad"] = dependenciaFacultad["Id"]
	}

	// Información de la dependencia del proyecto
	var dependencia map[string]interface{}
	errdependencia := request.GetJson("http://"+beego.AppConfig.String("OikosService")+"/dependencia/"+fmt.Sprintf("%.f", (*proyectobase)["DependenciaId"].(float64)), &dependencia)
	// if errdependencia["Type"] == "error" || errdependencia != nil || dependencia["Status"] == "404" || dependencia["Message"] != nil {
	if errdependencia == nil {
		(*proyecto)["TelefonoDependencia"] = dependencia["TelefonoDependencia"]
	}

	if (*proyectobase)["Oferta"] == true {
		(*proyecto)["OfertaLetra"] = "Si"
	} else if (*proyectobase)["Oferta"] == false {
		(*proyecto)["OfertaLetra"] = "No"
	}
	if (*proyectobase)["CiclosPropedeuticos"] == true {
		(*proyecto)["CiclosLetra"] = "Si"
	} else if (*proyectobase)["CiclosPropedeuticos"] == false {
		(*proyecto)["CiclosLetra"] = "NO"
	}
}

func ManejoProyectosGetOneId(proyectos *[]map[string]interface{}, unidades []map[string]interface{}, idUnidad float64) {
	for _, proyecto := range *proyectos {
		registros := proyecto["Registro"].([]interface{})
		proyectobase := proyecto["ProyectoAcademico"].(map[string]interface{})
		AsignarInfoProyectoGetOneId(&proyecto, &proyectobase)
		ManejoUnidadesGetOneId(unidades, idUnidad, proyectobase, &proyecto)
		ManejoRegistrosGetOneId(registros, proyecto)
	}
}

// FUNCIONES QUE SE USAN EN PUT INHABILITAR PROYECTO

func InhabilitarProyecto(alerta *models.Alert, alertas *[]interface{}, idStr string, ProyectoAcademico map[string]interface{}) {
	var resultadoProyecto map[string]interface{}
	errProyecto := request.SendJson("http://"+beego.AppConfig.String("ProyectoAcademicoService")+"/proyecto_academico_institucion/"+idStr, "PUT", &resultadoProyecto, ProyectoAcademico)
	if resultadoProyecto["Type"] == "error" || errProyecto != nil || resultadoProyecto["Status"] == "404" || resultadoProyecto["Message"] != nil {
		helpers.ManejoError(alerta, alertas, fmt.Sprintf("%v", resultadoProyecto))
	} else {
		*alertas = append(*alertas, ProyectoAcademico)
	}
}

// FUNCIONES QUE SE USAN EN PUT GET ONE REGISTRO POR ID

func ManejoRegistrosGetRegistroId(registros *[]map[string]interface{}) {
	if (*registros)[0]["Id"] != nil {
		for _, registro := range *registros {
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
}
