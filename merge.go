package main

import (
	"flag"
	"fmt"
	comp "github.com/edap/gomerger/composer"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"io"
	"io/ioutil"
	"log"
	//"math"
	"os"
	"path/filepath"
	//errors
)

var thumb_width = flag.Int("thumb_width", 120, "the width of a single thumb")
var thumb_height = flag.Int("thumb_height", 90, "the height of a single thumb")

// var allow_scaling = scale images if it is a prime number
var erase_original = flag.Bool("erase_original", false, "erase the original thumb after being merged into the grid")
var source_dir = flag.String("source_dir", "/home/da/to_merge", "the origin directory that contains the images to compose the grid")
var target_dir = flag.String("target_dir", "/home/da/to_merge/merged", "the destination directory that will containe the final grid")

// implementare log, o almeno, avere una politica coerente sugli errori

// questo sparira' non ha senso dire che 0 e' un error, quando di default, il pacchetto resize diche che se
// h o w sono 0, scala l'immagine in proporzione
type wrongArgumentError struct{ arg string }

func (a *wrongArgumentError) Error() string {
	return fmt.Sprintf("The argument %s can not be minor than 1", a.arg)
}

type Thumb struct {
	width          int
	height         int
	desired_width  int
	desired_height int
	img_name       string
}

func (t *Thumb) Scale(img_path string) error {
	are_equal, err := t.HasDesiredDimension()
	if err != nil {
		return err
	}
	if are_equal == true {
		return nil
	}

	file, err := os.Open(img_path)
	if err != nil {
		return err
	}

	img, err := jpeg.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()

	m := resize.Resize(uint(t.desired_width), uint(t.desired_height), img, resize.NearestNeighbor)
	out, err := os.Create(img_path)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	return jpeg.Encode(out, m, nil)
}

func (t *Thumb) Copy(src_path, dst_path string) (int64, error) {
	src_file, err := os.Open(src_path)
	if err != nil {
		return 0, err
	}
	defer src_file.Close()

	src_file_stat, err := src_file.Stat()
	if err != nil {
		return 0, err
	}
	if !src_file_stat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src_path)
	}

	dst_file, err := os.Create(dst_path)
	if err != nil {
		return 0, err
	}
	defer dst_file.Close()
	return io.Copy(dst_file, src_file)
}

func (t *Thumb) Move(src_path, dst_path string) error {
	return os.Rename(src_path, dst_path)
}

func (t *Thumb) HasDesiredDimension() (bool, error) {

	if t.desired_height < 1 {
		return false, &wrongArgumentError{"thumb_height"}
	}
	if t.desired_width < 1 {
		return false, &wrongArgumentError{"thumb_width"}
	}

	if t.desired_width == t.width && t.desired_height == t.height {
		return true, nil
	} else {
		return false, nil
	}
}

// func (t *Thumb) getNameAndDimension()(name string, height int, error) {
// 	return "fds",12, nil
// }

func main() {
	flag.Parse()
	if err := os.MkdirAll(*target_dir, 0755); err != nil {
		log.Fatal("impossible to create target directory")
	}
	rv := comp.CalcPrimeFactors(10)
	fmt.Println(rv)
	files, _ := ioutil.ReadDir(*source_dir)
	for _, imgFile := range files {
		if imgFile.IsDir() { // controllare se e' jpg, invece che se e' directory
			continue
		}
		img_name := imgFile.Name()
		dst_path := filepath.Join(*source_dir, img_name)
		src_path := filepath.Join(*target_dir, img_name)
		if reader, err := os.Open(src_path); err == nil {
			defer reader.Close()

			im, _, err := image.DecodeConfig(reader)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s: %v\n", imgFile.Name(), err)
				continue
			}

			// qui va implementata la logica che conta le immagini, calcola un quadro,
			// e dice su quante linee devono essere disposte le immagini

			thumb := &Thumb{
				width:          im.Width,
				height:         im.Height,
				desired_width:  *thumb_width,
				desired_height: *thumb_height,
				img_name:       img_name,
			}

			if *erase_original == true {
				thumb.Move(src_path, dst_path)
			} else {
				thumb.Copy(src_path, dst_path)
			}
			thumb.Scale(dst_path)

			fmt.Printf("%s %d %d\n", imgFile.Name(), thumb.width, thumb.height)
		} else {
			fmt.Println("no")
		}
	}

	// now go to the merged folder, count the images
	// use the composer
	// create an empty square
	// past the images in the empty square
}

// conta le immagini che ha nella directory, cerca di rispettare le proporzioni
// (date come parametro), e mette quello che ci sta. Quello che ci sta lo cancella (opzione delete?) quello che non ci sta viene salvato su un file, che fa da registro

// package main

// import (
// 	"fmt"
// 	"image"
// 	"image/draw"
// 	"image/jpeg"
// 	"os"
// )

// func main() {
// 	fImg1, _ := os.Open("/home/da/to_merge/Gwx_TPMT3Bg.jpg")
// 	defer fImg1.Close()
// 	img1, _, _ := image.Decode(fImg1)

// 	m := image.NewRGBA(image.Rect(0, 0, 800, 600))
// 	draw.Draw(m, m.Bounds(), img1, image.Point{0, 0}, draw.Src)
// 	//draw.Draw(m, m.Bounds(), img2, image.Point{-200,-200}, draw.Src)

// 	toimg, _ := os.Create("/home/da/to_merge/new.jpg")
// 	defer toimg.Close()

// 	jpeg.Encode(toimg, m, &jpeg.Options{jpeg.DefaultQuality})
// }

// http://golang.org/src/pkg/image/jpeg/reader.go?s=10946:10998#L343
// http://golang.org/src/pkg/image/jpeg/reader.go?s=10744:10789#L336
// func (t Thumb) Scale() (int, error) {
//   // if HasDesiredDimension
//   m := resize.Resize(300, 200, t.decoded_img, resize.NearestNeighbor)
//   out, _ := os.Create("test_resized.jpg")
//   // if err != nil {
//   //  log.Fatal(err)
//   // }
//   defer out.Close()

//   // write new image to file
//   jpeg.Encode(out, m, nil)
//   return 2, nil
// }

// CREATE IMAGE
// func createImage() (bool, error) {
// 	m := image.NewRGBA(image.Rect(0, 0, 120, 90))

// 	out, err := os.Create("test_images/120x90.jpg")
// 	if err != nil {
// 		fmt.Println("imm non create")
// 	}
// 	defer out.Close()

// 	// write new image to file
// 	jpeg.Encode(out, m, nil)
// 	return true, err
// }
