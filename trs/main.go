package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func main() {

	/*
		validID := regexp.MustCompile(`^[a-z]+\[[0-9]+\]$`)
		fmt.Println(validID.MatchString("adam[23]"))
		fmt.Println(validID.MatchString("eve[7]"))
		fmt.Println(validID.MatchString("Job[48]"))
		fmt.Println(validID.MatchString("snakey"))
	*/

	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Reemplaza un texto del stdin y envia al stdout\n")
		fmt.Fprintf(os.Stderr, "Debe ejcutar\n\ntrs \"textoorigen\" \"textodestino\"\n")
		return
	}
	scanner := bufio.NewScanner(os.Stdin)

	_, err := regexp.Compile(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "RegExp from Invalida")
		return
	}

	// cont := 0
	for scanner.Scan() {
		// cont++
		var line string
		line = scanner.Text() + "\n"
		result := replaceLine(line, os.Args[1], os.Args[2])
		fmt.Print(result)
	}

}

func replaceLine(line string, from string, to string) string {
	to = strings.Replace(to, "\\n", "\n", -1)
	to = strings.Replace(to, "\\r", "\r", -1)
	to = strings.Replace(to, "\\t", "\t", -1)
	rfrom, _ := regexp.Compile(from)
	ret := rfrom.ReplaceAllString(line, to)
	return ret
}
