package services

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/udistrital/sga_mid_proyecto_curricular/models"
	"github.com/udistrital/utils_oas/request"
)

func ManejoProyectos(proyectos *[]map[string]interface{}) {
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

		ManejoRegistros(registros, proyecto)
	}
}

func ManejoRegistros(registros []interface{}, proyecto map[string]interface{}) {
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

func ManejoError(alerta *models.Alert, alertas *[]interface{}, mensaje string, err ...error) {
	var msj string
	if len(err) > 0 && err[0] != nil {
		msj = mensaje + err[0].Error()
	} else {
		msj = mensaje
	}
	*alertas = append(*alertas, msj)
	(*alerta).Body = *alertas
	(*alerta).Type = "error"
	(*alerta).Code = "400"
}
