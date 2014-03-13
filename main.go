package main

import (
	"dna"
	"fmt"
)

func main() {
	//str := "Birney and Goldman"
	input := "ATAGTATATCGACTAGTACAGCGTAGCATCTCGCAGCGAGATACGCTGCTACGCAGCATGCTGTGAGTATCGATGACGAGTGACTCTGTACAGTACGTACGATACGTACGTACGTCGTCATCTCGCAGCGAGATACGCTGCTACGCAGCATGCTGTGAGTATCGATGACGAGTGACTCTGTACAGTACGTACGTACGTACGTACGTACGTACGACTATATACGTACGTACGAGC"
	//fmt.Println("Huff", huffmanEncode(str))
	//fmt.Println("Huff2", huffmanEncode("test test test"))
	//fmt.Println("dna", EncodeDNA(str))

	fmt.Println(dna.DecodeDNA(input))
}
