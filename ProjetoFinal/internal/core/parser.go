package core

import (
	"fmt"
	"strconv"
	"strings"
)

// Token representa um token na expressão
type Token struct {
	Type  string // "number", "operator", "paren"
	Value string
}

// StepDependency representa uma dependência de um resultado anterior
type StepDependency struct {
	Position  int    // Posição no array Numbers (0 ou 1)
	Reference string // Referência ao resultado (ex: "result_step0")
}

// Step representa uma operação atômica
type Step struct {
	ID        string
	Operation string
	Numbers   []float64
	DependsOn []StepDependency // Dependências de resultados anteriores
}

// Parser implementa o algoritmo Shunting Yard para converter expressões infix para RPN
type Parser struct{}

// NewParser cria um novo parser
func NewParser() *Parser {
	return &Parser{}
}

// Parse converte uma expressão infix em uma sequência de steps (RPN)
func (p *Parser) Parse(expression string) ([]Step, error) {
	// Tokeniza a expressão
	tokens, err := p.tokenize(expression)
	if err != nil {
		return nil, err
	}

	// Converte para RPN usando Shunting Yard
	rpn, err := p.toRPN(tokens)
	if err != nil {
		return nil, err
	}

	// Converte RPN em steps
	steps := p.rpnToSteps(rpn)
	return steps, nil
}

// ParseWithRPN converte uma expressão infix em steps e retorna também a string RPN
func (p *Parser) ParseWithRPN(expression string) ([]Step, string, error) {
	// Tokeniza a expressão
	tokens, err := p.tokenize(expression)
	if err != nil {
		return nil, "", err
	}

	// Converte para RPN usando Shunting Yard
	rpn, err := p.toRPN(tokens)
	if err != nil {
		return nil, "", err
	}

	// Converte RPN em string para exibição
	var rpnStr string
	for i, token := range rpn {
		if i > 0 {
			rpnStr += " "
		}
		rpnStr += token.Value
	}

	// Converte RPN em steps
	steps := p.rpnToSteps(rpn)
	return steps, rpnStr, nil
}

// tokenize divide a expressão em tokens
func (p *Parser) tokenize(expr string) ([]Token, error) {
	var tokens []Token
	expr = strings.ReplaceAll(expr, " ", "")

	i := 0
	for i < len(expr) {
		ch := expr[i]

		switch {
		case ch >= '0' && ch <= '9' || ch == '.':
			// Número
			j := i
			for j < len(expr) && (expr[j] >= '0' && expr[j] <= '9' || expr[j] == '.') {
				j++
			}
			tokens = append(tokens, Token{Type: "number", Value: expr[i:j]})
			i = j
		case ch == '+' || ch == '-' || ch == '*' || ch == '/':
			// Operador
			tokens = append(tokens, Token{Type: "operator", Value: string(ch)})
			i++
		case ch == '(' || ch == ')':
			// Parênteses
			tokens = append(tokens, Token{Type: "paren", Value: string(ch)})
			i++
		default:
			return nil, fmt.Errorf("caractere inválido: %c", ch)
		}
	}

	return tokens, nil
}

// precedence retorna a precedência de um operador
func (p *Parser) precedence(op string) int {
	switch op {
	case "+", "-":
		return 1
	case "*", "/":
		return 2
	default:
		return 0
	}
}

// toRPN converte tokens infix para RPN usando Shunting Yard
func (p *Parser) toRPN(tokens []Token) ([]Token, error) {
	var output []Token
	var stack []Token

	for _, token := range tokens {
		switch token.Type {
		case "number":
			output = append(output, token)
		case "operator":
			for len(stack) > 0 {
				top := stack[len(stack)-1]
				if top.Type != "operator" {
					break
				}
				if p.precedence(top.Value) < p.precedence(token.Value) {
					break
				}
				output = append(output, top)
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, token)
		case "paren":
			if token.Value == "(" {
				stack = append(stack, token)
			} else {
				found := false
				for len(stack) > 0 {
					top := stack[len(stack)-1]
					stack = stack[:len(stack)-1]
					if top.Value == "(" {
						found = true
						break
					}
					output = append(output, top)
				}
				if !found {
					return nil, fmt.Errorf("parênteses não balanceados")
				}
			}
		}
	}

	for len(stack) > 0 {
		top := stack[len(stack)-1]
		if top.Value == "(" {
			return nil, fmt.Errorf("parênteses não balanceados")
		}
		output = append(output, top)
		stack = stack[:len(stack)-1]
	}

	return output, nil
}

// rpnToSteps converte RPN em steps
func (p *Parser) rpnToSteps(rpn []Token) []Step {
	var steps []Step
	var stack []string
	stepCounter := 0

	for _, token := range rpn {
		if token.Type == "number" {
			stack = append(stack, token.Value)
		} else if token.Type == "operator" {
			if len(stack) < 2 {
				continue
			}

			b := stack[len(stack)-1]
			a := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			// Tenta converter para número, se falhar mantém como string
			num1, err1 := strconv.ParseFloat(a, 64)
			if err1 != nil {
				num1 = 0 // placeholder, será substituído depois
			}

			num2, err2 := strconv.ParseFloat(b, 64)
			if err2 != nil {
				num2 = 0 // placeholder, será substituído depois
			}

			operation := p.operatorToOperation(token.Value)
			stepID := fmt.Sprintf("step%d", stepCounter)
			stepCounter++

			step := Step{
				ID:        stepID,
				Operation: operation,
				Numbers:   []float64{num1, num2},
			}

			// Armazena informação sobre dependências
			if err1 != nil {
				// a é uma referência a resultado anterior
				step.DependsOn = append(step.DependsOn, StepDependency{Position: 0, Reference: a})
			}
			if err2 != nil {
				// b é uma referência a resultado anterior
				step.DependsOn = append(step.DependsOn, StepDependency{Position: 1, Reference: b})
			}

			steps = append(steps, step)

			// Empilha um placeholder para o resultado
			stack = append(stack, fmt.Sprintf("result_%s", stepID))
		}
	}

	return steps
}

// operatorToOperation converte símbolo de operador para nome da operação
func (p *Parser) operatorToOperation(op string) string {
	switch op {
	case "+":
		return "add"
	case "-":
		return "subtract"
	case "*":
		return "multiply"
	case "/":
		return "divide"
	default:
		return ""
	}
}
