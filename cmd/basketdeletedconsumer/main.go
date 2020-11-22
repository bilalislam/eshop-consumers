package main

import (
	"github.com/bilalislam/eshop-consumers"
	"github.com/bilalislam/torc/consumer"
	"github.com/bilalislam/torc/log"
	"github.com/bilalislam/torc/rabbitmq"
	"strings"
)

func main() {

	logger := log.GetLogger()
	configClient := configuration.NewConfigClient(".env", ".")
	h := configuration.ConfigHandler{
		ConfigClient: configClient,
	}
	config, _ := h.NewConfigHandler()
	_, eventsUseCase := handler.NewEventUseCase(config)
	_, projectionUseCase := handler.NewProjectionUseCase(config)
	measurement := utils.NewTimeMeasurement(logger)

	onConsumed := func(message rabbitmq.Message) error {
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
