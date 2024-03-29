package services

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/astaxie/beego"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/requestresponse"
)

// FUNCIONES QUE SE USAN EN GETALL

func PeticionProyectos() (APIResponseDTO requestresponse.APIResponse) {
	var proyectos []map[string]interface{}
	errproyecto := request.GetJson("http://"+beego.AppConfig.String("ProyectoAcademicoService")+"/tr_proyecto_academico/", &proyectos)

	if errproyecto == nil {
		manejoProyectosGetAll(&proyectos)
		APIResponseDTO = requestresponse.APIResponseDTO(true, 200, proyectos)
		return APIResponseDTO
	} else {
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, errproyecto.Error())
		return APIResponseDTO
	}
}

func manejoProyectosGetAll(proyectos *[]map[string]interface{}) {
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

		manejoRegistrosGetAll(registros, proyecto)
	}
}

func manejoRegistrosGetAll(registros []interface{}, proyecto map[string]interface{}) {
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

func manejoRegistrosGetOneId(registros []interface{}, proyecto map[string]interface{}) {
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

func manejoUnidadesGetOneId(unidades []map[string]interface{}, idUnidad float64, proyectobase map[string]interface{}, proyecto *map[string]interface{}) {
	for _, unidad := range unidades {
		unidadTem := unidad
		idUnidad = unidadTem["Id"].(float64)
		if proyectobase["UnidadTiempoId"].(float64) == idUnidad {
			(*proyecto)["NombreUnidad"] = unidadTem["Nombre"]
		}

	}
}

func asignarInfoProyectoGetOneId(proyecto *map[string]interface{}, proyectobase *map[string]interface{}) {
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

func manejoProyectosGetOneId(proyectos *[]map[string]interface{}, unidades []map[string]interface{}, idUnidad float64) {
	for _, proyecto := range *proyectos {
		registros := proyecto["Registro"].([]interface{})
		proyectobase := proyecto["ProyectoAcademico"].(map[string]interface{})
		asignarInfoProyectoGetOneId(&proyecto, &proyectobase)
		manejoUnidadesGetOneId(unidades, idUnidad, proyectobase, &proyecto)
		manejoRegistrosGetOneId(registros, proyecto)
	}
}

func validarProyecto(errproyecto error, errunidad interface{}, proyectos *[]map[string]interface{}, unidades []map[string]interface{}, idUnidad float64) (interface{}, bool) {
	if errproyecto == nil && errunidad == nil {
		manejoProyectosGetOneId(proyectos, unidades, idUnidad)
		return proyectos, true
	} else {
		return errproyecto, false
	}
}

func PeticionProyectosGetOneId(idStr string) (APIResponseDTO requestresponse.APIResponse) {
	// var idOikos float64
	var idUnidad float64
	var proyectos []map[string]interface{}
	// var dependencias []map[string]interface{}
	var unidades []map[string]interface{}

	errproyecto := request.GetJson("http://"+beego.AppConfig.String("ProyectoAcademicoService")+"/tr_proyecto_academico/"+idStr, &proyectos)
	errunidad := request.GetJson("http://"+beego.AppConfig.String("CoreService")+"/unidad_tiempo/", &unidades)

	if proyectos[0]["ProyectoAcademico"] != nil {
		response, ok := validarProyecto(errproyecto, errunidad, &proyectos, unidades, idUnidad)
		if ok {
			APIResponseDTO = requestresponse.APIResponseDTO(true, 200, response)
		} else {
			APIResponseDTO = requestresponse.APIResponseDTO(false, 400, response)
		}
	} else {
		APIResponseDTO = requestresponse.APIResponseDTO(false, 404, nil)
	}
	return APIResponseDTO
}

// FUNCIONES QUE SE USAN EN PUT INHABILITAR PROYECTO

func InhabilitarProyecto(idStr string, data []byte) (APIResponseDTO requestresponse.APIResponse) {
	var resultadoProyecto map[string]interface{}
	var proyectoAcademico map[string]interface{}
	if err := json.Unmarshal(data, &proyectoAcademico); err == nil {

		errProyecto := request.SendJson("http://"+beego.AppConfig.String("ProyectoAcademicoService")+"/proyecto_academico_institucion/"+idStr, "PUT", &resultadoProyecto, proyectoAcademico)
		if resultadoProyecto["Type"] == "error" || errProyecto != nil || resultadoProyecto["Status"] == "404" || resultadoProyecto["Message"] != nil {
			APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, errProyecto.Error())
		} else {
			APIResponseDTO = requestresponse.APIResponseDTO(true, 200, nil)
		}
	}else {
		APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, err.Error())
	}
	return APIResponseDTO
}

// FUNCIONES QUE SE USAN EN PUT GET ONE REGISTRO POR ID

func manejoRegistrosGetRegistroId(registros *[]map[string]interface{}) {
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

func PeticionRegistrosGetRegistroId(idStr string) (APIResponseDTO requestresponse.APIResponse) {
	var registros []map[string]interface{}

	errproyecto := request.GetJson("http://"+beego.AppConfig.String("ProyectoAcademicoService")+"/registro_calificado_acreditacion/?query=ProyectoAcademicoInstitucionId.Id:"+idStr, &registros)

	if errproyecto == nil {
		manejoRegistrosGetRegistroId(&registros)
		APIResponseDTO = requestresponse.APIResponseDTO(true, 200, registros)
    } else {
        APIResponseDTO = requestresponse.APIResponseDTO(false, 400, nil, errproyecto.Error())
    }
	return APIResponseDTO
}
