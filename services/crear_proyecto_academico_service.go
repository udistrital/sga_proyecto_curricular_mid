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

// FUNCIONES QUE SE USAN EN PUT GET ONE REGISTRO POR ID

func RegistrarCoordinador(alertas *[]interface{}, alerta *models.Alert, CoordinadorNuevo map[string]interface{}) {
	var resultadoCoordinadorNuevo map[string]interface{}
	CoordinadorNuevo["FechaFinalizacion"] = "0001-01-01T00:00:00-05:00"

	errRegistro := request.SendJson("http://"+beego.AppConfig.String("ProyectoAcademicoService")+"/proyecto_academico_rol_tercero_dependencia", "POST", &resultadoCoordinadorNuevo, CoordinadorNuevo)
	if resultadoCoordinadorNuevo["Type"] == "error" || errRegistro != nil || resultadoCoordinadorNuevo["Status"] == "404" || resultadoCoordinadorNuevo["Message"] != nil {
		helpers.ManejoError(alerta, alertas, fmt.Sprintf("%v", resultadoCoordinadorNuevo))
	} else {
		*alertas = append(*alertas, CoordinadorNuevo)
	}
}

func ManejoCoordinadorAntiguo(alertas *[]interface{}, alerta *models.Alert, CoordinadorAntiguos []map[string]interface{}) {
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
			}
		} else {
			fmt.Println("Todos los registros estan nulos")
		}

	}
}
