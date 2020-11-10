package api

import (
	"bytes"
	"encoding/json"
	"image"
	"image/jpeg"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/BurntSushi/graphics-go/graphics"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func gracefulExit(w http.ResponseWriter, text string) {
	log.Println(text)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// Handler telegram hook
func Handler(w http.ResponseWriter, r *http.Request) {
	token, ok := os.LookupEnv("BOT_TOKEN")
	if !ok {
		gracefulExit(w, "no telegram bot token")
		return
	}
	bot, err := tgbotapi.NewBotAPI(token)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		gracefulExit(w, "no body")
		return
	}
	var update tgbotapi.Update
	if err := json.Unmarshal(body, &update); err != nil {
		gracefulExit(w, "body parse error")
		return
	}
	message := update.Message
	if message.Document == nil {
		gracefulExit(w, "no file")
		return
	}
	url, err := bot.GetFileDirectURL(message.Document.FileID)
	if err != nil {
		gracefulExit(w, "file error")
		return
	}
	client := http.Client{}

	res, err := client.Get(url)
	if err != nil {
		gracefulExit(w, "file load error")
		return
	}
	defer res.Body.Close()

	buf := new(bytes.Buffer)
	err = blur(res.Body, buf)
	if err != nil {
		gracefulExit(w, "cannot blur this file")
		return
	}
	msg := tgbotapi.NewDocumentUpload(int64(message.Chat.ID), tgbotapi.FileBytes{
		Name:  "blured.png",
		Bytes: buf.Bytes(),
	})

	// todo: looks like send message can be done via response on this webhook request
	// but let's use this API for now
	bot.Send(msg)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func blur(file io.Reader, w io.Writer) error {
	srcImage, _, err := image.Decode(file)
	if err != nil {
		return err
	}
	dstImage := image.NewRGBA(srcImage.Bounds())
	graphics.Blur(dstImage, srcImage, &graphics.BlurOptions{StdDev: 5.5})
	jpeg.Encode(w, dstImage, &jpeg.Options{jpeg.DefaultQuality})
	return nil
}
