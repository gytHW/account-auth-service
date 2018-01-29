package api

import (
	"account-auth-service/logic"
	"fmt"
	"github.com/emicklei/go-restful"
	"io"
)

var (
	WsContainer = restful.NewContainer()
)

func init() {

	ws := new(restful.WebService)

	ws.Path("/").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/Hello").To(hello).Doc("keepalive").Operation("hello"))
	ws.Route(ws.POST("/").To(handle).Reads(logic.Req{}).Writes(logic.Resp{}))

	WsContainer.Add(ws)

	cors := restful.CrossOriginResourceSharing{
		AllowedHeaders: []string{"Content-Type", "Accept"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		CookiesAllowed: false,
		Container:      WsContainer}
	//WsContainer.Filter(cors.Filter)
	fmt.Println(cors)
	WsContainer.Filter(WsContainer.OPTIONSFilter)

}

type HelloMsg struct {
	Id     string `json:"id"`
	Count  int8   `json:"count"`
	Enable bool   `json:"enable"`
}

func hello(req *restful.Request, resp *restful.Response) {
	//resp.ResponseWriter.Header().Set("Content-Disposition", "attachment; filename=hello.txt")
	//resp.ResponseWriter.Header().Set("Content-Type", req.Request.Header.Get("Content-Type"))

	msg := HelloMsg{}
	//req.ReadEntity(&msg)
	fmt.Println(msg)
	io.WriteString(resp, "world")
}

func handle(req *restful.Request, resp *restful.Response) {

	req_data := &logic.Req{}
	err := req.ReadEntity(req_data)

	fmt.Println(req_data)

	if err != nil {
		resp.WriteEntity(&logic.Resp{1, err.Error(), 0, nil})
		return
	}
	switch req_data.Action {
	case "CreatePassword":
		logic.CreatePassword(req_data, resp)
		break
	case "ValidatePassword":
		logic.ValidatePassword(req_data, resp)
		break
	case "CreateInitialPassword":
		logic.CreateInitialPassword(req_data, resp)
		break
	case "ChangePassword":
		logic.ChangePassword(req_data, resp)
		break
	case "CreateVertificationCode":
		logic.CreateVertificationCode(req_data, resp)
		break
	case "CheckInitialPassword":
		logic.CheckInitialPassword(req_data, resp)
		break
	case "ResetPassword":
		logic.ResetPassword(req_data, resp)
		break
	default:
		resp.WriteEntity(&logic.Resp{1, "incorrect api name", 0, nil})
		break
	}

	return
}
