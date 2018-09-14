package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Fprintf(os.Stderr, "\x1b[93;1mAgrega un texto antes o despues de una determinada linea\n")
		fmt.Fprintf(os.Stderr, "el primer parametro es -a , -b , -r (After/Before/Replace Line)\n")
		fmt.Fprintf(os.Stderr, "el segundo parametro es la expresion regular a buscar\n")
		fmt.Fprintf(os.Stderr, "el tercero y sucesivos parametro es el texto a agregar\n")
		fmt.Fprintf(os.Stderr, "\n\x1b[37;1madds [-a,-b] \"regexporigen\" \"textoagregado\" ... \"textoagregado_n\"\x1b[0m\n\n")
		return
	}

	pos := os.Args[1]
	in, err := regexp.Compile(os.Args[2])
	ou := os.Args[3]

	if err != nil {
		fmt.Fprintf(os.Stderr, "Expresion Regular Invalida\n")
		return
	}

	if !(pos == "-a" || pos == "-b" || pos == "-r") {
		fmt.Fprintf(os.Stderr, "Opcion Invalida, opciones posibles -a -b -r\n")
		return
	}

	ou = strings.Join(os.Args[3:], "\n")
	ou = strings.Replace(ou, "\\t", "\t", -1)
	scanner := bufio.NewScanner(os.Stdin)
	line := ""
	for scanner.Scan() {
		printLine(line, pos, in, ou)
		line = scanner.Text()
	}
	printLine(line, pos, in, ou)

}

func printLine(line string, pos string, in *regexp.Regexp, ou string) {
	ok := in.MatchString(line)
	switch {
	case ok && pos == "-a":
		fmt.Println(line)
		fmt.Println(ou)
		return
	case ok && pos == "-b":
		fmt.Println(ou)
		fmt.Println(line)
		return
	case ok && pos == "-r":
		fmt.Println(ou)
		return
	default:
		fmt.Println(line)
	}
}
