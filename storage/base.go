package storage

import (
	"encoding/json"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type NFTStorage interface {
	Write(data interface{}, output string) error
	WriteImage(data *image.RGBA, output string) error
	Read(path string) ([]byte, error)
	ReadImage(path string) (image.Image, error)
}

type LocalStorage struct{}

func (LocalStorage) Write(data interface{}, output string) error {
	bytes, _ := json.Marshal(data)
	err := ioutil.WriteFile(output, bytes, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (LocalStorage) WriteImage(data *image.RGBA, output string) error {
	file, err := os.Create(output)
	log.Println("create file:", err)
	if err != nil {
		return err
	}
	return png.Encode(file, data)
}

func (LocalStorage) Read(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return io.ReadAll(file)
}

func (LocalStorage) ReadImage(path string) (image.Image, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	imgFile, err := os.Open(absPath)
	if err != nil {
		return nil, err
	}
	defer imgFile.Close()

	var img image.Image
	splittedPath := strings.Split(path, ".")
	ext := splittedPath[len(splittedPath)-1]

	if ext == "jpg" || ext == "jpeg" {
		img, err = jpeg.Decode(imgFile)
	} else {
		img, err = png.Decode(imgFile)
	}

	if err != nil {
		return nil, err
	}

	return img, nil
}