package main

import "fmt"

func main() {
	var line string

	names := []string{"Sammy", "Jessica", "Drew", "Jamie"}

	line = join(",", names...)
	fmt.Println(line)
}

func join(del string, values ...string) string {
	var line string
	fmt.Println(values)
	for i, v := range values {
		line = line + v
		if i != len(values)-1 {
			line = line + del
		}
	}
	return line
}
