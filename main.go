package main

import (
	"flag"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"log"
	"os"
	"path/filepath"
)

func main() {
	var (
		thumb_width  = flag.Int("thumb_width", 120, "the width of a single thumb")
		thumb_height = flag.Int("thumb_height", 90, "the height of a single thumb")
		source_dir   = flag.String("source_dir", ".", "the origin directory that contains the images to compose the grid")
		dest_dir     = flag.String("dest_dir", ".", "the destination directory that will contain the grid")
		log_file     = flag.String("log_file", "stdout", "specify a log file, as default it will print on stdout")
	)
	flag.Parse()

	// set a log file if it's required
	if *log_file != "stdout" {
		f, err := os.OpenFile(*log_file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		defer f.Close()
		log.SetOutput(f)
	}

	tot, images := CountImagesAndVerifyPreconditions(*source_dir, *dest_dir)

	//calculate the dimension of the rectangle
	res := map[string]int{"area": tot, "height": 0, "base": 0, "skipped": 0}
	rect := calculateRectangle(res)
	log.Printf("%d images will be skipped", res["skipped"])
	log.Printf("%d images will be merged together", res["area"])

	// calculate the position of each image in the final canvas
	positions := calculatePositions(rect, images, *thumb_width, *thumb_height)

	// give a name to the canvas file and prepare it
	canvas_filename := filepath.Join(*dest_dir, randStr(20)+".jpg")
	canvas_image := image.NewRGBA(image.Rect(0, 0, *thumb_width*res["base"], *thumb_height*res["height"]))

	// sei arrivato qui, devi decidere come gestire gli errori nella goroutine che vuoi creare

	c_errore := make(chan error)
	c_immagine := make(chan Img)

	// iterate through the images, resize if necessary, decode and add to the canvas
	for _, image_path := range images {
		// thumb := NewThumb(
		// 	*thumb_width, *thumb_height, image_path, )
		// thumb.SetDimension()

		thumb := NewThumb(
			*thumb_width,
			*thumb_height,
			image_path,
		)
		go DoSmth(thumb)

		img, err := DecodedThumb(thumb)
		CopyToCanvas(err, positions, image_path, canvas_image, img)
	}

	toimg, _ := os.Create(canvas_filename)
	defer toimg.Close()
	jpeg.Encode(toimg, canvas_image, &jpeg.Options{jpeg.DefaultQuality})

	log.Printf("canvas %s succesfully created", canvas_filename)
}

func DecodedThumb(i Img) (image.Image, error) {
	i.SetDimension()
	return i.DecodeIt()
}

func CopyToCanvas(err error, positions map[string][2]int, image_path string, canvas_image *image.RGBA, img image.Image) error {
	if err != nil {
		log.Printf("it was not possible to decode the image %s: %v", image_path, err)
	} else {
		x := positions[image_path][0]
		y := positions[image_path][1]
		draw.Draw(canvas_image, canvas_image.Bounds(), img, image.Point{x, y}, draw.Src)
	}
	return err
}

func DoSmth(c chan Img) {
	t := <-c
	fmt.Println(t.Width())
}

func CountImagesAndVerifyPreconditions(source_dir string, destination_dir string) (int, []string) {
	// At least 2 images has to be present in the source directory
	tot, images := listFiles(source_dir)
	if tot < 2 {
		log.Fatal("There are less than two images in this folder")
	}

	// create the destination directory
	err := createDirectory(destination_dir)
	if err != nil {
		log.Fatalf("impossible to create destination directory: %v", destination_dir)
	}
	return tot, images
}
