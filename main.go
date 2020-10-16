package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	var (
		input  string
		tweets []string
		err    error
	)

	// checking args
	args := os.Args[1:]
	if len(args) >= 1 {
		input = strings.Join(os.Args[1:], " ")
	} else {
		Instructions()
		return
	}

	tweets, err = GenerateTweets(input)
	if err != nil {
		log.Panic(fmt.Errorf("generateTweets(): %w", err))
	}

	PrintTweets(tweets)
}

// constant with tweet lenght
const maxLenght = 140

//GenerateTweets split text into an 140 characters string array and returns a string slice with tweets and index
func GenerateTweets(input string) (tweets []string, err error) {

	var (
		inputLenght, i int
	)

	inputLenght = len(input)

	if inputLenght > maxLenght {

		//string reader
		reader := strings.NewReader(input)

		tweetSize, maxIndexSize := MapInput(inputLenght)

		readerOffSet := maxLenght - maxIndexSize

		for i = 0; i < tweetSize; i++ {

			var (
				textPart []byte
			)

			// creates a prefix with index of text
			index := strconv.Itoa(i+1) + "/" + strconv.Itoa(tweetSize) + " "

			// creates a reader offset
			offset, err := reader.Seek(int64(i*readerOffSet), 0)
			if err != nil {
				return nil, fmt.Errorf("reader.Seek(): %w", err)
			}

			// validates the index to set text
			if int(offset)+readerOffSet > inputLenght {
				textPart = make([]byte, int64(inputLenght)-offset)
			} else {
				textPart = make([]byte, readerOffSet)
			}

			// set io to read at least 1 char
			if read, err := io.ReadAtLeast(reader, textPart, 1); err != nil && read == 0 {
				return nil, fmt.Errorf("io.ReadAtLeast(): %w", err)
			}

			// concatenates index and build the final tweet
			tweet := index + string(textPart)

			// append the tweet into an array of string
			tweets = append(tweets, tweet)
		}
	} else {
		tweets = append(tweets, input)
	}

	return
}

//MapInput used to separate in quantity of tweets with chars limit and then returns infos for prefixing each tweet
func MapInput(inputLenght int) (tweetSize int, maxIndexSize int) {

	//calculate total tweets from inputLenght
	tweetSize = inputLenght / maxLenght

	//multiply total of chars by 2 considering prefixed text with index
	maxIndexSize = (len(strconv.Itoa(tweetSize)) * 2) + 2

	//recalculate tweetSize considering prefixed text for each tweet
	tweetSize = (inputLenght + maxIndexSize*tweetSize) / maxLenght

	if inputLenght%maxLenght != 0 {
		tweetSize++
	}

	return
}

//PrintTweets print tweets on a tweet slice, from the last to the first part
func PrintTweets(tweets []string) {
	for i := len(tweets) - 1; i >= 0; i-- {
		fmt.Println(tweets[i])
	}
}

//Instructions ...
func Instructions() {
	fmt.Println("Use the app by executing the program followed by your input text --> go run main.go \"the text\"")
}
