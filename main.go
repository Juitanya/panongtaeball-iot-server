package main

import (
	"GoAgent/pkg/discordbot"
	"GoAgent/pkg/hwinfo"
	"GoAgent/pkg/response"
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"time"
)

func main() {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalln("can't read from env")
	}

	connStr := viper.GetString("PSQL_URL")
	if connStr == "" {
		log.Fatalln("NO DATABASE URL")
	}
	db, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatal(err)
	}

	hwClient, _ := hwinfo.NewSystemInfo()

	startUpReport(&hwClient, db)

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
	log.Println(fmt.Sprintf("HTTP server listening on port %s", appPort))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", appPort), r))
}

func startUpReport(hw *hwinfo.SystemInfo, db *pgxpool.Pool) {
	hw.PrintHostInfo(true)
	dcID := viper.GetString("DISCORD_BOT_ID")
	dcSecret := viper.GetString("DISCORD_BOT_SECRET")
	if dcID == "" || dcSecret == "" || !viper.GetBool("DISCORD_PROMPT_STARTUP") {
		return
	}
	dcBot := discordbot.NewDiscordClient(dcID, dcSecret, true, db)
	for {
		reports := hw.ToReports(true)
		vals := make([]discordbot.Embed, len(reports))
		for i := range reports {
			vals[i].Title = reports[i].Topic
			vals[i].Description = reports[i].Content
		}
		err := dcBot.SendMessage(discordbot.ThePayload{
			Content: "StartUp Report",
			Embeds:  vals,
		})
		if err == nil {
			break
		}
		time.Sleep(5 * time.Second)
	}

}
