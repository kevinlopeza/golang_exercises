package main

import (
	"image"
	"image/color"
	"image/gif"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strconv"
)

var greenColor = color.RGBA{0, 0xFF, 0, 0xFF}
var palette = []color.Color{color.Black, greenColor}

const (
	blackIndex = 0 //First color in palette
	greenIndex = 1 //Second color in palette
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			log.Println("Error while parsing:", err)
		}
		//Some default values
		cycles := 5
		size := 100

		if r.Form["cycles"] != nil {
			cycles, _ = strconv.Atoi(r.Form["cycles"][0]) //Ignoring a possible error from Atoi
		}

		if r.Form["size"] != nil {
			size, _ = strconv.Atoi(r.Form["size"][0])
		}

		lissajous(w, cycles, size)
	})
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func lissajous(out io.Writer, cycles, size int) {
	const (
		res     = 0.001
		nframes = 64
		delay   = 8
	)

	freq := rand.Float64() * 3.0
	anim := gif.GIF{LoopCount: nframes}
	phase := 0.0
	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
		img := image.NewPaletted(rect, palette)
		for t := 0.0; t < float64(cycles)*2*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			img.SetColorIndex(size+int(x*float64(size)+0.5), size+int(y*float64(size)+0.5), greenIndex)
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}
	gif.EncodeAll(out, &anim)
}
