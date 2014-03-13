package dna

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"strconv"
	"strings"
	"unicode/utf8"
)

func EncodeDNA(input string) string {
	trits := huffmanEncode(input)
	n := len(trits)
	s2 := base10toBase3str(n)

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
	fmt.Println("s3", s3)
	fmt.Println("len(s4)", len(s4))
	fmt.Println("s4", s4)
	s5 := base3ToDNA(s4)

	//s5 = nucleotides (nt)

	//1.5
	N := len(s5)
	fmt.Println("s5", s5)
	fmt.Println("n", n)
	//ID := "10" // 2 trit identifier for orig input unique to runtime
	//TODO: this will break if s5 < 75
	segment_length := (N / 25) - 3
	fmt.Println("segment len", segment_length)
	segments := make([]string, segment_length)

	index := 0
	for index < segment_length {

		pos := index * 25
		end := pos + 100
		// if end > (len(s5) - 1) {
		// 	end = len(s5) - 1
		// }

		segments[index] = s5[pos:end]
		fmt.Println("f", index, segments[index])
		index++
	}

	for index, segment := range segments {
		if index != 0 && index%2 != 0 {
			segments[index] = ReverseComplement(segments[index])
			fmt.Println("f", index, segments[index])
		}

		i3 := base10toBase3str(index)
		//padd i3 to len(12)
		for len(i3) < 12 {
			i3 = "0" + i3
		}

		//TODO: id should be computed per file and be unqiue
		// per encoding batch
		id := "12"
		fmt.Println("i3", i3)
		//Odd trits non-zero indexed
		p := (int(id[1-1]) + int(i3[1-1]) + int(i3[3-1]) + int(i3[5-1]) + int(i3[7-1]) + int(i3[9-1]) + int(i3[11-1])) % 3
		fmt.Println("p", p)
		ix := id + i3 + strconv.Itoa(p)
		fmt.Println("ix", ix)
		//append dna encoded ix to fi
		//start with last char of existing Fi o get Fi'
		seed, _ := utf8.DecodeLastRuneInString(segments[index])
		fmt.Println("seed", seed)
		ixe := base3ToDNAStart(ix, seed)
		fmt.Println("ix encoded", ixe)
		fiComp := segment + ixe

		fmt.Println("fi'", index, fiComp)
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

		segments[index] = fiComp

	}

	return strings.Join(segments, "")
}

func DecodeDNA(dna string) string {
	if len(dna)%117 != 0 {
		panic("Invalid dna sequence")
	}

	segments := make([]string, len(dna)/117)
	for i := 0; i < len(dna)-117; i += 117 {
		start := i * 117
		segments[i] = dna[start : start+117]

		//check for A|T or C|G and remove

	}

	fmt.Println("len", len(segments))

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
	runes := []rune(dna)
	var result bytes.Buffer
	for i := len(runes) - 1; i >= 0; i -= 1 {
		result.WriteRune(complement[runes[i]])
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
			temp := dnaTable[start][r]
			result.WriteRune(temp)
			prev = temp
			continue
		}

		next := dnaTable[prev][r]
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
func huffmanEncode(input string) string {
	//load dict file
	initializeDict()
	var result bytes.Buffer
	for _, char := range input {
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
