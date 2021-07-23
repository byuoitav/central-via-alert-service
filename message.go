package main

import (
	"fmt"
	"strings"
)

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
			chunks = append(chunks, s[currentStart:i])
			currentLen = 0
			currentStart = i
		}
		currentLen++
	}
	chunks = append(chunks, s[currentStart:])
	return chunks
}

func longWords(s string, maxlength int) (string, string) {
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

	return best, message
}
