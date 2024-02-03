package publisherapi

import (
	"andrewsaputra/go-mqtt-demo/internal/mqtt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type PublisherApiHandler struct {
	MqttClient *mqtt.MqttClient
	validate   *validator.Validate
}

func NewPublisherApiHandler(brokerUrl string) (*PublisherApiHandler, error) {
	client, err := mqtt.NewMqttClient(mqtt.MqttConnectConfig{
		BrokerUrl:             brokerUrl,
		ClientID:              "test-client-1",
		CleanStart:            true,
		KeepAlive:             60,
		SessionExpiryInterval: 120,
	})
	if err != nil {
		return nil, err
	}

	return &PublisherApiHandler{
		MqttClient: client,
		validate:   validator.New(validator.WithRequiredStructEnabled()),
	}, nil
}

func (this *PublisherApiHandler) Publish(c *gin.Context) {
	request := PublishRequest{Qos: 1, MessageExpiry: 3600} //initialize defaults
	if err := c.BindJSON(&request); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	if err := this.validate.Struct(request); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	err := this.MqttClient.Publish(mqtt.MqttPublish{
		Topic:         request.Topic,
		Payload:       []byte(request.Payload),
		Qos:           request.Qos,
		MessageExpiry: request.MessageExpiry,
		Metadata:      request.Metadata,
	})
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.String(http.StatusOK, "message published")
}
