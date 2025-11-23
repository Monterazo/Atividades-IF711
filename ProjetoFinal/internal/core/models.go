package core

// ExpressionRequest representa uma requisição de expressão
type ExpressionRequest struct {
	ExpressionID string
	Expression   string
	DeadlineMs   int64
}

// ExpressionResponse representa uma resposta de expressão
type ExpressionResponse struct {
	ExpressionID string
	Result       float64
	Error        *ErrorInfo
}

// OperationRequest representa uma requisição de operação
type OperationRequest struct {
	ExpressionID string
	StepID       string
	Operation    string
	Numbers      []float64
	DeadlineMs   int64
}

// OperationResponse representa uma resposta de operação
type OperationResponse struct {
	ExpressionID string
	StepID       string
	Result       float64
	Error        *ErrorInfo
}

// ErrorInfo representa informações de erro
type ErrorInfo struct {
	Code    string
	Message string
}
