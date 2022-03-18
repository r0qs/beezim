package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

type inputReader interface {
	ReadInput() (input string, err error)
}

func readInput() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	input = strings.Replace(input, "\n", "", -1)
	return input, err
}

func terminalPrompt(r inputReader, msg string) (input string, err error) {
	fmt.Print(msg + ": ")
	input, err = readInput()
	if err != nil {
		return "", err
	}
	return input, nil
}

var AbortCmd error = errors.New("abort")

type ConfirmationInputReader struct {
	action string
	do     func() error
}

func NewConfirmationInputReader(action string, do func() error) *ConfirmationInputReader {
	return &ConfirmationInputReader{
		action: action,
		do:     do,
	}
}

func (c ConfirmationInputReader) ReadInput() (string, error) {
	msg := fmt.Sprintf("%s (y/[n] or q to abort)", c.action)
	answer, err := terminalPrompt(c, msg)
	if err != nil {
		return "", err
	}

	switch answer {
	case "y", "yes":
		return answer, c.do()
	case "q", "skip", "abort":
		return answer, AbortCmd
	default:
	}
	return answer, nil
}
