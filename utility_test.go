package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestRandStr(t *testing.T) {
	Convey("Given the lenght of the string to be generated", t, func() {
		Convey("generates a random string", func() {
			str := randStr(10)
			So(len(str), ShouldEqual, 10)
		})
	})
}

func TestIsImage(t *testing.T) {
	Convey("Given a name file", t, func() {

		Convey("return true if the extension is an image extension", func() {
			jpegs := "image.jpg"
			So(isImage(jpegs), ShouldBeTrue)
			pngs := "image.png"
			So(isImage(pngs), ShouldBeTrue)
		})

		Convey("return true if the extension is an image extension, also for uppercase", func() {
			str := "image.PNG"
			So(isImage(str), ShouldBeTrue)
		})

		Convey("return false if the extension name is wrong written", func() {
			str := "imagejpg"
			So(isImage(str), ShouldBeFalse)
		})

		Convey("return false if the extension is not an image extension", func() {
			str := "text.doc"
			So(isImage(str), ShouldBeFalse)
		})
	})
}
