package speech

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

type SpeechHandler struct {
}

func (h SpeechHandler) Test(w http.ResponseWriter, r *http.Request) {
	filePath := "./iot/speech/sounds/timeout.mp3"

	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	defer streamer.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))

	<-done

	// HTTP Response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Speech successfully played"))
}
