package main

import (
	"bufio"
	"fmt"
	"os"
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
	for scanner.Scan() {
		result := strings.Replace(scanner.Text()+"\n", strings.Replace(os.Args[1], "\\n", "\n", -1), os.Args[2], -1)
		fmt.Print(result)
	}

}
