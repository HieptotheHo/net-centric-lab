package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func luhnAlgoChecker(numbers string) string {
	splitter := regexp.MustCompile("[ \n\r\t\f]")
	trimmedNumbers := strings.Join(splitter.Split(strings.TrimSpace(numbers), -1), "")
	if len(trimmedNumbers) <= 1 {
		return "invalid length"
	}
	sum := 0
	isSecondNum := false
	for i := len(trimmedNumbers) - 1; i >= 0; i-- {
		char := string(trimmedNumbers[i])

		number, err := strconv.Atoi(char)

		if err != nil {
			fmt.Println("Error during conversion")
			return "wrong input"
		}
		if isSecondNum {
			number *= 2
			if number > 9 {
				number -= 9
			}
		}
		sum += number
		isSecondNum = !isSecondNum
	}
	if sum%10 == 0 {
		return "valid"
	} else {
		return "invalid"
	}
}

func main() {
	fmt.Println("4539    3195     0343     6467: ", luhnAlgoChecker("4539 3195 0343 6467"))
	fmt.Println("2: ", luhnAlgoChecker("2"))
	fmt.Println("8273 1232 7352 0569: ", luhnAlgoChecker("8273 1232 7352 0569"))
}
