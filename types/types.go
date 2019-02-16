package types

import (
	"bytes"
	"fmt"
	"golang.org/x/text/encoding/charmap"
)

type Header struct {
	Version    int32
	FirstRows  int32
	SecondRows int32
	Year       int32
	Issue      int32
}

type First struct {
	X        int16
	Y        int16
	Count    int16
	Position int32
}

func (f First) String() string {
	return fmt.Sprintf("%v\t%v\t%v\t%v", f.X, f.Y, f.Position, f.Count)
}

type Second struct {
	Lon     int32
	Lat     int32
	Angle   int16
	Dir     int8
	Type    int8
	Speed   int8
	Comment [19]byte
}

var decoder = charmap.Windows1251.NewDecoder()

func (s Second) String() string {
	b, _ := decoder.Bytes(s.Comment[:])
	comment := string(bytes.Trim(b, "\x00"))

	return fmt.Sprintf("{%v,%v\t%v\t%v\t%v\t%v\t%s}", float32(s.Lon)/10000, float32(s.Lat)/10000, s.Angle, s.Dir, s.Type, s.Speed, comment)
}

type PointsGroup struct {
	Points []Second
}
