package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"os"
	"bytes"
	"strings"
	"sync"
	"encoding/gob"
	"time"
	
	"github.com/labstack/echo"
)


type log_write_data struct {
	Content string
	FileName string
	ComCommand string
	ExeCommand string
}

type post_data struct {
	Lang       string  `json:"lang"`
	FileData   string  `json:"fdata"`
	FileName   string  `json:"fname"`
	Input 	   string  `json:"input"`
	Extension  string  `json:"ext"`
	ComCommand string  `json:"ccommand"`
	ExeCommand string  `json:"ecommand"`
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
		return c.File("./public/index.html")
	})

	e.GET("/main.js",func(c echo.Context) error {
		return c.File("./public/main.js")
	})

	e.GET("/execute",func(c echo.Context) error{
		return c.String(http.StatusOK,"POST Only Link")
	})

	e.POST("/execute",func(c echo.Context) error{
		var buff bytes.Buffer 
		enc := gob.NewEncoder(&buff) 
		d := new(post_data)
		err1 := enc.Encode(d)
		if err1 != nil {
			return err1;
		}

		if err := c.Bind(d); err != nil {
			return err
		}
		data := []byte(d.FileData)
		err := ioutil.WriteFile("./" + d.FileName + d.Extension, data, 0644)
		if err != nil {
			panic(err)
		}

		wg := new(sync.WaitGroup)
		wg.Add(2)
		compile := exe_cmd(d.ComCommand, wg)

		t := time.Now().UTC()
		timeStamp := t.Format("20060102150405")
		os.Mkdir("./log/" + timeStamp,0644)
		err2 := ioutil.WriteFile("./log/" + timeStamp +  "/" + "source_" +  d.FileName + d.Extension ,data,0644 )
		if err2 != nil {
			panic(err2)
		}
		execute := exe_cmd(d.ExeCommand, wg)

		wg.Wait()
		out := &response_data{Error: compile, Output: execute}

		return c.JSON(http.StatusOK, out)
	})
	e.Logger.Fatal(e.Start(":8080"))
}
