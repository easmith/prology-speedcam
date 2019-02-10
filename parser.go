package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"golang.org/x/text/encoding/charmap"
	"log"
	"os"
)

type Header struct {
	Version    int32
	FirstRows  int32
	SecondRows int32
	Year       int32
	Issue      int32
}

type First struct {
	Unknown  int32
	Count    int16
	Position int32
}

//ImLon/ImLat/sangle/cdir/ctype/cspeed
type Second struct {
	Lon     int32
	Lat     int32
	Angle   int16
	Dir     int8
	Type    int8
	Speed   int8
	Comment [19]byte
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

	readSecond(header, file)

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

	total := 0

	for i := 0; i < int(header.FirstRows); i++ {
		//for i := 0; i < 10; i++ {

		//first := make([]byte, 10)
		//err := binary.Read(fr, binary.LittleEndian, first)
		first := &First{}
		err := binary.Read(fr, binary.LittleEndian, first)
		if err != nil {
			log.Fatal(err)
		}
		if i%1000 == 0 {
			fmt.Println(first)
		}
		total += int(first.Count)
	}

	fmt.Println(total)
}

func readSecond(header *Header, file *os.File) {
	fr := bufio.NewReader(file)

	// Пропускаем заголовок и первый блок
	file.Seek(20+int64(header.FirstRows)*10, 0)

	decoder := charmap.Windows1251.NewDecoder()

	types := make(map[int]map[string]int)

	for i := 0; i < int(header.SecondRows); i++ {
		//for i := 0; i < 10; i++ {

		second := Second{}
		err := binary.Read(fr, binary.LittleEndian, &second)
		if err != nil {
			log.Fatal(err)
		}

		b, _ := decoder.Bytes(second.Comment[:])
		comment := string(bytes.Trim(b, "\x00"))

		if _, ok := types[int(second.Type)]; !ok {
			types[int(second.Type)] = make(map[string]int)
			fmt.Printf("new type: %v\n", second.Type)
		}

		types[int(second.Type)][comment]++

		//types[int(second.Type)][comment]++

		//if second.Type == 20 {
		//	fmt.Printf("TYPE: %v", second)
		//}

		//if i%5000 == 0 {
		//	fmt.Println(second)
		//	b, _ := decoder.Bytes(second.Comment[:])
		//	fmt.Println(string(bytes.Trim(b, "\x00")))
		//}
	}

	for t := 1; t < 1000; t++ {
		if val, ok := types[t]; ok {
			fmt.Printf("%v\t\n", t)
			for k, v := range val {
				fmt.Printf("\t%v\t%v\n", k, v)
			}
		}

	}

}
