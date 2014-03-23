go-dna
======

golang library for encoding/decoding information in DNA.

The algorithm is based on the method described in this Nature paper: http://www.nature.com/nature/journal/v494/n7435/full/nature11875.html . Pseudo code and details can be found here: http://www.nature.com/nature/journal/vaop/ncurrent/extref/nature11875-s2.pdf , the required huffman table is included in the repository.

Inspired by Allan Costa's python implementation: https://github.com/allanino/DNA

Example:

Encoding:
	str := "some string to encode in DNA"
	dna := dna.EncodeDNA(str)
	fmt.Println("Result: ", dna)

Decoding:
	dna := "ATAGTATATCGACTAGTACAGCGTAGCATCTCGCAGCGAGATACGCTGCTACGCAGCATGCTGTGAGTATCGATGACGAGTGACTCTGTACAGTACGTACGATACGTACGTACGTCGTATAGTCGTACGTACGTACGTACGTACGTACGTACTGTACAGAGTCACTCGTCATCGATACTCACAGCATGCTGCGTAGCAGCGTATCTCGCTGCGAGATGATACGTACGTACGAGC"

	str := dna.DecodeDNA(dna)
	fmt.Println("Result",str)