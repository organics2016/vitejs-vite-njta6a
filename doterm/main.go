package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"net/http"
	"organics.ink/doterm/connector"
)

type ConnParam struct {
	Token string `form:"token" binding:"required"`
	Param string `form:"param"`
}

func doterm(c *gin.Context) {

	websocket := connector.InitWebSocket(c.Writer, c.Request, context.Background())
	defer websocket.Close()

	var connParam ConnParam
	if err := c.ShouldBindQuery(&connParam); err != nil {
		websocket.OutputError(err)
		c.Status(http.StatusBadRequest)
		return
	}

	if err := conn(&connParam, websocket); err != nil {
		websocket.OutputError(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

type BootOptions struct {
	Server string
}

var bootOptions = &BootOptions{}

func parseFlags() {
	flag.StringVar(&bootOptions.Server, "s", "http://127.0.0.1/", "aaaaaa")
	flag.Parse()
}

func main() {
	parseFlags()
	fmt.Printf("%v\n", bootOptions)

	server := gin.Default()
	server.Use(static.Serve("/", static.LocalFile("../client/dist", false)))
	server.GET("/doterm", doterm)
	server.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
