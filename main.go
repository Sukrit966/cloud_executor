package main

import (
	"github.com/labstack/echo"
	"github.com/sukrit966/cloud_executor/server"
)

type log_write_data struct {
	Content    string
	FileName   string
	ComCommand string
	ExeCommand string
}

func main() {
	e := echo.New()
	e.GET("/", server.Home)
	e.GET("/ping", server.Ping)
	e.GET("/main.js", server.JsFile)

	e.GET("/execute", server.GetExectue)

	e.POST("/execute", server.Execute)
	e.Logger.Fatal(e.Start(":8080"))
}
