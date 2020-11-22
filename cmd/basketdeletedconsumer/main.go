package main

import (
	"context"
	"encoding/json"
	events "github.com/bilalislam/eshop-consumers/pkg/domain"
	"github.com/bilalislam/eshop-consumers/pkg/handler/configuration"
	"github.com/bilalislam/eshop-consumers/pkg/repository"
	"github.com/bilalislam/eshop-consumers/pkg/utils"
	"github.com/bilalislam/torc/consumer"
	"github.com/bilalislam/torc/log"
	"github.com/bilalislam/torc/rabbitmq"
	"strings"
	"time"
)

func main() {

	logger := log.GetLogger()
	configClient := configuration.NewConfigClient(".env", ".")
	h := configuration.ConfigHandler{
		ConfigClient: configClient,
	}
	config, _ := h.NewConfigHandler()
	measurement := utils.NewTimeMeasurement(logger)

	redisClient := repository.NewRedisClient(context.TODO(), config)
	onConsumed := func(message rabbitmq.Message) error {

		defer measurement.TimeTrack(time.Now(), "[Basket Deleted Consumer]", "Consumed")
		ctx := context.Background()

		orderStartedEvent, err := parseOrderStarted(message.Payload)
		if err != nil {
			logger.Error(err.Error())
			return err
		}

		_, err = redisClient.Delete(ctx, orderStartedEvent.UserId)
		if err != nil {
			logger.Error(err.Error())
			return err
		}

		return nil

	}

	r, c := consumer.AddConsumer(consumer.Request{
		Uri:           strings.Split(config.GetString("rabbitmq.uri"), ","),
		UserName:      config.GetString("rabbitmq.username"),
		Password:      config.GetString("rabbitmq.password"),
		RetryCount:    config.GetInt("rabbitmq.retry-count"),
		PrefetchCount: config.GetInt("rabbitmq.prefetch-count"),
		Exchange:      config.GetString("rabbitmq.exchange"),
		ExchangeType:  config.GetInt("rabbitmq.exchange-type"),
		Queue:         config.GetString("rabbitmq.queue"),
		RoutingKey:    config.GetString("rabbitmq.routing-key"),
	})

	c.HandleConsumer(onConsumed)
	_ = r.RunConsumers()
}

func parseOrderStarted(payload []byte) (*events.OrderStartedIntegrationEvent, error) {
	var orderStarted events.OrderStartedIntegrationEvent
	err := json.Unmarshal(payload, &orderStarted)
	if err != nil {
		return nil, err
	}
	return &orderStarted, nil
}
