package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// IMPROVEMENT: Make these functions receive os.Stdin as arguments - like that it becomes trivial to test!

func GetUserInput(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s\n> ", prompt)
	text, _ := reader.ReadString('\n')
	text = strings.TrimRight(text, "\r\n")
	return text
}

func GetUserInputWithPrefix(prompt string, prefix string) (string, string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s\n%s", prompt, prefix)
	text, _ := reader.ReadString('\n')
	text = strings.TrimRight(text, "\r\n")
	return prefix + text, text
}
