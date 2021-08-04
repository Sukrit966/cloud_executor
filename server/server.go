package server

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/labstack/echo"
)

type post_data struct {
	Lang       string `json:"lang"`
	FileData   string `json:"fdata"`
	FileName   string `json:"fname"`
	MaxTime    int    `json:"maxTime"`
	Input      string `json:"input"`
	Extension  string `json:"ext"`
	ComCommand string `json:"ccommand"`
	ExeCommand string `json:"ecommand"`
}

type response_data struct {
	Error  string `json:"error"`
	Output string `json:"output"`
}

type version_response struct {
	Version string `json:"version"`
	CPU_ID  string `json:"CPU_ID"`
}

type error_response struct {
	Error string `json:"error"`
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

func Home(c echo.Context) error {
	return c.File("./public/index.html")
}

func JsFile(c echo.Context) error {
	return c.File("./public/main.js")
}

func Ping(c echo.Context) error {
	vr := &version_response{Version: "0.1", CPU_ID: "x86"}
	return c.JSON(http.StatusOK, vr)
}

func GetExectue(c echo.Context) error {
	return c.String(http.StatusMethodNotAllowed, "POST Only Link")
}

func Execute(c echo.Context) error {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	d := new(post_data)
	err1 := enc.Encode(d)
	if err1 != nil {
		e := &error_response{Error: err1.Error()}
		return c.JSON(http.StatusOK, e)
	}

	if err := c.Bind(d); err != nil {
		e := &error_response{Error: err.Error()}
		return c.JSON(http.StatusOK, e)
	}
	data := []byte(d.FileData)
	err := ioutil.WriteFile("./"+d.FileName+d.Extension, data, 0644)
	if err != nil {
		e := &error_response{Error: err.Error()}
		return c.JSON(http.StatusOK, e)
	}

	wg := new(sync.WaitGroup)
	wg.Add(2)
	compile := exe_cmd(d.ComCommand, wg)

	t := time.Now().UTC()
	timeStamp := t.Format("20060102150405")
	os.Mkdir("./log/"+timeStamp, 0644)
	err2 := ioutil.WriteFile("./log/"+timeStamp+"/"+"source_"+d.FileName+d.Extension, data, 0644)
	if err2 != nil {
		e := &error_response{Error: err2.Error()}
		return c.JSON(http.StatusOK, e)
	}
	execute := exe_cmd(d.ExeCommand, wg)

	wg.Wait()
	out := &response_data{Error: compile, Output: execute}

	return c.JSON(http.StatusOK, out)
}
