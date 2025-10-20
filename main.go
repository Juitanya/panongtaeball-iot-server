package main

import (
	"Panong/iot/light"
	"Panong/iot/valve"
	"Panong/pkg/hwinfo"
	"Panong/pkg/response"
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

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("X-Auth-Token")

		if token == "" || token != viper.GetString("HEADER_SECRET_AUTH") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalln("can't read from env")
	}

	hwClient, _ := hwinfo.NewSystemInfo()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(AuthMiddleware)
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
	opts.SetAutoReconnect(true)
	opts.SetResumeSubs(true)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	r.Mount("/light", LightRoutes(client))
	r.Mount("/valve", ValveRoutes(client))

	log.Printf("HTTP server listening on port %s", appPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", appPort), r))
}

func LightRoutes(mqttClient mqtt.Client) chi.Router {
	r := chi.NewRouter() // สร้าง router ใหม่
	lightHandler := light.LightHandler{
		MqttClient: mqttClient,
	}

	r.Get("/{light}", lightHandler.Light)
	r.Put("/{light}/{action}", lightHandler.UpdateLight)
	return r
}

func ValveRoutes(mqttClient mqtt.Client) chi.Router {
	r := chi.NewRouter() // สร้าง router ใหม่
	valveHandler := valve.ValveHandler{
		MqttClient: mqttClient,
	}

	r.Get("/{valve}", valveHandler.Valve)
	r.Put("/{valve}/{action}", valveHandler.UpdateValve)
	return r
}
