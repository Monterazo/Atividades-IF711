package main

import (
	"fmt"
	"github.com/Monterazo/Atividades-IF711/ProjetoFinal/internal/core"
)

func main() {
	parser := core.NewParser()

	testCases := []struct {
		expr     string
		expected string
	}{
		{"10+20*3", "Esperado: 70 (10 + (20*3))"},
		{"((4+3)*2)/5", "Esperado: 2.8 ((7*2)/5)"},
		{"100/4+50*2", "Esperado: 125 ((100/4) + (50*2))"},
		{"2+3*4-5", "Esperado: 9 (2 + (3*4) - 5)"},
	}

	fmt.Println("===========================================")
	fmt.Println("Teste do Parser - Verificando Precedencia")
	fmt.Println("===========================================")

	for _, tc := range testCases {
		fmt.Printf("\nExpressao: %s\n", tc.expr)
		fmt.Printf("%s\n", tc.expected)

		steps, err := parser.Parse(tc.expr)
		if err != nil {
			fmt.Printf("ERRO: %v\n", err)
			continue
		}

		fmt.Printf("Steps gerados:\n")
		for i, step := range steps {
			fmt.Printf("  %d. %s(%v) -> %s\n", i+1, step.Operation, step.Numbers, step.ID)
		}
	}
}
