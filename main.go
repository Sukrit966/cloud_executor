package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
	"sync"
	
	"github.com/labstack/echo"
)

type post_data struct {
	Lang       string
	FileData   string
	FileName   string
	Input 	   string
	ComCommand string
	ExeCommand string
}

type response_data struct {
	Error  string `json:"error"`
	Output string `json:"output"`
}



func exe_cmd(cmd string, wg *sync.WaitGroup) string {
	fmt.Println("command is ", cmd)
	parts := strings.Fields(cmd)
	head := parts[0]
	parts = parts[1:len(parts)]

	out, err := exec.Command(head, parts...).Output()
	if err != nil {
		fmt.Printf("%s", err)
	}

	fmt.Printf("%s", out)
	wg.Done()
	return string(out)
}

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.POST("/exe",func(c echo.Context) error{
		d := new(post_data)
		if err := c.Bind(d); err != nil {
			return err
		}
		data := []byte(d.FileData)
		err := ioutil.WriteFile("./" + d.FileName, data, 0644)
		if err != nil {
			panic(err)
		}
		wg := new(sync.WaitGroup)
		wg.Add(2)
		compile := exe_cmd(d.ComCommand, wg)
		execute := exe_cmd(d.ExeCommand, wg)

		wg.Wait()
		out := &response_data{Error: compile, Output: execute}

		return c.JSON(http.StatusOK, out)
	})
	e.Logger.Fatal(e.Start(":8080"))
}