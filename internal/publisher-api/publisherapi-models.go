package publisherapi

type PublishRequest struct {
	Topic         string            `json:"topic" validate:"required"`
	Payload       string            `json:"payload" validate:"required"`
	Qos           uint8             `json:"qos" validate:"gte=0,lte=2"`
	MessageExpiry uint32            `json:"message_expiry" validate:"gt=0"`
	Metadata      map[string]string `json:"metadata"`
}
