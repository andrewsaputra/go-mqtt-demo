package mqtt

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
)

type MqttConnectConfig struct {
	BrokerUrl             string
	ClientID              string
	CleanStart            bool
	KeepAlive             uint16
	SessionExpiryInterval uint32
}

type MqttPublish struct {
	Topic         string
	Payload       []byte
	Qos           uint8
	MessageExpiry uint32
	Metadata      map[string]string
}

type MqttClient struct {
	connManager *autopaho.ConnectionManager
	timeout     time.Duration
}

func NewMqttClient(config MqttConnectConfig) (*MqttClient, error) {
	u, err := url.Parse(config.BrokerUrl)
	if err != nil {
		return nil, err
	}

	clientConfig := autopaho.ClientConfig{
		ServerUrls:                    []*url.URL{u},
		KeepAlive:                     config.KeepAlive,
		CleanStartOnInitialConnection: config.CleanStart,
		SessionExpiryInterval:         config.SessionExpiryInterval,

		ClientConfig: paho.ClientConfig{
			ClientID: config.ClientID,
			OnClientError: func(err error) {
				log.Printf("OnClientError: %s\n", err.Error())
			},
			OnServerDisconnect: func(d *paho.Disconnect) {
				if d.Properties == nil {
					log.Printf("OnServerDisconnect: Reason Code: %d\n", d.ReasonCode)
				} else {
					log.Printf("OnServerDisconnect: Reason String: %s\n", d.Properties.ReasonString)
				}
			},

			OnPublishReceived: []func(paho.PublishReceived) (bool, error){
				func(pr paho.PublishReceived) (bool, error) {
					fmt.Printf("OnPublishReceived: %s\n", pr.Packet.Payload)
					fmt.Println(pr)
					return true, nil
				},
			},
		},

		OnConnectionUp: func(cm *autopaho.ConnectionManager, c *paho.Connack) {
			log.Println("OnConnectionUp: mqtt connection established")

			//if _, err := cm.Subscribe(context.Background(), &paho.Subscribe{
			//	Subscriptions: []paho.SubscribeOptions{
			//		{Topic: "test-topic", QoS: 1},
			//	},
			//}); err != nil {
			//	fmt.Printf("failed to subscribe (%s). This is likely to mean no messages will be received.", err)
			//}
		},
		OnConnectError: func(err error) {
			log.Printf("OnConnectError : %s\n", err.Error())
		},
	}

	connectionManager, err := autopaho.NewConnection(context.Background(), clientConfig)
	if err != nil {
		return nil, err
	}

	if err = connectionManager.AwaitConnection(context.Background()); err != nil {
		return nil, err
	}

	return &MqttClient{
		connManager: connectionManager,
	}, nil
}

func (this *MqttClient) Publish(pub MqttPublish) error {
	userProps := paho.UserProperties{}
	for key, val := range pub.Metadata {
		userProps = append(userProps, paho.UserProperty{Key: key, Value: val})
	}
	_, err := this.connManager.Publish(context.Background(), &paho.Publish{
		Topic:   pub.Topic,
		Payload: pub.Payload,
		QoS:     pub.Qos,
		Properties: &paho.PublishProperties{
			MessageExpiry: &pub.MessageExpiry,
			User:          userProps,
		},
	})
	if err != nil {
		return err
	}

	//fmt.Println("response", response)
	//fmt.Println(err, pub.Payload)
	//if response.ReasonCode != 0 {
	//	return fmt.Errorf("ReasonCode: %d, Reason String: %s", response.ReasonCode, response.Properties.ReasonString)
	//}

	return nil
}

func (this *MqttClient) Subscribe(topic string, qos uint8) error {
	res, err := this.connManager.Subscribe(context.Background(), &paho.Subscribe{
		Subscriptions: []paho.SubscribeOptions{
			{Topic: topic, QoS: qos},
		},
	})
	if err != nil {
		return err
	}

	switch res.Reasons[0] {
	case 0, 1, 2:
		return nil
	default:
		return fmt.Errorf("Reason String: %s\n", res.Properties.ReasonString)
	}
}
