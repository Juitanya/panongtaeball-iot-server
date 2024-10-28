package discordbot

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"io"
	"log"
	"net/http"
	"runtime"
	"time"
)

type DiscordClient struct {
	ID    string `json:"id"`
	Token string `json:"token"`

	actionBy    int64
	logFunction bool
	db          *pgxpool.Pool
}

type ThePayload struct {
	Content string  `json:"content"`
	Embeds  []Embed `json:"embeds,omitempty"`
}

type Embed struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type T struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func NewDiscordClient(id, token string, lf bool, db *pgxpool.Pool) DiscordClient {
	if db == nil {
		log.Println("No Postgres No log")
	}
	return DiscordClient{
		ID:          id,
		Token:       token,
		db:          db,
		logFunction: lf,
	}
}

func (dc *DiscordClient) SetActByWho(id int64) {
	dc.actionBy = id
}

func (dc *DiscordClient) SendMessage(tp ThePayload) error {
	if dc.actionBy == 0 {
		log.Println("No ActBy Or System Called")
	}

	const timeOut = 90 * time.Second
	url := fmt.Sprintf("https://discord.com/api/webhooks/%s/%s", dc.ID, dc.Token)

	client := &http.Client{
		Timeout: timeOut,
	}

	payloadBuf := new(bytes.Buffer)
	err := json.NewEncoder(payloadBuf).Encode(tp)
	if err != nil {
		fmt.Println(err)
		return err
	}
	req, err := http.NewRequest(http.MethodPost, url, payloadBuf)
	if err != nil {
		fmt.Println(err)
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}

	body, err := io.ReadAll(io.NopCloser(res.Body))
	if err != nil {
		return err
	}
	log.Println("RawBody --> ", string(body))
	if dc.db != nil {
		dc.logMsg(tp)
		if dc.logFunction {
			pc, file, line, ok := runtime.Caller(1)
			if ok {
				fn := runtime.FuncForPC(pc).Name()
				dc.logAction(fn, file, line)
			}
		}
	}

	return nil
}

func (dc *DiscordClient) logMsg(tp ThePayload) {
	payloadJSON, _ := json.Marshal(tp)

	// Insert query
	query := "INSERT INTO discord_toggle_histories (action_by, payload) VALUES ($1, $2) ;"

	dc.db.QueryRow(context.Background(), query, dc.actionBy, payloadJSON)

}

func (dc *DiscordClient) logAction(name, file string, line int) {
	associateWith := "discord_toggle_histories"
	calledByFunction := name
	lineNum := int64(line)
	fileLocation := file

	query := `
    INSERT INTO function_histories (associate_with, called_by_function, line, file_location)
    VALUES ($1, $2, $3, $4);
    `

	dc.db.QueryRow(context.Background(), query, associateWith, calledByFunction, lineNum, fileLocation)

}
