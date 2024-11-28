package main

import (
	"fmt"
	"log"
	"os"

	"github.com/pchchv/mp3"
)

func run() error {
	f, err := os.Open("classic.mp3")
	if err != nil {
		return err
	}
	defer f.Close()

	d, err := mp3.NewDecoder(f)
	if err != nil {
		return err
	}

	fmt.Printf("Length: %d[bytes]\n", d.Length())

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
