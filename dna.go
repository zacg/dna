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

	for index, _ := range segments {
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
		fiComp := segments[index] + ixe

		fmt.Println("fi'", index, fiComp)
		//1.9 prepend AT and append CG to mark

		if fiComp[0] == 'A' {
			fiComp = "T" + fiComp
		} else if fiComp[0] == 'T' {
			fiComp = "A" + fiComp
		} else {
			//choose at random
			fmt.Println("choosing randmon AT")
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
			fmt.Println("choosing random GC")
			//choose at random
			if rand.Intn(1) == 0 {
				fiComp += "G"
			} else {
				fiComp += "C"
			}
		}

		fmt.Println("adjusted ficomp", fiComp)

		segments[index] = fiComp

	}

	return strings.Join(segments, "")
}

var dnaTritTbl = map[byte]map[byte]byte{
	'A': {'T': '0', 'G': '1', 'C': '2'},
	'C': {'A': '0', 'T': '1', 'G': '2'},
	'G': {'C': '0', 'A': '1', 'T': '2'},
	'T': {'G': '0', 'C': '1', 'A': '2'},
}

func DecodeDNA(dna string) string {
	if len(dna)%117 != 0 {
		panic("Invalid dna sequence")
	}

	segments := make([]string, len(dna)/117)
	for i := 0; i < len(segments); i++ {
		fmt.Println("loop", i)
		start := i * 117
		end := start + 117

		fmt.Println("start", start)
		fmt.Println("end", end)
		segment := dna[start:end]
		fmt.Println("cut seg", segment)

		//check for reverse complement
		if segment[0] != 'T' && segment[0] != 'A' {
			segment = ReverseComplement(segment)
		}

		//Trim A|T or C|G
		segment = segment[1:116]

		//recreate ix and fi
		ix := segment[len(segment)-15:]
		Fi := segment[:len(segment)-15]

		fmt.Println("ix %v len %v", ix, len(ix))
		fmt.Println("Fi %v len %v", Fi, len(Fi))

		// # Convert ix to trits (IX)
		lastFi := Fi[len(Fi)-1]

		IX := string(dnaTritTbl[ix[0]][lastFi])
		for x := 1; x < 15; x++ {
			IX += string(dnaTritTbl[ix[x]][ix[x-1]])
		}

		fmt.Println("IX", IX)

		//Extract ID
		ID := IX[:2]

		fmt.Println("ID", ID)

		//#Extract i3 and i
		i3 := IX[2 : len(IX)-1]
		extractedI := base3toBase10(i3)

		fmt.Println("i", extractedI)

		//parity check
		//temp := strconv.Itoa(IX[len(IX)-1])
		P, _ := strconv.Atoi(string(IX[len(IX)-1]))
		//P := utf8.DecodeRuneInString(IX[len(IX)-1])
		Pexpected := (int(ID[1-1]) + int(i3[1-1]) + int(i3[3-1]) +
			int(i3[5-1]) + int(i3[7-1]) + int(i3[9-1]) + int(i3[11-1])) % 3

		fmt.Println("P", P)
		fmt.Println("Pexpected", Pexpected)

		if P != Pexpected {
			//panic("Corrupted segment:\nID = %s\ni = %d")
			panic("corrupt segment " + strconv.Itoa(P) + " " + strconv.Itoa(Pexpected))
		} else {
			//reverse complement odd fi
			if extractedI%2 == 1 {
				segment = ReverseComplement(Fi)
				fmt.Println("segment '", segment)
			} else {
				segment = Fi
			}
		}

		fmt.Println("segment output", segment)
		segments[i] = segment
	}

	fmt.Println("len", len(segments))

	//process back to s0
	s5 := fiToS5(segments)
	s4 := s5Tos4(s5)
	s0 := s4Tos0(s4)
	return s0
}

func fiToS5(fi []string) string {
	fmt.Println("len(fi)", len(fi))
	fmt.Println("len(fi)[0]", len(fi[1]))
	s5 := fi[0][0:75]
	for _, segment := range fi {
		s5 += segment[len(segment)-25:]
	}
	return s5
}

func s5Tos4(s5 string) string {

	bytes := []byte(s5)
	s4 := make([]byte, len(s5)+1)
	for x := len(s5) - 1; x > 1; x-- {
		s4[x] = dnaTritTbl[bytes[x]][bytes[x-1]]
	}
	s4[1] = dnaTritTbl[bytes[1]][bytes[len(bytes)-1]]
	s4[0] = dnaTritTbl[bytes[0]]['A']

	fmt.Println("s4", s4)
	return string(s4)
}

func s4Tos0(s4 string) string {
	//last 20 trits = s2
	s2 := s4[len(s4)-20:]

	// n = len(s1)
	n := base3toBase10(s2)

	// first n trits make up s1
	s1 := s4[:n]

	//convert trits to data using huffman
	s0 := huffmanDecode(s1)

	return s0
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

	content, err := ioutil.ReadFile("/home/zac/dev/go/src/github.com/zacg/dna/huff3.dict")
	if err != nil {
		panic("io error" + err.Error())
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

func huffmanDecode(input string) string {
	//load dict file
	initializeDict()
	//build inverse dict
	invDict := make(map[string]int)
	for key, value := range hDict {
		invDict[value] = key
	}

	var result bytes.Buffer
	x := 0
	for x < len(input) {
		//result.WriteString(hDict[int(char)])
		if val, ok := invDict[input[x:x+5]]; ok {
			result.WriteByte(byte(val))
			x += 5
		} else {
			result.WriteByte(byte(invDict[input[x:x+6]]))
			x += 6
		}

		//result.Write(invDict[char])
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

func base3toBase10(input string) int {
	//TODO: this can be removed when utf code cleaned up
	input = strings.TrimRight(input, "\x00")

	n, err := strconv.Atoi(input)

	if err != nil {
		fmt.Println("input %s", input)
		fmt.Println(err.Error())
		panic("invalid base3 number")
	}
	if n == 0 {
		return 0
	}

	res := 0
	b := 1
	for n != 0 {
		res = res + (n%10)*b
		n = n / 10
		b = 3 * b
	}

	return res
}
