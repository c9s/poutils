package main

import (
	"bufio"
	"fmt"
	"io"
	// "bytes"
	"log"
	"os"
)

type POScanner struct {
	Input chan *string
}

func (self *POScanner) Close() {
	self.Input <- nil
}

func (self *POScanner) Run() {
	var line *string
	for {
		line = <-self.Input
		fmt.Print(*line)
	}
}

func ParseFile(filename string, input chan *string) error {
	var err error
	var f *os.File
	f, err = os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	reader := bufio.NewReader(f)

	for err == nil {
		var line string
		if line, err = reader.ReadString('\n'); err != nil {
			break
		}
		// fmt.Print(line)
		input <- &line
	}
	if err == io.EOF {
		return nil
	}
	return err
}

func main() {
	scanner := POScanner{}
	scanner.Input = make(chan *string, 10)
	go scanner.Run()

	err := ParseFile("locale/en/LC_MESSAGES/jifty.po", scanner.Input)
	if err != nil {
		log.Fatal(err)
	}
	scanner.Close()

	/*
		line, isPrefix, err := reader.ReadLine()
		if err != nil {
			log.Fatal(err)
		}
	*/
	/*
		data := make(byte, 1024)
		for ; count , err := f.Read(data) ; err != nil {

		}
	*/

	// reader := bytes.NewReader(f)
	// _ = reader
}
