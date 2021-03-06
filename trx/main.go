package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"runtime"
	"strings"
)

func main() {
	// Valido Argumentos
	if len(os.Args) < 4 {
		color1, color2, color3 := "\x1b[0;36m", "\x1b[37;1m", "\x1b[0m"
		if runtime.GOOS == "windows" {
			color1, color2, color3 = "", "", ""
		}
		fmt.Fprintf(os.Stderr, color1+"Agrega un texto antes o despues de una determinada linea\n")
		fmt.Fprintf(os.Stderr, "el primer parametro es -a , -b , -r, -x (After/Before/Replace Line/Replace Text RexExp)\n")
		fmt.Fprintf(os.Stderr, "el segundo parametro es la expresion regular a buscar\n")
		fmt.Fprintf(os.Stderr, "el tercero y sucesivos parametro es el texto a agregar\n")
		fmt.Fprintf(os.Stderr, "\nEjemplo de reemplado exp regular\n")
		fmt.Fprintf(os.Stderr, "\n"+`echo La Direccion del Server=serverIp y el host es serverIp | trx -x '^(.*?Server=)\bserverIp\b(.*?)$' '${1}192.168.99.100${2}'`)
		fmt.Fprintf(os.Stderr, "\n\n"+color2+"trx [-a,-b,-r,-x] \"regexp\" \"txt1\" ... \"txtn\"\n\n"+color3)
		return
	}

	option := os.Args[1]
	in, err := regexp.Compile(os.Args[2])
	ou := os.Args[3]

	if err != nil {
		fmt.Fprintf(os.Stderr, "Expresion Regular Invalida\n")
		return
	}

	if !(option == "-a" || option == "-b" || option == "-r" || option == "-x") {
		fmt.Fprintf(os.Stderr, "Opcion Invalida, opciones posibles -a -b -r -x\n")
		return
	}

	ou = strings.Join(os.Args[3:], "\n")
	ou = strings.Replace(ou, "\\t", "\t", -1)
	scanner := bufio.NewScanner(os.Stdin)
	line := ""
	for scanner.Scan() {
		print(printLine(line, option, in, ou))
		line = scanner.Text()
	}
	print(printLine(line, option, in, ou))

}

// Imprimo las lineas
func print(s1, s2 string) {
	fmt.Println(s1)
	if s2 != "" {
		fmt.Println(s2)
	}
}

/*
	Genero la linea a imprimir segun el parametro seleccionado -a -b -r -x
	retorna dos lineas, la primera siempre tiene valor y la segunda es opcional segun
	el tipo
*/
func printLine(line string, pos string, in *regexp.Regexp, ou string) (string, string) {
	ok := in.MatchString(line)
	switch {
	case ok && pos == "-a":
		return line, ou // Agrego linea After
	case ok && pos == "-b":
		return ou, line // Agrego linea Before
	case ok && pos == "-r":
		return ou, "" // Reemplazo la Linea entera por el nuevo texto
	case ok && pos == "-x":
		return in.ReplaceAllString(line, ou), "" // Modifico el texto de la linea original con el nuevo segun una expresion regular
	default:
		return line, ""
	}
}
