package assets

import (
	"fmt"

	rice "github.com/GeertJohan/go.rice"
)

func loadFromRice(filename, description string) ([]byte, error) {
	box, err := riceBox()
	if err != nil {
		return nil, err
	}
	data, err := box.Bytes(filename)
	if err != nil {
		return nil, fmt.Errorf("could not load %s: %s", description, err)
	}
	return data, nil

}

func riceBox() (*rice.Box, error) {
	box, err := rice.FindBox("")
	if err != nil {
		return nil, fmt.Errorf("rice find box failed: %s", err)
	}
	return box, nil
}

// Scissors returns the scissors image as bytes.
func Scissors() ([]byte, error) {
	return loadFromRice("scissors.png", "scissors image")
}

// CHCross returns the CH-cross image as bytes.
func CHCross() ([]byte, error) {
	return loadFromRice("ch-cross.png", "ch-cross image")
}
