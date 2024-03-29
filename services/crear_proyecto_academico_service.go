package services

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/requestresponse"
	"github.com/udistrital/utils_oas/time_bogota"
)

// FUNCIONES QUE SON LLAMADAS POR LOS SERVICES DE CADA ENDPOINT

func ManejoPeticionesProyecto(data []byte) (APIResponseDTO requestresponse.APIResponse) {
	var Proyecto_academico *map[string]interface{}

	if err := json.Unmarshal(data, &Proyecto_academico); err == nil {
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

		if !peticionOikos(&resultadoOikos, Proyecto_academico_oikosPost, Proyecto_academico, &Proyecto_academicoPost) {
			APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil)
			return APIResponseDTO
		}

		if !peticionProyecto(&resultadoProyecto, Proyecto_academicoPost, *Proyecto_academico) {
			APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil)
			return APIResponseDTO
		}
		APIResponseDTO = requestresponse.APIResponseDTO(true, 200, nil)
		return APIResponseDTO
	} else {
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, err.Error())
		return APIResponseDTO
	}
}

func CoordinaById(data []byte) (APIResponseDTO requestresponse.APIResponse) {
	var CoordinadorNuevo map[string]interface{}
	var resultado map[string]interface{}
	var status = 0

	if err := json.Unmarshal(data, &CoordinadorNuevo); err == nil {
		if resultado["Type"] != "error" {
			var CoordinadorAntiguos []map[string]interface{}
			idStr := fmt.Sprintf("%v", CoordinadorNuevo["ProyectoAcademicoInstitucionId"].(map[string]interface{})["Id"])

			errcordinador := request.GetJson("http://"+beego.AppConfig.String("ProyectoAcademicoService")+"/proyecto_academico_rol_tercero_dependencia/?query=ProyectoAcademicoInstitucionId.Id:"+idStr, &CoordinadorAntiguos)
			if errcordinador == nil {
				if CoordinadorAntiguos[0]["Id"] != nil {
					if exito := ManejoCoordinadorAntiguo(CoordinadorAntiguos); !exito {
						status = 400
					}
					if exito := RegistrarCoordinador(CoordinadorNuevo); !exito {
						status = 400
					}

					APIResponseDTO = requestresponse.APIResponseDTO(true, 200, nil)
					return APIResponseDTO
				} else {
					if err := json.Unmarshal(data, &CoordinadorNuevo); err == nil {
						if exito := RegistrarCoordinador(CoordinadorNuevo); !exito {
							status = 400
						} else {
							status = 200
						}
					} else {
						status = 400
					}
				}
			} else {
				APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, errcordinador.Error())
				return APIResponseDTO
			}
		} else {
			if resultado["Body"] == "<QuerySeter> no row found" {
				APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil)
				return APIResponseDTO
			} else {
				APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, resultado["Body"])
				return APIResponseDTO
			}
		}
	} else {
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, err.Error())
		return APIResponseDTO
	}
	if status != 400 {
		APIResponseDTO = requestresponse.APIResponseDTO(true, 200, nil)
	} else if status == 404 {
		APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil)
	} else {
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil)
	}
	return APIResponseDTO
}

func PostRegistroCalificadoById(data []byte) requestresponse.APIResponse {
	var Registro_nuevo map[string]interface{}
	var resultado map[string]interface{}

	// Deserializar el JSON de entrada
	if err := json.Unmarshal(data, &Registro_nuevo); err != nil {
		return requestresponse.APIResponseDTO(false, 400, "Error en JSON Unmarshal: "+err.Error())
	}

	// Verificar si hay un error en el resultado
	if resultado["Type"] == "error" {
		return requestresponse.APIResponseDTO(false, 400, "Error en el resultado: "+fmt.Sprintf("%v", resultado["Body"]))
	}

	// Obtener registros antiguos
	idStr := fmt.Sprintf("%v", Registro_nuevo["ProyectoAcademicoInstitucionId"].(map[string]interface{})["Id"])
	var registros_antiguos_acreditacion []map[string]interface{}
	erregistro := request.GetJson("http://"+beego.AppConfig.String("ProyectoAcademicoService")+"/registro_calificado_acreditacion/?query=ProyectoAcademicoInstitucionId:"+idStr+",TipoRegistroId.Id:1", &registros_antiguos_acreditacion)

	if erregistro != nil {
		return requestresponse.APIResponseDTO(false, 400, "Error al obtener registros antiguos: "+erregistro.Error())
	}

	// Verificar y actualizar registros antiguos si es necesario
	if fmt.Sprintf("%v", registros_antiguos_acreditacion) != "[map[]]" {
		for _, registro := range registros_antiguos_acreditacion {
			registro["Activo"] = false
			idRegistro := registro["Id"].(float64)

			var resultado map[string]interface{}
			errregistrocambiado := request.SendJson("http://"+beego.AppConfig.String("ProyectoAcademicoService")+"/registro_calificado_acreditacion/"+strconv.FormatFloat(idRegistro, 'f', -1, 64), "PUT", &resultado, registro)

			if errregistrocambiado != nil || resultado["Type"] == "error" || resultado["Status"] == "404" || resultado["Message"] != nil {
				return requestresponse.APIResponseDTO(false, 400, "Error al actualizar registro: "+fmt.Sprintf("%v", resultado))
			}
		}
	}

	// Crear un nuevo registro
	var resultadoRegistroNuevo map[string]interface{}
	errRegistro := request.SendJson("http://"+beego.AppConfig.String("ProyectoAcademicoService")+"/registro_calificado_acreditacion", "POST", &resultadoRegistroNuevo, Registro_nuevo)

	if errRegistro != nil || resultadoRegistroNuevo["Type"] == "error" || resultadoRegistroNuevo["Status"] == "404" || resultadoRegistroNuevo["Message"] != nil {
		return requestresponse.APIResponseDTO(false, 400, "Error al enviar nuevo registro: "+fmt.Sprintf("%v", resultadoRegistroNuevo))
	}

	return requestresponse.APIResponseDTO(true, 200, resultadoRegistroNuevo)
}

func PostRegistroAltaCalidadById(id string, data []byte) requestresponse.APIResponse {
	var Registro_nuevo map[string]interface{}
	var resultado map[string]interface{}

	if resultado["Type"] != "error" {
		var registros_antiguos_alta_calidad []map[string]interface{}

		erregistro := request.GetJson("http://"+beego.AppConfig.String("ProyectoAcademicoService")+"/registro_calificado_acreditacion/?query=ProyectoAcademicoInstitucionId:"+id+",TipoRegistroId.Id:2", &registros_antiguos_alta_calidad)
		if erregistro != nil {
			return requestresponse.APIResponseDTO(false, 400, erregistro.Error())
		}

		if len(registros_antiguos_alta_calidad) > 0 && registros_antiguos_alta_calidad[0]["Id"] != nil {
			for _, registro := range registros_antiguos_alta_calidad {
				registro["Activo"] = false
				idRegistro := registro["Id"].(float64)

				var resultado map[string]interface{}
				errregistrocambiado := request.SendJson("http://"+beego.AppConfig.String("ProyectoAcademicoService")+"/registro_calificado_acreditacion/"+strconv.FormatFloat(idRegistro, 'f', -1, 64), "PUT", &resultado, registro)

				if errregistrocambiado != nil || resultado["Type"] == "error" || resultado["Status"] == "404" || resultado["Message"] != nil {
					return requestresponse.APIResponseDTO(false, 400, "Error al actualizar registro: "+fmt.Sprintf("%v", resultado))
				}
			}
		}

		if err := json.Unmarshal(data, &Registro_nuevo); err != nil {
			return requestresponse.APIResponseDTO(false, 400, "Error en JSON Unmarshal: "+err.Error())
		}

		var resultadoRegistroNuevo map[string]interface{}
		errRegistro := request.SendJson("http://"+beego.AppConfig.String("ProyectoAcademicoService")+"/registro_calificado_acreditacion", "POST", &resultadoRegistroNuevo, Registro_nuevo)

		if errRegistro != nil || resultadoRegistroNuevo["Type"] == "error" || resultadoRegistroNuevo["Status"] == "404" || resultadoRegistroNuevo["Message"] != nil {
			return requestresponse.APIResponseDTO(false, 400, "Error al enviar nuevo registro: "+fmt.Sprintf("%v", resultadoRegistroNuevo))
		}

		return requestresponse.APIResponseDTO(true, 200, "Registro creado con éxito")
	}

	return requestresponse.APIResponseDTO(false, 400, "Error en la respuesta: "+fmt.Sprintf("%v", resultado))
}

// FUNCIONES QUE SE USAN EN PUT GET ONE POST COORDINADOR BY ID

func asignarProyectoAcademico(Proyecto_academico *map[string]interface{}, resultadoOikos map[string]interface{}, Proyecto_academicoPost *map[string]interface{}) {
	idDependenciaProyecto := resultadoOikos["HijaId"].(map[string]interface{})["Id"]
	(*Proyecto_academicoPost)["ProyectoAcademicoInstitucion"].(map[string]interface{})["DependenciaId"] = idDependenciaProyecto
}

func peticionOikos(resultadoOikos *map[string]interface{}, Proyecto_academico_oikosPost interface{}, Proyecto_academico *map[string]interface{}, Proyecto_academicoPost *map[string]interface{}) bool {
	errOikos := request.SendJson("http://"+beego.AppConfig.String("OikosService")+"/dependencia_padre/tr_dependencia_padre", "POST", resultadoOikos, Proyecto_academico_oikosPost)
	if (*resultadoOikos)["Type"] == "error" || errOikos != nil || (*resultadoOikos)["Status"] == "404" || (*resultadoOikos)["Message"] != nil {

		return false
	} else {
		asignarProyectoAcademico(Proyecto_academico, *resultadoOikos, Proyecto_academicoPost)
		return true
	}
}

func peticionProyecto(resultadoProyecto *map[string]interface{}, Proyecto_academicoPost map[string]interface{}, Proyecto_academico map[string]interface{}) bool {
	errProyecto := request.SendJson("http://"+beego.AppConfig.String("ProyectoAcademicoService")+"/tr_proyecto_academico", "POST", resultadoProyecto, Proyecto_academicoPost)
	if (*resultadoProyecto)["Type"] == "error" || errProyecto != nil || (*resultadoProyecto)["Status"] == "404" || (*resultadoProyecto)["Message"] != nil {

		return false
	} else {
		return true
	}
}

// FUNCIONES QUE SE USAN EN PUT GET ONE POST COORDINADOR BY ID

func RegistrarCoordinador(CoordinadorNuevo map[string]interface{}) bool {
	var resultadoCoordinadorNuevo map[string]interface{}
	CoordinadorNuevo["FechaFinalizacion"] = "0001-01-01T00:00:00-05:00"

	errRegistro := request.SendJson("http://"+beego.AppConfig.String("ProyectoAcademicoService")+"/proyecto_academico_rol_tercero_dependencia", "POST", &resultadoCoordinadorNuevo, CoordinadorNuevo)
	if resultadoCoordinadorNuevo["Type"] == "error" || errRegistro != nil || resultadoCoordinadorNuevo["Status"] == "404" || resultadoCoordinadorNuevo["Message"] != nil {

		return false
	} else {

		return true
	}
}

func ManejoCoordinadorAntiguo(CoordinadorAntiguos []map[string]interface{}) bool {
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
				return false
			}
		} else {
			fmt.Println("Todos los registros estan nulos")
		}

	}
	return true
}
