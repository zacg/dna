package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"strconv"
	"strings"
	"unicode/utf8"
)

func main() {

	fmt.Println(huffmanEncode("test"))
	fmt.Println(EncodeDNA("test"))

}

func strToInt(str string) int {
	result, _ := strconv.Atoi(str)
	return result
}

func EncodeDNA(input string) string {
	trits := huffmanEncode(input)
	n := len(trits)
	s2 := strconv.Itoa(n)

	for len(s2) < 20 {
		s2 = "0" + s2
	}

	if len(s2) > 20 {
		panic("s2 too long")
	}

	s3 := ""
	s4 := trits + s3 + s2
	for len(s4)%25 > 0 && len(s3) <= 24 {
		s3 += "0"
		s4 = trits + s3 + s2
	}

	s5 := base3ToDNA(s4)

	//s5 = nucleotides (nt)

	//1.5
	N := len(s5)
	//ID := "10" // 2 trit identifier for orig input unique to runtime
	//TODO: this will break if s5 < 75
	segment_length := (N / 25) - 3
	fmt.Println("segment len", segment_length)
	segments := make([]string, segment_length)

	index := 0
	for index < segment_length {

		pos := index * 25
		end := pos + 100
		if end > (len(s5) - 1) {
			end = len(s5) - 1
		}

		segments[index] = s5[pos:end]
		index++
	}

	for index, segment := range segments {
		if index%2 != 0 {
			segments[index] = ReverseComplement(segments[index])

		}

		i3 := base10toBase3str(index)
		//padd i3 to len(12)
		for len(i3) < 12 {
			i3 = "0" + i3
		}

		//TODO: id should be computed per file and be unqiue
		// per encoding batch
		id := "10"
		//p := strToInt(id) + int(i3[1])
		p := (strToInt(id) + int(i3[1]) + int(i3[3]) + int(i3[5]) + int(i3[7]) + int(i3[9]) + int(i3[11])) % 3
		ix := id + i3 + strconv.Itoa(p)
		//append dna encoded ix to fi
		//start with last char of existing Fi o get Fi'
		seed, _ := utf8.DecodeLastRuneInString(segment)
		fiComp := segment + base3ToDNAStart(ix, seed)

		//1.9 prepend AT and append CG to mark
		//rand.Seed(42)

		if fiComp[0] == 'A' {
			fiComp = "T" + fiComp
		} else if fiComp[0] == 'T' {
			fiComp = "A" + fiComp
		} else {
			//choose at random
			if rand.Intn(1) == 0 {
				fiComp = "T" + fiComp
			} else {
				fiComp = "A" + fiComp
			}
		}

		if fiComp[len(fiComp)-1] == 'C' {
			fiComp += "G"
		} else if fiComp[len(fiComp)-1] == 'G' {
			fiComp += "C"
		} else {
			//choose at random
			if rand.Intn(1) == 0 {
				fiComp += "G"
			} else {
				fiComp += "C"
			}
		}

	}

	return s5
}

func DecodeDNA(dna string) string {
	return "blah"
}

//Returns reverse complement of specified DNA string
func ReverseComplement(dna string) string {
	complement := map[rune]rune{
		'A': 'T',
		'C': 'G',
		'G': 'C',
		'T': 'A',
	}

	var result bytes.Buffer
	for _, rune := range dna {
		//TODO: this should iterate backwards
		result.WriteRune(complement[rune])
	}

	return result.String()
}

var dnaTable = map[rune]map[rune]rune{
	'A': {'0': 'C', '1': 'G', '2': 'T'},
	'C': {'0': 'G', '1': 'T', '2': 'A'},
	'G': {'0': 'T', '1': 'A', '2': 'C'},
	'T': {'0': 'A', '1': 'C', '2': 'G'},
}

func base3ToDNA(base3 string) string {
	//first trit encoded with "A"
	return base3ToDNAStart(base3, 'A')
}

func base3ToDNAStart(base3 string, start rune) string {
	var result bytes.Buffer
	//first trit encoded with start
	prev := '0'
	for index, r := range base3 {
		if index == 0 {
			fmt.Println("rune", r)
			temp := dnaTable[start][r]
			result.WriteRune(temp)
			prev = temp
			continue
		}

		next := dnaTable[prev][r]
		//fmt.Println("next", next)
		result.WriteRune(next)
		prev = next
	}

	return result.String()
}

var hDict map[int]string

func initializeDict() {
	if hDict != nil {
		return
	}

	content, err := ioutil.ReadFile("/home/zac/dev/go/src/go-dna/huff3.dict")
	if err != nil {
		panic("io error")
	}
	lines := strings.Split(string(content), "\n")
	hDict = make(map[int]string, len(lines))
	for _, element := range lines {
		temp := strings.Split(element, ",")
		if len(temp) < 2 {
			continue
		}
		i, _ := strconv.Atoi(temp[0])
		hDict[i] = temp[1]
	}
}

//Encodes a string to base3 using huffman
//digits referred to as 'trit'
func huffmanEncode(input string) string {
	//load dict file
	initializeDict()
	var result bytes.Buffer
	//var str string = "Hello"
	//u //nicodeCodePoints := []int(str)
	var test string = "test"
	iis := []byte(test)

	for _, char := range iis {
		result.WriteString(hDict[int(char)])
	}

	return result.String()
}

func base10toBase3str(num int) string {
	digits := ""
	for num > 0 {
		digit := num % 3
		digits = strconv.Itoa(digit) + digits
		num = num / 3
	}
	return digits
}
