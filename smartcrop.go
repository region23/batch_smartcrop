package main

import (
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/muesli/smartcrop"
	"github.com/nfnt/resize"
)

// SubImager это черт его знает что за конструкция
type SubImager interface {
	SubImage(r image.Rectangle) image.Image
}

func writeImageToJpeg(img *image.Image, name string) {
	fso, err := os.Create(name)
	if err != nil {
		panic(err)
	}
	defer fso.Close()

	jpeg.Encode(fso, (*img), &jpeg.Options{Quality: 100})
}

func writeImageToPng(img *image.Image, name string) {
	fso, err := os.Create(name)
	if err != nil {
		panic(err)
	}
	defer fso.Close()

	png.Encode(fso, (*img))
}

func main() {
	// принимаем параметрами директорию с изображениями, префикс для новых изображений и новые размеры
	pathPtr := flag.String("path", ".", "Path to directory with images. If param is not set then used current directory")
	prefixPtr := flag.String("prefix", "thumb", "Prefix for thumbnail images")
	widthPtr := flag.Int("width", 220, "width of thumbnail image")
	heightPtr := flag.Int("height", 124, "height of thumbnail image")

	flag.Parse()

	// открываем заданную директорию. Если директория не задана то открываем текущую директорию
	var dir string
	var err error
	if *pathPtr == "." {
		dir = *pathPtr
	} else {
		dir, err = filepath.Abs(os.Getenv("HOME") + *pathPtr)

		if err != nil {
			log.Fatal(err)
		}
	}

	files, err := ioutil.ReadDir(dir)

	if err != nil {
		log.Fatal(err)
	}

	// извлекаем список всех изображений
	for _, f := range files {
		ext := filepath.Ext(f.Name())

		if ext == ".jpg" {
			// в цикле уменьшаем изображение и сохраняем с префиксом в текущую директорию

			fi, _ := os.Open(filepath.Join(dir, f.Name()))
			defer fi.Close()

			img, _, _ := image.Decode(fi)
			topCrop, _ := smartcrop.SmartCrop(&img, *widthPtr, *heightPtr)
			fmt.Printf("Top crop: %+v\n", topCrop)
			fmt.Println(f.Name())

			sub, ok := img.(SubImager)
			if ok {
				newName := filepath.Join(dir, *prefixPtr+"_"+f.Name())
				cropImage := sub.SubImage(image.Rect(topCrop.X, topCrop.Y, topCrop.Width+topCrop.X, topCrop.Height+topCrop.Y))
				// resize
				newImage := resize.Resize(uint(*widthPtr), uint(*heightPtr), cropImage, resize.Lanczos3)
				writeImageToJpeg(&newImage, newName)

			} else {
				log.Fatal("No SubImage support")
			}
		}

	}
}
