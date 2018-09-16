package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"runtime"
)

var fw = bufio.NewWriter(os.Stdout)

type listDir struct {
	root      string
	filter    *regexp.Regexp
	prefix    string
	posfix    string
	printFile bool
	printDirs bool
	readDirFn func(dirname string) ([]os.FileInfo, error)
}

func (par *listDir) process() {
	par.readDir(par.root)
}

func (par *listDir) readDir(pth string) {

	fmt.Println(par)
	if par.printDirs {
		par.print(pth)
	}

	fi, err := par.readDirFn(pth)
	if err != nil {
		return
	}
	for _, f := range fi {
		if f.IsDir() {
			par.readDir(path.Join(pth, f.Name()))
		} else {
			if par.printFile {
				par.print(path.Join(pth, f.Name()))
			}
		}
	}
}

func (par *listDir) print(line string) {
	defer fw.Flush()
	var ok = true
	if par.filter != nil {
		ok = par.filter.MatchString(line)
	}
	if ok {
		line := fmt.Sprintf("%s%s%s\n", par.prefix, line, par.posfix)
		fw.WriteString(line)
		fw.Flush()
	}

}

func main() {
	if len(os.Args) < 2 {
		color1, color2, color3 := "\x1b[0;36m", "\x1b[37;1m", "\x1b[0m"
		if runtime.GOOS == "windows" {
			color1, color2, color3 = "", "", ""
		}
		fmt.Fprintf(os.Stderr, color1+"Busca y lista Archivos y Directorios en un arbol de directorios\n")
		fmt.Fprintf(os.Stderr, "el primer parametro el el path en donde buscar\n")
		fmt.Fprintf(os.Stderr, "el segundo parametro -d -f -a Imprmir (directory/file/all) (Opcional)\n")
		fmt.Fprintf(os.Stderr, "el tercer parametro es una espresion regular para filtrar el listado (Opcional)\n")
		fmt.Fprintf(os.Stderr, "el cuarto es el prefijo a agregar a cada nombre de archivo (Opcional)\n")
		fmt.Fprintf(os.Stderr, "el quinto es el posfijo a agregar a cada nombre de archivo (Opcional)\n")
		fmt.Fprintf(os.Stderr, "\n\n"+color2+"lsf ./ -f \".go$\" \"<\" \">\"\n\n"+color3)
		return
	}

	pt := os.Args[1]

	var par = &listDir{pt, nil, "", "", true, false, ioutil.ReadDir}

	if len(os.Args) > 2 {
		option := os.Args[2]
		if !(option == "-a" || option == "-d" || option == "-f") {
			fmt.Fprintf(os.Stderr, "Opcion Invalida, opciones posibles -f -d -a\n")
			return
		}
		par.printFile = option == "-a" || option == "-f"
		par.printDirs = option == "-a" || option == "-d"
	}

	if len(os.Args) > 3 {
		var err error
		par.filter, err = regexp.Compile(os.Args[3])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Filtro Invalido\n")
			return
		}
	}
	if len(os.Args) > 4 {
		par.prefix = os.Args[4]
	}
	if len(os.Args) > 5 {
		par.posfix = os.Args[5]
	}

	par.process()

}
