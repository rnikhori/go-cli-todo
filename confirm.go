package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

var stdin io.Reader = os.Stdin

func confirm(prompt string) (bool, error) {
	fmt.Printf("%s [y/N]: ", prompt)
	line, err := bufio.NewReader(stdin).ReadString('\n')
	if err != nil && err != io.EOF {
		return false, err
	}
	answer := strings.ToLower(strings.TrimSpace(line))
	return answer == "y" || answer == "yes", nil
}
