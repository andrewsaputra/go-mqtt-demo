package main

import (
	publisherapi "andrewsaputra/go-mqtt-demo/internal/publisher-api"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const ADDRESS = ":3001"

var startTime time.Time = time.Now()

func main() {
	//todo : move to config file
	handler, err := publisherapi.NewPublisherApiHandler("mqtt://localhost:1883")
	if err != nil {
		log.Panic(err)
	}

	err = handler.MqttClient.Subscribe("test-topic-1", 1)
	if err != nil {
		log.Panic(err)
	}
	handler.MqttClient.Subscribe("test-topic-2", 1)
	handler.MqttClient.Subscribe("test-topic-3", 1)

	router := gin.Default()
	router.GET("/status", statusCheck)
	router.POST("/publish", handler.Publish)

	router.Run(ADDRESS)
}

func statusCheck(c *gin.Context) {
	body := make(map[string]string)
	body["status"] = "Healthy"
	body["started_at"] = startTime.Format(time.RFC822Z)

	c.JSON(http.StatusOK, body)
}
