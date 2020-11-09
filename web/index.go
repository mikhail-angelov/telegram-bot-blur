package web

import (
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"net/http"

	"github.com/BurntSushi/graphics-go/graphics"
	"github.com/gorilla/mux"
)

// GetRouter returns the router for the API
func GetRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", WebHandler).Methods(http.MethodPost)
	return r
}

// Handler entry point
func WebHandler(w http.ResponseWriter, r *http.Request) {
	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	r.ParseMultipartForm(10 << 20)
	// FormFile returns the first file for the given key `file`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	w.Header().Set("Content-Type", "image/jpeg")
	err = blur(file, w)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func blur(file io.Reader, w io.Writer) error {
	srcImage, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("decode error")
		return err
	}
	dstImage := image.NewRGBA(srcImage.Bounds())
	graphics.Blur(dstImage, srcImage, &graphics.BlurOptions{StdDev: 5.5})
	return jpeg.Encode(w, dstImage, &jpeg.Options{jpeg.DefaultQuality})
}
