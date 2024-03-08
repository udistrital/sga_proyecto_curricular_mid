package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/udistrital/sga_proyecto_curricular_mid/controllers:ProyectoAcademicoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_proyecto_curricular_mid/controllers:ProyectoAcademicoController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: "/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_proyecto_curricular_mid/controllers:ProyectoAcademicoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_proyecto_curricular_mid/controllers:ProyectoAcademicoController"],
        beego.ControllerComments{
            Method: "PostProyecto",
            Router: "/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_proyecto_curricular_mid/controllers:ProyectoAcademicoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_proyecto_curricular_mid/controllers:ProyectoAcademicoController"],
        beego.ControllerComments{
            Method: "GetOnePorId",
            Router: "/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_proyecto_curricular_mid/controllers:ProyectoAcademicoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_proyecto_curricular_mid/controllers:ProyectoAcademicoController"],
        beego.ControllerComments{
            Method: "PutInhabilitarProyecto",
            Router: "/:id/inhabilitar",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_proyecto_curricular_mid/controllers:ProyectoAcademicoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_proyecto_curricular_mid/controllers:ProyectoAcademicoController"],
        beego.ControllerComments{
            Method: "PostRegistroAltaCalidadById",
            Router: "/:id/registros-alta-calidad/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_proyecto_curricular_mid/controllers:ProyectoAcademicoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_proyecto_curricular_mid/controllers:ProyectoAcademicoController"],
        beego.ControllerComments{
            Method: "PostCoordinadorById",
            Router: "/coordinador",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_proyecto_curricular_mid/controllers:ProyectoAcademicoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_proyecto_curricular_mid/controllers:ProyectoAcademicoController"],
        beego.ControllerComments{
            Method: "GetOneRegistroPorId",
            Router: "/registro/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_proyecto_curricular_mid/controllers:ProyectoAcademicoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_proyecto_curricular_mid/controllers:ProyectoAcademicoController"],
        beego.ControllerComments{
            Method: "PostRegistroCalificadoById",
            Router: "/registros-calificados/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
