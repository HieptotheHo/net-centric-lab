package main

import "fmt"

type Stack []interface{}

func (s *Stack) Push(item interface{}) {
	*s = append(*s, item)
}

func (s *Stack) Pop() interface{} {
	if len(*s) == 0 {
		return nil
	}
	lastIndex := len(*s) - 1
	item := (*s)[lastIndex]
	*s = (*s)[:lastIndex]
	return item
}

func (s *Stack) Peek() interface{} {
	if len(*s) == 0 {
		return nil
	}
	return (*s)[len(*s)-1]
}

func matchingBracketsCheck(text string) string {
	stack := make(Stack, 0)
	for i := 0; i < len(text); i++ {
		char := string(text[i])
		if char == "{" || char == "(" || char == "[" {
			stack.Push(char)
		} else {
			if char == "}" || char == "]" || char == ")" {
				openBracket := stack.Pop()
				fmt.Print(openBracket)
				if (char == "}" && openBracket != "{") || (char == ")" && openBracket != "(") || (char == "]" && openBracket != "[") {
					return "INVALID BRACKET PLACING!"
				}
			}
		}
	}
	return "ALL BRACKETS MATCHING"
}

func main() {
	fmt.Println("fmt.Println(a.TypeOf(xyz)){[ ]}: ", matchingBracketsCheck("fmt.Println(a.TypeOf(xyz)){[ ]}"))
	fmt.Println("io[({12hro12b:})]", matchingBracketsCheck("io[({12hro12b:})]"))

	fmt.Println("fmt.Println(a.TypeOf(xyz)){[ }: ", matchingBracketsCheck("fmt.Println(a.TypeOf(xyz)){[ }"))
	fmt.Println("io[({12hro12b:}]", matchingBracketsCheck("io[({12hro12b:}]"))

}
