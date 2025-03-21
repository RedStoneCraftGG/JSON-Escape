package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func decode(inputStr string) (string, error) {
	var decoded string
	if err := json.Unmarshal([]byte(fmt.Sprintf(`"%s"`, inputStr)), &decoded); err != nil {
		return inputStr, nil
	}
	return decoded, nil
}

func removeComments(jsonStr string) string {
	var result []string
	lines := strings.Split(jsonStr, "\n")
	inMultilineComment := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			result = append(result, line)
			continue
		}

		if inMultilineComment {
			if strings.Contains(line, "*/") {
				inMultilineComment = false
				line = line[strings.Index(line, "*/")+2:]
			}
			continue
		}

		if strings.Contains(line, "//") {
			line = line[:strings.Index(line, "//")]
		}

		if strings.Contains(line, "/*") {
			inMultilineComment = true
			line = line[:strings.Index(line, "/*")]
		}

		result = append(result, line)
	}

	return strings.Join(result, "\n")
}

func encode(inputStr string) string {
	var result strings.Builder
	inString := false
	wordBuffer := ""
	inNumber := false
	runes := []rune(inputStr)

	for i := 0; i < len(runes); i++ {
		char := runes[i]

		if !inString {
			if char == '/' && i+1 < len(runes) {
				nextChar := runes[i+1]
				if nextChar == '/' {
					for i < len(runes) && runes[i] != '\n' {
						i++
					}
					continue
				} else if nextChar == '*' {
					for i < len(runes)-1 && !(runes[i] == '*' && runes[i+1] == '/') {
						i++
					}
					if i < len(runes) {
						i += 2
					}
					continue
				}
			}
		}

		if char == '"' && (i == 0 || runes[i-1] != '\\') {
			inString = !inString
			wordBuffer = ""
			inNumber = false
			result.WriteRune(char)
			continue
		}

		if inString {
			if char == '\\' && i+1 < len(runes) {
				nextChar := runes[i+1]
				if strings.ContainsRune(`"\\/bfnrt`, nextChar) {
					result.WriteRune(char)
					result.WriteRune(nextChar)
					i++
					continue
				}
			}
			result.WriteString(fmt.Sprintf("\\u%04x", char))
			continue
		}

		if strings.ContainsRune("{}[],: ", char) || char <= ' ' {
			wordBuffer = ""
			inNumber = false
			result.WriteRune(char)
			continue
		}

		if char == '-' && i+1 < len(runes) && strings.ContainsRune("0123456789", runes[i+1]) {
			inNumber = true
			result.WriteRune(char)
			continue
		}

		if strings.ContainsRune("0123456789.", char) {
			if !inNumber && strings.ContainsRune("0123456789", char) {
				inNumber = true
			}
			if inNumber {
				result.WriteRune(char)
				continue
			}
		}

		if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') {
			wordBuffer += string(char)
			if wordBuffer == "true" || wordBuffer == "false" || wordBuffer == "null" {
				result.WriteString(wordBuffer)
				wordBuffer = ""
				continue
			}
			if len(wordBuffer) > 4 {
				result.WriteString(wordBuffer)
				wordBuffer = ""
				continue
			}
			continue
		}

		result.WriteString(fmt.Sprintf("\\u%04x", char))
	}

	return result.String()
}

func processEncode(inputPath, outputPath string) error {
	data, err := ioutil.ReadFile(inputPath)
	if err != nil {
		return err
	}

	content := string(data)
	cleanJSON := removeComments(content)
	output := encode(cleanJSON)

	return ioutil.WriteFile(outputPath, []byte(output), 0644)
}

func processDecode(inputPath, outputPath string) error {
	data, err := ioutil.ReadFile(inputPath)
	if err != nil {
		return err
	}

	content := string(data)
	cleanJSON := removeComments(content)

	var parsedData interface{}
	if err := json.Unmarshal([]byte(cleanJSON), &parsedData); err != nil {
		return err
	}

	outputData, err := json.MarshalIndent(parsedData, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(outputPath, outputData, 0644)
}

func main() {
	mode := flag.String("mode", "", "Operation mode: encode/decode (alias: e/d)")
	inputPath := flag.String("input", "", "Input file path (alias: i)")
	outputPath := flag.String("output", "", "Output file path (alias: o)")
	flag.StringVar(mode, "m", "", "Operation mode: encode/decode (alias: e/d)")
	flag.StringVar(inputPath, "i", "", "Input file path")
	flag.StringVar(outputPath, "o", "", "Output file path")
	flag.Parse()

	switch *mode {
	case "e":
		*mode = "encode"
	case "d":
		*mode = "decode"
	}

	if *mode != "encode" && *mode != "decode" {
		fmt.Println("Error: Mode must be 'encode' or 'decode' (or use alias 'e'/'d')")
		os.Exit(1)
	}

	if *inputPath == "" || *outputPath == "" {
		fmt.Println("Error: Input and output paths are required")
		os.Exit(1)
	}

	var err error
	if *mode == "encode" {
		err = processEncode(*inputPath, *outputPath)
	} else {
		err = processDecode(*inputPath, *outputPath)
	}

	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	fmt.Println("Output:", *outputPath)
}