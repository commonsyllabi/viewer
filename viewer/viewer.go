package viewer

import (
	"bufio"
	"fmt"
	"os"
)

func ReadFile() error {
	file, err := os.Open("manifest.imscc")
	defer file.Close()

	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	return nil
}
