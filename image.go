package main

import (
	"encoding/base64"
	"fmt"
	"os"
)

var (
	imgFilePath = "/tmp/"
	// imgFilePath = fixPath(os.Getenv("IMG_FILE_PATH"))
)

// fixPath ensures that a path string ends in '/', adding it if necessary
func fixPath(s string) string {
	fmt.Println(s)
	if s == "" {
		return "/tmp/"
	} else if s[len(s)-1] != '/' {
		return s + "/"
	}
	return s
}

// Uint64ToFilename takes a uint64 and format and returns the corresponding filename
func Uint64ToFilename(ii uint64, format string) string {
	s := fmt.Sprintf("img%#x.%s", ii, format)
	return s
}

// MakeUniqueFilenameGenerator creates a function to generate unique filenames given a `format`
func MakeUniqueFilenameGenerator(start uint64) func(string) (string, error) {
	c := make(chan uint64)

	go func() {
		var i uint64
		for i = start; ; i++ {
			c <- i
		}
	}()
	return func(format string) (string, error) {
		i := <-c
		return Uint64ToFilename(i, format), nil
	}
}

// SaveToFile saves `data` to `filename`
func SaveToFile(data []byte, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}

// SaveImgToFile saves image `data` to file `filename`
func SaveImgToFile(img Image) error {
	sd, err := base64.StdEncoding.DecodeString(img.Data)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return SaveToFile(sd, imgFilePath+img.Filename)
}

// RetreiveImg reads from a filename and returns an Image or error
func RetreiveImg(filename string) (Image, error) {
	// TODO: this
	return Image{}, nil
}
