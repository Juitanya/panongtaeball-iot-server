package light

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-chi/chi/v5"
	"github.com/spf13/viper"
)

type LightHandler struct {
	MqttClient mqtt.Client
}

type Payload struct {
	State string `json:"state"`
}

func publish(client mqtt.Client, action string) {
	payload, err := json.Marshal(Payload{State: action})
	if err != nil {
		fmt.Println("light publish error: ", err)
	}

	token := client.Publish(fmt.Sprintf("zigbee2mqtt/%s/set", viper.GetString("FIRST_LIGHT")), 0, false, string(payload))
	token.Wait()
	time.Sleep(time.Second)
}

func (l LightHandler) UpdateLight(w http.ResponseWriter, r *http.Request) {
	action := chi.URLParam(r, "action")
	publish(l.MqttClient, action)

	return
}
