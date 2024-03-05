package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/udistrital/sga_proyecto_curricular_mid/controllers:ConsultaProyectoAcademicoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_proyecto_curricular_mid/controllers:ConsultaProyectoAcademicoController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: "/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_proyecto_curricular_mid/controllers:ConsultaProyectoAcademicoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_proyecto_curricular_mid/controllers:ConsultaProyectoAcademicoController"],
        beego.ControllerComments{
            Method: "GetOnePorId",
            Router: "/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_proyecto_curricular_mid/controllers:ConsultaProyectoAcademicoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_proyecto_curricular_mid/controllers:ConsultaProyectoAcademicoController"],
        beego.ControllerComments{
            Method: "GetOneRegistroPorId",
            Router: "/get_registro/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_proyecto_curricular_mid/controllers:ConsultaProyectoAcademicoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_proyecto_curricular_mid/controllers:ConsultaProyectoAcademicoController"],
        beego.ControllerComments{
            Method: "PutInhabilitarProyecto",
            Router: "/inhabilitar_proyecto/:id",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_proyecto_curricular_mid/controllers:CrearProyectoAcademicoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_proyecto_curricular_mid/controllers:CrearProyectoAcademicoController"],
        beego.ControllerComments{
            Method: "PostProyecto",
            Router: "/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_proyecto_curricular_mid/controllers:CrearProyectoAcademicoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_proyecto_curricular_mid/controllers:CrearProyectoAcademicoController"],
        beego.ControllerComments{
            Method: "PostCoordinadorById",
            Router: "/coordinador",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
