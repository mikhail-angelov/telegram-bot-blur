package api

import (
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/BurntSushi/graphics-go/graphics"
)

// Handler entry point
func Handler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("TelegramWebhookHandler %s\n", body)

	// telegram_api.HandleUpdateFromTelegram(body)
}

func blur(file io.Reader, w io.Writer) {
	srcImage, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("decode error")
		fmt.Println(err)
		return
	}
	dstImage := image.NewRGBA(srcImage.Bounds())
	graphics.Blur(dstImage, srcImage, &graphics.BlurOptions{StdDev: 5.5})
	jpeg.Encode(w, dstImage, &jpeg.Options{jpeg.DefaultQuality})
}
