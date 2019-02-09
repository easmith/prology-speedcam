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

func main() {
	file, err := os.Open("examples/speedcam.bin")
	defer file.Close()

	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	header := readHeader(file)
	fmt.Println(header)

	//readFirst(header, file)

	//readSecond(header, file)

}

func readHeader(file *os.File) (header *Header) {
	fr := bufio.NewReader(file)
	header = &Header{}
	err := binary.Read(fr, binary.LittleEndian, header)
	if err != nil {
		log.Fatal(err)
	}
	return
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

func readFirst(header *Header, file *os.File) {
	bufLen := 10
	buf := make([]byte, bufLen)
	ui := int64(20)
	l := int32(0)
	file.Seek(20, 0)
	for {
		_, err := file.Read(buf)

		if err != nil {
			log.Fatal(err)
		}
		ui += int64(bufLen)
		_, err = file.Seek(ui, 0)
		if err != nil {
			log.Fatal(err)
		}

		l++
		//if l > (header.UnknownDB - 15000) && l < (header.UnknownDB - 5) {
		//	fmt.Printf("%s %s %s\n",  hex.EncodeToString(buf[0:4]), hex.EncodeToString(buf[4:6]), hex.EncodeToString(buf[6:10]))
		//	fmt.Printf("%s %s\n", hex.EncodeToString(buf[0:5]), hex.EncodeToString(buf[5:10]))
		//fmt.Printf("%s %s\n", ui, hex.EncodeToString(buf))

		//var i16 int32
		//fmt.Printf("[%v %v %v]\n",
		//	binary.LittleEndian.Uint16(buf[0:2]),
		//
		//	binary.LittleEndian.Uint32(buf[2:6]),
		//	binary.LittleEndian.Uint32(buf[6:10]),
		//)
		//fmt.Printf("[%v %v %v]\n",
		//
		//	binary.LittleEndian.Uint32(buf[0:4]),
		//	binary.LittleEndian.Uint32(buf[4:8]),
		//
		//	binary.LittleEndian.Uint16(buf[8:10]),
		//)

		fmt.Printf("%v\t%v\t%v",
			binary.LittleEndian.Uint32(buf[0:4]),
			binary.LittleEndian.Uint32(buf[6:10]),
			binary.LittleEndian.Uint16(buf[4:6]),
		)

		fmt.Println("")

		//}

		if l > (header.UnknownDB - 2) {
			break
		}

		if ui > int64(bufLen)*20000000 {
			file.Seek(0, 0)
			break
		}
	}
}
