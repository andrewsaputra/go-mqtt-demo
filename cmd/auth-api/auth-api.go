package main

import (
	authapi "andrewsaputra/go-mqtt-demo/internal/auth-api"
	"encoding/base64"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const ADDRESS = ":3000"

var startTime time.Time = time.Now()

func main() {
	//todo : load key from parameter store
	key := "/2VwP8jzGw6TUpAAYinR5SJamGLaIhlScQhXhd+/bik="
	keyBytes, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		log.Panic(err)
	}

	//todo : move duration to config
	sessionDurationSecs := 120
	sessionManager := authapi.NewSessionManagerImpl(int64(sessionDurationSecs))
	handler, err := authapi.NewAuthApiHandler(keyBytes, sessionManager)
	if err != nil {
		log.Panic(err)
	}

	router := gin.Default()
	router.GET("/status", statusCheck)
	router.GET("/newsession", handler.NewSession)

	router.Run(ADDRESS)
}

func statusCheck(c *gin.Context) {
	body := make(map[string]string)
	body["status"] = "Healthy"
	body["started_at"] = startTime.Format(time.RFC822Z)

	c.JSON(http.StatusOK, body)
}
