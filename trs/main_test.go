package main

import (
	"testing"
)

func Test_replaceLine(t *testing.T) {
	type args struct {
		line string
		from string
		to   string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"Replace_Char_1_linea",
			args{
				"Linea1",
				"a",
				"_",
			},
			"Line_1",
		},
		{
			"Replace_Char_2_lineas",
			args{
				"Linea1\nLinea2",
				"a",
				"_",
			},
			"Line_1\nLine_2",
		},
		{
			"Replace_Empty_Line",
			args{
				"\n",
				"^\\n",
				"",
			},
			"",
		},
		{
			"Replace_Empty_Line_1",
			args{
				"Hola Men\n",
				"^\\n",
				"",
			},
			"Hola Men\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := replaceLine(tt.args.line, tt.args.from, tt.args.to); got != tt.want {
				t.Errorf("replaceLine() = %v, want %v", got, tt.want)
			}
		})
	}
}
