package types

import (
	"bufio"
	"encoding/binary"
	"log"
	"os"
)

type BinReader struct {
	file   *os.File
	header *Header
}

func OpenBinReader(filePath string) BinReader {
	file, _ := os.Open(filePath)
	binReader := BinReader{
		file:   file,
		header: &Header{},
	}
	binReader.readHeader()
	return binReader
}

func (br BinReader) readHeader() *Header {
	br.file.Seek(0, 0)
	fr := bufio.NewReader(br.file)
	err := binary.Read(fr, binary.LittleEndian, br.header)
	if err != nil {
		log.Fatal(err)
	}
	return br.header
}

func (br BinReader) readPointsBySector() []Second {
	points := make([]Second, br.header.SecondRows)

	// Пропускаем заголовок
	br.file.Seek(20, 0)
	p := 0
	for i := 0; i < int(br.header.FirstRows); i++ {
		first := &First{}
		br.file.Seek(20+10*int64(i), 0)

		err := binary.Read(br.file, binary.LittleEndian, first)
		if err != nil {
			log.Fatal(err)
		}
		point := br.readPoint(first.Position, first.Count)
		for _, v := range point {
			points[p] = v
			p++
		}
	}
	return points
}

func (br BinReader) GetPoints() (points []Second) {
	return br.readPointsBySector()
}

func (br BinReader) GetTotal() int32 {
	return br.header.SecondRows
}

func (br BinReader) readPoint(p int32, c int16) []Second {
	br.file.Seek(int64(p), 0)
	res := make([]Second, c)
	for i := int16(0); i < c; i++ {
		second := Second{}
		binary.Read(br.file, binary.LittleEndian, &second)
		res[i] = second
	}
	return res
}
