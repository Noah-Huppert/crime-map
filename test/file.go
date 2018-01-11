package test

import (
	"bufio"
	"fmt"
	"os"
)

// ReadFile returns a []string of lines contained in the specified file. An
// error is returned if one occurs. Nil on success.
func ReadFile(path string) ([]string, error) {
	// strs is the list of strings the file contains
	strs := []string{}

	// Read file
	file, err := os.Open(path)
	if err != nil {
		return strs, fmt.Errorf("error opening specified test fields "+
			"file: %s", err.Error())
	}

	// Read lines
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		strs = append(strs, scanner.Text())
	}

	if err = scanner.Err(); err != nil {
		return strs, fmt.Errorf("error scanning file lines: %s",
			err.Error())
	}

	// Close file
	if err = file.Close(); err != nil {
		return strs, fmt.Errorf("error closing specified test fields"+
			" file: %s", err.Error())
	}

	// Success
	return strs, nil
}
