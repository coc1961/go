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
		fmt.Fprintf(os.Stderr, "el primer parametro es -a , -b , -r, -x (After/Before/Replace Line/Replace RexExp)\n")
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

	if !(pos == "-a" || pos == "-b" || pos == "-r" || pos == "-x") {
		fmt.Fprintf(os.Stderr, "Opcion Invalida, opciones posibles -a -b -r -x\n")
		return
	}

	ou = strings.Join(os.Args[3:], "\n")
	ou = strings.Replace(ou, "\\t", "\t", -1)
	scanner := bufio.NewScanner(os.Stdin)
	line := ""
	for scanner.Scan() {
		print(printLine(line, pos, in, ou))
		line = scanner.Text()
	}
	print(printLine(line, pos, in, ou))

}

func print(s1, s2 string) {
	fmt.Println(s1)
	if s2 != "" {
		fmt.Println(s2)
	}
}
func printLine(line string, pos string, in *regexp.Regexp, ou string) (string, string) {
	ok := in.MatchString(line)
	switch {
	case ok && pos == "-a":
		return line, ou
	case ok && pos == "-b":
		return ou, line
	case ok && pos == "-r":
		return ou, ""
	case ok && pos == "-x":
		return in.ReplaceAllString(line, ou), ""
	default:
		return line, ""
	}
}
