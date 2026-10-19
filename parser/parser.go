package parser

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Instruction struct {
	Command string
	Args    []string
	Line    int
}

func ParseScript(filepath string) ([]Instruction, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open script: %w", err)
	}
	defer file.Close()

	var instructions []Instruction
	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		text := scanner.Text()
		trimmed := strings.TrimSpace(text)

		if lineNum == 1 {
			if trimmed != "# moschito:1" {
				return nil, fmt.Errorf("invalid or missing parser directive. First line must be '# moschito:1'")
			}
			continue
		}

		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		parts := strings.SplitN(trimmed, " ", 2)
		command := strings.ToUpper(parts[0])
		var args []string
		
		if len(parts) > 1 {
			if command == "RUN" || command == "SH" || command == "CHECK" {
				args = []string{parts[1]}
			} else {
				args = strings.Fields(parts[1])
			}
		}

		instructions = append(instructions, Instruction{
			Command: command,
			Args:    args,
			Line:    lineNum,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading script: %w", err)
	}

	return instructions, nil
}
