package main

import (
	"fmt"

	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

func main() {

	bs_UTF16LE, _, _ := transform.Bytes(unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewEncoder(), []byte("1"))
	bs_UTF16BE, _, _ := transform.Bytes(unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM).NewEncoder(), []byte("1"))

	bs_UTF16LEN, _ := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewEncoder().Bytes([]byte("1"))
	bs_UTF16BEN, _ := unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM).NewEncoder().Bytes([]byte("1"))

	bs_UTF8LE, _, _ := transform.Bytes(unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder(), bs_UTF16LE)
	bs_UTF8BE, _, _ := transform.Bytes(unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM).NewDecoder(), bs_UTF16BE)

	fmt.Printf("%v\n%v\n%v\n%v\n%v\n%v\n", bs_UTF16LE, bs_UTF16BE, bs_UTF16LEN, bs_UTF16BEN, bs_UTF8LE, bs_UTF8BE)
}
