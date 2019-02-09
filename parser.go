package main

import (
	"bufio"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"os"
)

type Header struct {
	Version   int32
	UnknownDB int32
	GpsDB     int32
	Year      int32
	Issue     int32
}

type First struct {
	Unknown  int32
	Count    int16
	Position int32
}

func main() {
	//file, err := os.Open("examples/speedcam20190131.bin")
	file, err := os.Open("examples/speedcam20190114.bin")
	defer file.Close()

	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	//buf := make([]byte, 20)
	//file.Read(buf)
	//
	//fmt.Printf("%s\n", hex.EncodeToString(buf))

	header := readHeader(file)
	fmt.Println(header)

	readFirst(header, file)

	//readSecond(header, file)

}

func readHeader(file *os.File) (header *Header) {
	file.Seek(0, 0)
	fr := bufio.NewReader(file)
	header = &Header{}
	err := binary.Read(fr, binary.LittleEndian, header)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func readFirst(header *Header, file *os.File) {
	fr := bufio.NewReader(file)

	// Пропускаем заголовок
	file.Seek(20, 0)

	//for i := 0; i < int(header.UnknownDB); i++ {
	for i := 0; i < 10; i++ {

		//first := make([]byte, 10)
		//err := binary.Read(fr, binary.LittleEndian, first)
		first := &First{}
		err := binary.Read(fr, binary.LittleEndian, first)
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Println(hex.EncodeToString(first))
		fmt.Println(first)
	}
}

func readSecond(header *Header, file *os.File) {
	position := 20 + int64(header.UnknownDB)*10

	file.Seek(position, 0)

	for bufLen := 12; bufLen < 40; bufLen++ {
		fmt.Printf("buf %v\n", bufLen)
		buf := make([]byte, bufLen)

		i := 0
		for {
			_, err := file.Read(buf)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("%s\n", hex.EncodeToString(buf))
			position += int64(bufLen)
			_, err = file.Seek(position, 0)
			i++
			if i > 5 {
				break
			}
		}
	}
}
