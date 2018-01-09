package core

import (
	"io/ioutil"
	"os"

	"github.com/h2non/bimg"
)

const tmpFilePrefix = "kunai"

func Process(filename string, options *Options) (string, error) {
	filePath := FullAssetPath(filename)

	// Read image buffer from file path
	buffer, err := bimg.Read(filePath)
	if err != nil {
		return "", err
	}

	// Construct image DSL object
	image := bimg.NewImage(buffer)

	// Transform image with the options
	resultImage, err := transformImage(image, options)

	// Create a temporary file
	tmpFile, err := ioutil.TempFile(os.TempDir(), tmpFilePrefix)
	if err != nil {
		return "", err
	}

	// Write to temporary file
	err = bimg.Write(tmpFile.Name(), resultImage)
	if err != nil {
		return "", err
	}

	return tmpFile.Name(), nil
}
