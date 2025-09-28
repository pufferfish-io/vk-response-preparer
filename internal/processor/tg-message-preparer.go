package processor

import (
	"context"
	"encoding/json"

	"vk-response-preparer/internal/contract"
)

type Producer interface {
	Send(_ context.Context, topic string, data []byte) error
}

type TgMessagePreparer struct {
	producer   Producer
	kafkaTopic string
}

func NewVkMessagePreparer(kt string, p Producer) *TgMessagePreparer {
	return &TgMessagePreparer{
		producer:   p,
		kafkaTopic: kt,
	}
}

func (t *TgMessagePreparer) Handle(ctx context.Context, raw []byte) error {
    var requestMessage contract.NormalizedResponse
    if err := json.Unmarshal(raw, &requestMessage); err != nil {
        return err
    }

    msg := contract.SendMessageRequest{
        PeerID:  requestMessage.ChatID,
        Message: requestMessage.Text,
    }

    out, err := json.Marshal(msg)
    if err != nil {
        return err
    }

	return t.producer.Send(ctx, t.kafkaTopic, out)
}
