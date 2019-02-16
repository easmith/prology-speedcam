package main

import (
	. "./types"
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"golang.org/x/text/encoding/charmap"
	"log"
	"os"
)

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
	// Пропускаем заголовок
	file.Seek(20, 0)

	for i := 0; i < int(header.FirstRows); i++ {
		//for i := 0; i < 10; i++ {

		//first := make([]byte, 10)
		//err := binary.Read(fr, binary.LittleEndian, first)
		first := &First{}
		file.Seek(20+10*int64(i), 0)

		err := binary.Read(file, binary.LittleEndian, first)
		if err != nil {
			log.Fatal(err)
		}

		//if i > 10 {
		//	break
		//}

		points := readPoint(file, first.Position, first.Count)

		if first.Count > 0 && len(points) == int(first.Count) {
			//fmt.Printf("%v\t%v\t%v\n", i, float32(first.Y)/float32(points[0].Lon), float32(first.X)/float32(points[0].Lat))
			fmt.Printf("%v\t%v\t%v\n", i, points[0].Lon, points[0].Lat)
			//fmt.Printf("%v\t%v\n", float32(first.X)/float32(points[0].Lon), float32(first.Y)/float32(points[0].Lat))
		} else {
			fmt.Printf(">>>>>>> %s\n", first)
		}

	}

	//fmt.Printf("max: %s \n", max)
	//fmt.Println(total)
}

func readPoint(file *os.File, p int32, c int16) []Second {
	file.Seek(int64(p), 0)
	res := make([]Second, c)
	for i := int16(0); i < c; i++ {
		second := Second{}
		binary.Read(file, binary.LittleEndian, &second)
		res[i] = second
	}
	return res
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

		if i < 10 || i > (int(header.SecondRows)-10) {
			fmt.Println(second)
			b, _ := decoder.Bytes(second.Comment[:])
			fmt.Println(string(bytes.Trim(b, "\x00")))
		}
	}

	//for t := 1; t < 1000; t++ {
	//	if val, ok := types[t]; ok {
	//		fmt.Printf("%v\t\n", t)
	//		for k, v := range val {
	//			fmt.Printf("\t%v\t%v\t%v\n", t, k, v)
	//		}
	//	}
	//
	//}

}
