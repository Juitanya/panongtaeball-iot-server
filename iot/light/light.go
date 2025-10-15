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

var ErrFriendlyNameNotFound = errors.New("friendly name not found.")

func (l LightHandler) Lights() []string {
	return []string{
		viper.GetString("FIRST_LIGHT"),
		viper.GetString("SECOND_LIGHT"),
		viper.GetString("THIRD_LIGHT"),
		viper.GetString("FOURTH_LIGHT"),
		viper.GetString("IN_FRONT_OF_CLUBHOUSE_LOGO_LIGHT"),
		viper.GetString("TRIPLE_PLUGS"),
	}
}

func (l LightHandler) getFriendlyName(light string) (string, error) {
	lights := l.Lights()
	switch light {
	case lights[0]:
		return "ไฟสนาม1", nil
	case lights[1]:
		return "ไฟสนาม2", nil
	case lights[2]:
		return "ไฟสนาม3", nil
	case lights[3]:
		return "ไฟสนาม4", nil
	case lights[4]:
		return "ไฟโลโก้หน้าคลับเฮ้าส์", nil
	default:
		return "", ErrFriendlyNameNotFound
	}
}

func (l LightHandler) getZigbee2MQTTLightStatus(client mqtt.Client, light string) (string, error) {
	statusCh := make(chan string)

	friendlyName, err := l.getFriendlyName(light)
	if err != nil {
		return "", err
	}

	subscribedTopic := fmt.Sprintf("zigbee2mqtt/%s", friendlyName)
	getTopic := fmt.Sprintf("zigbee2mqtt/%s/get", light)

	token := client.Subscribe(subscribedTopic, 0, func(client mqtt.Client, msg mqtt.Message) {
		statusCh <- string(msg.Payload())
	})
	token.Wait()
	if err := token.Error(); err != nil {
		return "", err
	}

	payload, _ := json.Marshal(map[string]string{"state": ""})
	pubToken := client.Publish(getTopic, 0, false, payload)
	pubToken.Wait()
	if err := pubToken.Error(); err != nil {
		return "", err
	}

	select {
	case status := <-statusCh:
		unsubToken := client.Unsubscribe(subscribedTopic)
		unsubToken.Wait()
		return status, nil
	case <-time.After(1 * time.Minute):
		return "", errors.New("timeout waiting for status")
	}
}

func (l LightHandler) updateZigbee2MQTTLight(client mqtt.Client, action, light string) error {
	payload, err := json.Marshal(Payload{State: action})
	if err != nil {
		return err
	}

	if !slices.Contains(l.Lights(), light) {
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

	err := l.updateZigbee2MQTTLight(l.MqttClient, action, light)
	if err != nil {
		http.Error(w, "Failed to publish message", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Light updated successfully"))
}

func (l LightHandler) Light(w http.ResponseWriter, r *http.Request) {
	light := chi.URLParam(r, "light")
	if light == "" {
		http.Error(w, "light param required", http.StatusBadRequest)
		return
	}

	status, err := l.getZigbee2MQTTLightStatus(l.MqttClient, light)
	if err != nil {
		if errors.Is(err, ErrFriendlyNameNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(status))
}
