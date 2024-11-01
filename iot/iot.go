package iot

import (
	"GoAgent/pkg/response"
	"fmt"
	"net/http"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type IotHandler struct {
	MqttClient mqtt.Client
}

func publish(client mqtt.Client) {
	num := 5
	for i := 0; i < num; i++ {
		text := fmt.Sprintf("Message %d", i)
		token := client.Publish("topic/test", 0, false, text)
		token.Wait()
		time.Sleep(time.Second)
	}
}

func sub(client mqtt.Client) {
	topic := "topic/test"
	token := client.Subscribe(topic, 1, nil)
	token.Wait()
	fmt.Printf("Subscribed to topic %s", topic)
}

func (i IotHandler) UpdateLight(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	sub(i.MqttClient)
	fmt.Println("test params: ", i.MqttClient)
	render.JSON(w, r, response.HTTPResponse{
		Data:  id,
		Error: nil,
	})
}
