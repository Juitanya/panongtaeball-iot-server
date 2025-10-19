package valve

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

type ValveHandler struct {
	MqttClient mqtt.Client
}

type Payload struct {
	State string `json:"state"`
}

var ErrFriendlyNameNotFound = errors.New("friendly name not found.")

func (v ValveHandler) Valves() []string {
	return []string{
		viper.GetString("WATER_VALVE"),
	}
}

func (v ValveHandler) getFriendlyName(light string) (string, error) {
	valves := v.Valves()
	switch light {
	case valves[0]:
		return "วาล์วน้ำ", nil
	default:
		return "", ErrFriendlyNameNotFound
	}
}

func (v ValveHandler) getZigbee2MQTTValveStatus(client mqtt.Client, valve string) (string, error) {
	// 1. ตรวจสอบว่ามีการเชื่อมต่อ MQTT หรือไม่
	if !client.IsConnected() {
		token := client.Connect()
		token.Wait()
		if token.Error() != nil {
			return "", fmt.Errorf("MQTT client not connected: %w", token.Error())
		}
	}

	// 2. แปลงชื่อ friendly name
	friendlyName, err := v.getFriendlyName(valve)
	if err != nil {
		return "", err
	}

	subscribedTopic := fmt.Sprintf("zigbee2mqtt/%s", friendlyName)
	getTopic := fmt.Sprintf("zigbee2mqtt/%s/get", valve)

	// 3. สร้าง channel สำหรับรับสถานะ
	statusCh := make(chan string, 1)

	// 4. Subscribe topic ของ light
	token := client.Subscribe(subscribedTopic, 0, func(client mqtt.Client, msg mqtt.Message) {
		statusCh <- string(msg.Payload())
	})
	token.Wait()
	if err := token.Error(); err != nil {
		return "", fmt.Errorf("failed to subscribe: %w", err)
	}

	// 5. Publish request สถานะ
	payload, _ := json.Marshal(map[string]string{"state": "", "battery": "", "flow": ""})
	pubToken := client.Publish(getTopic, 0, false, payload)
	pubToken.Wait()
	if err := pubToken.Error(); err != nil {
		client.Unsubscribe(subscribedTopic)
		return "", fmt.Errorf("failed to publish get request: %w", err)
	}

	// 6. รอ response จาก channel
	select {
	case status := <-statusCh:
		unsubToken := client.Unsubscribe(subscribedTopic)
		unsubToken.Wait()
		return status, nil
	case <-time.After(30 * time.Second):
		unsubToken := client.Unsubscribe(subscribedTopic)
		unsubToken.Wait()
		return "", errors.New("timeout waiting for valve status")
	}
}

func (v ValveHandler) updateZigbee2MQTTValve(client mqtt.Client, action, valve string) error {
	payload, err := json.Marshal(Payload{State: action})
	if err != nil {
		return err
	}

	if !slices.Contains(v.Valves(), valve) {
		return errors.New("valve not found")
	}

	token := client.Publish(fmt.Sprintf("zigbee2mqtt/%s/set", valve), 0, false, string(payload))
	token.Wait()
	time.Sleep(time.Second)

	return nil
}

func (v ValveHandler) UpdateValve(w http.ResponseWriter, r *http.Request) {
	valve := chi.URLParam(r, "valve")
	action := chi.URLParam(r, "action")

	err := v.updateZigbee2MQTTValve(v.MqttClient, action, valve)
	if err != nil {
		http.Error(w, "Failed to publish message", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Valve updated successfully"))
}

func (v ValveHandler) Valve(w http.ResponseWriter, r *http.Request) {
	fmt.Println("get VALVEEEEE...")
	valve := chi.URLParam(r, "valve")
	if valve == "" {
		http.Error(w, "valve param required", http.StatusBadRequest)
		return
	}

	status, err := v.getZigbee2MQTTValveStatus(v.MqttClient, valve)
	if err != nil {
		if errors.Is(err, ErrFriendlyNameNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		fmt.Println("error:", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(status))
}
