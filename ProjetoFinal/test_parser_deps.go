package main

import (
	"fmt"
	"github.com/Monterazo/Atividades-IF711/ProjetoFinal/internal/core"
)

func main() {
	parser := core.NewParser()

	expr := "10+20*3"
	fmt.Println("===========================================")
	fmt.Printf("Testando: %s\n", expr)
	fmt.Println("Esperado: 70")
	fmt.Println("===========================================")

	steps, err := parser.Parse(expr)
	if err != nil {
		fmt.Printf("ERRO: %v\n", err)
		return
	}

	fmt.Printf("\nSteps gerados:\n")
	for i, step := range steps {
		fmt.Printf("\nStep %d: %s\n", i, step.ID)
		fmt.Printf("  Operacao: %s\n", step.Operation)
		fmt.Printf("  Numbers: %v\n", step.Numbers)
		if len(step.DependsOn) > 0 {
			fmt.Printf("  Dependencias:\n")
			for _, dep := range step.DependsOn {
				fmt.Printf("    - Position %d depende de %s\n", dep.Position, dep.Reference)
			}
		} else {
			fmt.Printf("  Dependencias: nenhuma\n")
		}
	}

	fmt.Println("\n===========================================")
	fmt.Println("Simulando execucao:")
	fmt.Println("===========================================")

	results := make(map[string]float64)

	for i, step := range steps {
		// Copia os números
		numbers := make([]float64, len(step.Numbers))
		copy(numbers, step.Numbers)

		// Substitui dependências
		for _, dep := range step.DependsOn {
			if result, ok := results[dep.Reference]; ok {
				numbers[dep.Position] = result
				fmt.Printf("\nStep %d: Substituindo position %d com resultado de %s = %v\n",
					i, dep.Position, dep.Reference, result)
			}
		}

		// Executa a operação
		var result float64
		switch step.Operation {
		case "add":
			result = numbers[0] + numbers[1]
		case "subtract":
			result = numbers[0] - numbers[1]
		case "multiply":
			result = numbers[0] * numbers[1]
		case "divide":
			result = numbers[0] / numbers[1]
		}

		fmt.Printf("Step %d: %s(%v) = %v\n", i, step.Operation, numbers, result)
		results[fmt.Sprintf("result_%s", step.ID)] = result
	}

	finalResult := results[fmt.Sprintf("result_step%d", len(steps)-1)]
	fmt.Printf("\n===========================================\n")
	fmt.Printf("Resultado Final: %v\n", finalResult)
	if finalResult == 70 {
		fmt.Println("✅ CORRETO!")
	} else {
		fmt.Println("❌ INCORRETO!")
	}
	fmt.Printf("===========================================\n")
}
