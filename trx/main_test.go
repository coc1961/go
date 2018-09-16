package main

import (
	"regexp"
	"testing"
)

func Test_printLine(t *testing.T) {
	type args struct {
		line string
		pos  string
		in   *regexp.Regexp
		ou   string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 string
	}{
		{
			// Prueba con dos grupos y un texto
			// (.*?Server=) Grupo 1
			// (.*?) Grupo2
			// \bpapa\b Palabra Prueba
			// Imprimo Grupo1 nueva_palabra grupo2
			// ${1}Prueba${2}
			// Prueba Manual
			//
			// echo Server=papa papa | ./trx -x '^(.*?Server=)\bpapa\b(.*?)$' '${1}Prueba${2}'
			"Test Replace Regexp",
			args{
				"Server=papa papa",
				"-x",
				regexp.MustCompile(`^(.*?Server=)\bpapa\b(.*?)$`),
				`${1}Prueba${2}`,
			},
			"Server=Prueba papa",
			"",
		},
		{
			"Test Replace Line",
			args{
				"Server=papa papa",
				"-r",
				regexp.MustCompile(`Server`),
				`Prueba`,
			},
			"Prueba",
			"",
		},
		{
			"Test Replace Line",
			args{
				"Server=papa papa",
				"-b",
				regexp.MustCompile(`Server`),
				`Prueba`,
			},
			"Prueba",
			"Server=papa papa",
		},
		{
			"Test Replace Line",
			args{
				"Server=papa papa",
				"-a",
				regexp.MustCompile(`Server`),
				`Prueba`,
			},
			"Server=papa papa",
			"Prueba",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := printLine(tt.args.line, tt.args.pos, tt.args.in, tt.args.ou)
			if got != tt.want {
				t.Errorf("printLine() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("printLine() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
