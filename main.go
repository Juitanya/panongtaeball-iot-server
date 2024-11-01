package main

import (
	"GoAgent/iot"
	"GoAgent/pkg/hwinfo"
	"GoAgent/pkg/response"
	"fmt"
	"log"
	"net/http"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/spf13/viper"
)

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}

func main() {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalln("can't read from env")
	}

	hwClient, _ := hwinfo.NewSystemInfo()

	r := chi.NewRouter()
	// Add middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(1 * time.Minute))

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, response.HTTPResponse{
			Data:  hwClient.Host,
			Error: nil,
		})
	})

	appPort := viper.GetString("APP_PORT")
	if appPort == "" {
		appPort = "5000"
	}

	var broker = viper.GetString("BROKER")
	var port = 1883
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetClientID(viper.GetString("CLIENT_ID"))
	opts.SetUsername(viper.GetString("USERNAME"))
	opts.SetPassword(viper.GetString("PASSWORD"))
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	r.Mount("/iot", IotRoutes(client))

	// sub(client)
	// publish(client)

	// client.Disconnect(3000)

	log.Println(fmt.Sprintf("HTTP server listening on port %s", appPort))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", appPort), r))
}

func IotRoutes(mqttClient mqtt.Client) chi.Router {
	r := chi.NewRouter()
	iotHandler := iot.IotHandler{
		MqttClient: mqttClient,
	}
	r.Get("/{id}", iotHandler.UpdateLight)
	return r
}
