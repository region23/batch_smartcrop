package main

import (
	"flag"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/muesli/smartcrop"
	"github.com/nfnt/resize"

)

// SubImager - I don't understand this construction
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
	// get params
	pathPtr := flag.String("path", ".", "Path to directory with images. If param is not set then used current directory")
	prefixPtr := flag.String("prefix", "thumb", "Prefix for thumbnail images")
	widthPtr := flag.Int("width", 220, "width of thumbnail image")
	heightPtr := flag.Int("height", 124, "height of thumbnail image")

	flag.Parse()

	// open directory with files
	var dir string
	var err error
	// if path not set using current directory
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

	// get all files
	for _, f := range files {
		ext := strings.ToLower(filepath.Ext(f.Name()))

		// proceed only images (jpg or png)
		if ext == ".jpg" || ext == ".png" {
			fi, _ := os.Open(filepath.Join(dir, f.Name()))
			defer fi.Close()

			// crop
			img, _, err := image.Decode(fi)
            if err != nil {
		      log.Fatal(err)
	        }
    
			topCrop, err := smartcrop.SmartCrop(img, *widthPtr, *heightPtr)
            if err != nil {
		      log.Fatal(err)
	        }
			//fmt.Printf("Top crop: %+v\n", topCrop)
			//fmt.Println(f.Name())

			sub, ok := img.(SubImager)
			if ok {
				newName := filepath.Join(dir, *prefixPtr+"_"+f.Name())
				cropImage := sub.SubImage(image.Rect(topCrop.X, topCrop.Y, topCrop.Width+topCrop.X, topCrop.Height+topCrop.Y))
				// resize
				newImage := resize.Resize(uint(*widthPtr), uint(*heightPtr), cropImage, resize.Lanczos3)
				if ext == ".jpg" {
					writeImageToJpeg(&newImage, newName)
				}
				if ext == ".png" {
					writeImageToPng(&newImage, newName)
				}

			} else {
				log.Fatal("No SubImage support")
			}
		}

	}
}
