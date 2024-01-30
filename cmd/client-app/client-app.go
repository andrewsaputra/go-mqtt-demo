package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const ADDRESS = ":4000"

var startTime time.Time = time.Now()

func main() {
	router := gin.Default()
	router.GET("/status", statusCheck)

	router.Run(ADDRESS)
}

func statusCheck(c *gin.Context) {
	body := make(map[string]string)
	body["status"] = "Healthy"
	body["started_at"] = startTime.Format(time.RFC822Z)

	c.JSON(http.StatusOK, body)
}
