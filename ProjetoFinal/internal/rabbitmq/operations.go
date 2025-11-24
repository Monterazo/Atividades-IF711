package rabbitmq

import (
	"errors"
	"fmt"
)

// ExecuteOperation executa uma operação matemática
func ExecuteOperation(operation string, numbers []float64) (float64, error) {
	if len(numbers) != 2 {
		return 0, errors.New("operação requer exatamente 2 números")
	}

	a, b := numbers[0], numbers[1]

	switch operation {
	case "add":
		return a + b, nil
	case "subtract":
		return a - b, nil
	case "multiply":
		return a * b, nil
	case "divide":
		if b == 0 {
			return 0, errors.New("divisão por zero")
		}
		return a / b, nil
	default:
		return 0, fmt.Errorf("operação desconhecida: %s", operation)
	}
}
