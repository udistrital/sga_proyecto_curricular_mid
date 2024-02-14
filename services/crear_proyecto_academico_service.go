package services

import (
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/udistrital/sga_mid_proyecto_curricular/helpers"
	"github.com/udistrital/sga_mid_proyecto_curricular/models"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/time_bogota"
)

// FUNCIONES QUE SE USAN EN PUT GET ONE POST COORDINADOR BY ID

func asignarProyectoAcademico(alertas *[]interface{}, Proyecto_academico *map[string]interface{}, resultadoOikos map[string]interface{}, Proyecto_academicoPost *map[string]interface{}) {
	*alertas = append(*alertas, *Proyecto_academico)
	idDependenciaProyecto := resultadoOikos["HijaId"].(map[string]interface{})["Id"]
	(*Proyecto_academicoPost)["ProyectoAcademicoInstitucion"].(map[string]interface{})["DependenciaId"] = idDependenciaProyecto
}

func peticionOikos(resultadoOikos *map[string]interface{}, Proyecto_academico_oikosPost interface{}, alerta *models.Alert, alertas *[]interface{}, Proyecto_academico *map[string]interface{}, Proyecto_academicoPost *map[string]interface{}) bool {
	errOikos := request.SendJson("http://"+beego.AppConfig.String("OikosService")+"/dependencia_padre/tr_dependencia_padre", "POST", resultadoOikos, Proyecto_academico_oikosPost)
	if (*resultadoOikos)["Type"] == "error" || errOikos != nil || (*resultadoOikos)["Status"] == "404" || (*resultadoOikos)["Message"] != nil {
		helpers.ManejoError(alerta, alertas, fmt.Sprintf("%v", resultadoOikos), errOikos)
		return false
	} else {
		asignarProyectoAcademico(alertas, Proyecto_academico, *resultadoOikos, Proyecto_academicoPost)
		return true
	}
}

func peticionProyecto(resultadoProyecto *map[string]interface{}, Proyecto_academicoPost map[string]interface{}, alerta *models.Alert, alertas *[]interface{}, Proyecto_academico map[string]interface{}) bool {
	errProyecto := request.SendJson("http://"+beego.AppConfig.String("ProyectoAcademicoService")+"/tr_proyecto_academico", "POST", resultadoProyecto, Proyecto_academicoPost)
	if (*resultadoProyecto)["Type"] == "error" || errProyecto != nil || (*resultadoProyecto)["Status"] == "404" || (*resultadoProyecto)["Message"] != nil {
		helpers.ManejoError(alerta, alertas, "", errProyecto)
		return false
	} else {
		*alertas = append(*alertas, Proyecto_academico)
		return true
	}
}

func ManejoPeticionesProyecto(Proyecto_academico *map[string]interface{}, alerta *models.Alert, alertas *[]interface{}) bool {
	Proyecto_academicoPost := make(map[string]interface{})
	Proyecto_academicoPost = map[string]interface{}{
		"ProyectoAcademicoInstitucion": (*Proyecto_academico)["ProyectoAcademicoInstitucion"],
		"Enfasis":                      (*Proyecto_academico)["Enfasis"],
		"Registro":                     (*Proyecto_academico)["Registro"],
		"Titulaciones":                 (*Proyecto_academico)["Titulaciones"],
	}

	Proyecto_academico_oikosPost := (*Proyecto_academico)["Oikos"]

	var resultadoOikos map[string]interface{}
	var resultadoProyecto map[string]interface{}

	if !peticionOikos(&resultadoOikos, Proyecto_academico_oikosPost, alerta, alertas, Proyecto_academico, &Proyecto_academicoPost) {
		return false
	}

	if !peticionProyecto(&resultadoProyecto, Proyecto_academicoPost, alerta, alertas, *Proyecto_academico) {
		return false
	}
	return true
}

// FUNCIONES QUE SE USAN EN PUT GET ONE POST COORDINADOR BY ID

func RegistrarCoordinador(alertas *[]interface{}, alerta *models.Alert, CoordinadorNuevo map[string]interface{}) bool {
	var resultadoCoordinadorNuevo map[string]interface{}
	CoordinadorNuevo["FechaFinalizacion"] = "0001-01-01T00:00:00-05:00"

	errRegistro := request.SendJson("http://"+beego.AppConfig.String("ProyectoAcademicoService")+"/proyecto_academico_rol_tercero_dependencia", "POST", &resultadoCoordinadorNuevo, CoordinadorNuevo)
	if resultadoCoordinadorNuevo["Type"] == "error" || errRegistro != nil || resultadoCoordinadorNuevo["Status"] == "404" || resultadoCoordinadorNuevo["Message"] != nil {
		helpers.ManejoError(alerta, alertas, fmt.Sprintf("%v", resultadoCoordinadorNuevo))
		return false
	} else {
		*alertas = append(*alertas, CoordinadorNuevo)
		return true
	}
}

func ManejoCoordinadorAntiguo(alertas *[]interface{}, alerta *models.Alert, CoordinadorAntiguos []map[string]interface{}) bool {
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
				helpers.ManejoError(alerta, alertas, fmt.Sprintf("%v", resultado))
				return false
			}
		} else {
			fmt.Println("Todos los registros estan nulos")
		}

	}
	return true
}
