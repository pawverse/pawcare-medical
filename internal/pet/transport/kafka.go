package transport

import (
	"fmt"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v3/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	watermilltransport "github.com/pawverse/pawcare-core/pkg/kit/transport/watermill"
	"github.com/pawverse/pawcare-medical/internal/pet/endpoint"
	"github.com/pawverse/pawcare-profiles/pkg/events"
)

func RegisterKafkaRoutes(router *message.Router, subscriber *kafka.Subscriber, publisher *kafka.Publisher, endpoints endpoint.Set) error {
	petCreateSubscriber := watermilltransport.NewSubscriber(
		endpoints.CreateEndpoint,
		watermilltransport.DecodeJSONMessage[endpoint.CreateRequest],
		watermilltransport.EncodeResponse,
	)

	handler := router.AddNoPublisherHandler("create_pet", fmt.Sprintf("profiles.%s", events.EventPetCreated), subscriber, petCreateSubscriber.Handle)
	poisonMiddleware, err := middleware.PoisonQueue(publisher, fmt.Sprintf("profiles.%s.poison", events.EventPetCreated))
	if err != nil {
		return err
	}

	handler.AddMiddleware(poisonMiddleware,
		middleware.Retry{
			MaxRetries:      3,
			InitialInterval: 50 * time.Millisecond,
			Logger: router.Logger().With(watermill.LogFields{
				"middleware": "retry",
			}),
		}.Middleware,
	)

	return nil
}
