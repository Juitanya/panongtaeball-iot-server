package light

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"slices"
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

func publish(client mqtt.Client, action, light string) error {
	payload, err := json.Marshal(Payload{State: action})
	if err != nil {
		fmt.Println("light publish error: ", err)
	}

	lights := []string{viper.GetString("FIRST_LIGHT"), viper.GetString("SECOND_LIGHT"), viper.GetString("THIRD_LIGHT"), viper.GetString("FOURTH_LIGHT")}

	if !slices.Contains(lights, light) {
		return errors.New("light not found")
	}

	token := client.Publish(fmt.Sprintf("zigbee2mqtt/%s/set", light), 0, false, string(payload))
	token.Wait()
	time.Sleep(time.Second)

	return nil
}

func (l LightHandler) UpdateLight(w http.ResponseWriter, r *http.Request) {
	light := chi.URLParam(r, "light")
	action := chi.URLParam(r, "action")
	err := publish(l.MqttClient, action, light)
	if err != nil {
		http.Error(w, "Failed to publish message", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Light updated successfully"))
}
