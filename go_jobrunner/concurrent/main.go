package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/bamzi/jobrunner"
	"github.com/gin-gonic/gin"
)

// MyJob ...
type MyJob struct {
}

func main() {
	jobrunner.Start(10, 5)
	jobrunner.Schedule("@every 5s", MyJob{})

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET("/jobrunner/status", JobJSON)
	r.Run(":8080")
}

// JobJSON ...
func JobJSON(c *gin.Context) {
	c.JSON(http.StatusOK, jobrunner.StatusJson())
}

// Run ...
func (e MyJob) Run() {
	fmt.Println("[Start] Run MyJob!")
	time.Sleep(30 * time.Second)
	fmt.Println("[End] Run MyJob!")
}
