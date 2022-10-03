package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

var configMap = make(map[string]map[string]any)

func initJson() {

	file, err := os.ReadFile("./host.json")
	if err != nil {
		panic(err)
	}

	var res []map[string]any
	json.Unmarshal(file, &res)

	for _, container := range res {
		configMap[container["token"].(string)] = container
	}
}

type ConnParam struct {
	Token string `form:"token" binding:"required"`
	Param string `form:"param"`
}

func jsonS(c *gin.Context) {

	var connParam ConnParam
	if err := c.ShouldBindQuery(&connParam); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	container, ok := configMap[connParam.Token]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not find container"})
		return
	}

	c.JSON(http.StatusOK, container)
}

type BootOptions struct {
	Port int
}

var bootOptions = &BootOptions{}

func parseFlags() {
	flag.IntVar(&bootOptions.Port, "p", 2234, "The server listening port (default 2234)")
	flag.Parse()
}

func main() {
	parseFlags()
	initJson()
	server := gin.Default()
	server.GET("/json", jsonS)
	server.Run(fmt.Sprintf(":%d", bootOptions.Port)) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
