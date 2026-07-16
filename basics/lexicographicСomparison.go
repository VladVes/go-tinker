package basics

import "fmt"

func getRune(word []rune, index int) rune {
	if index > len(word)-1 {
		return ' '
	}
	return word[index]
}

func appendRune(runes []rune, r rune) []rune {
	if r == 32 {
		return runes
	}
	return append(runes, r)
}

func StrStepByStepComparsionDemo() {
	wordA := []rune("helloworld")
	wordB := []rune("приве")
	lenA := len(wordA)
	lenB := len(wordB)
	maxLen := max(lenA, lenB)
	fmt.Println(maxLen)
	cmpA := []rune{}
	cmpB := []rune{}

	for i := range maxLen {
		runeA := getRune(wordA, i)
		runeB := getRune(wordB, i)
		fmt.Printf("word A %d:%v\n", i, runeA)
		fmt.Printf("word B %d:%v\n", i, runeB)
		cmpA = appendRune(cmpA, runeA)
		cmpB = appendRune(cmpB, runeB)
		fmt.Printf("%c == %c:%v\n", cmpA, cmpB, string(cmpA) == string(cmpB))
		fmt.Printf("%c > %c:%v\n", cmpA, cmpB, string(cmpA) > string(cmpB))
		fmt.Printf("%c < %c:%v\n", cmpA, cmpB, string(cmpA) < string(cmpB))

	}
}
