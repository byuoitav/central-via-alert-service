package message

import (
	"fmt"
	"strings"
)

// A small genius little bit of code found on our friend - StackOverFlow
// Takes and chunks larger things into smaller things
func Chunks(s string, chunkSize int) []string {
	if len(s) == 0 {
		return nil
	}
	if chunkSize >= len(s) {
		return []string{s}
	}
	var chunks []string = make([]string, 0, (len(s)-1)/chunkSize+1)
	currentLen := 0
	currentStart := 0
	for i := range s {
		if currentLen == chunkSize {
			// interate over the string until you hit the limit and cut it at that point.
			chunks = append(chunks, s[currentStart:i])
			currentLen = 0
			currentStart = i
		}
		currentLen++
	}
	chunks = append(chunks, s[currentStart:])
	return chunks
}

// Let's break up large words over a certain amount into smaller chunks and rebuild the message
func LongWords(s string, maxlength int) string {
	var message string
	for _, word := range strings.Split(s, " ") {
		if len(word) > maxlength {
			cs := Chunks(word, maxlength)
			for _, c := range cs {
				message = message + c + " "
			}
		} else {
			message = message + word + " "
		}
	}

	return message
}

// Sometimes the simplest logic can help get you out of a bind.
// The question though that is peaking on everyones mind is can it be done.  The answer is yes but messy.  Maybe if I had more time and more programming experience.
func WordChunks(s string, chunkSize int) []string {
	var message string
	if len(s) == 0 {
		return nil
	}
	if chunkSize >= len(s) {
		return []string{s}
	}
	var messchunks []string = make([]string, 0, (len(s)-1)/chunkSize+1)
	//currentLen := 0
	//currentStart := 0
	for _, word := range strings.Split(s, " ") {
		if len(message+word+" ") >= chunkSize {
			messchunks = append(messchunks, message)
			message = ""
			fmt.Printf("len=%d cap=%d %v\n", len(messchunks), cap(messchunks), messchunks)
		} else if word == "" {
			messchunks = append(messchunks, message)
			fmt.Printf("len=%d cap=%d %v\n", len(messchunks), cap(messchunks), messchunks)
		} else {
			message = message + word + " "
		}
	}
	return messchunks
}
