package rabbitmq

// ExpressionRequest representa uma requisição de expressão via RabbitMQ
type ExpressionRequest struct {
	ExpressionID string `json:"expression_id"`
	Expression   string `json:"expression"`
	DeadlineMs   int64  `json:"deadline_ms"`
}

// ExpressionResponse representa uma resposta de expressão via RabbitMQ
type ExpressionResponse struct {
	ExpressionID string     `json:"expression_id"`
	Result       float64    `json:"result"`
	Error        *ErrorInfo `json:"error,omitempty"`
}

// OperationRequest representa uma requisição de operação via RabbitMQ
type OperationRequest struct {
	ExpressionID string    `json:"expression_id"`
	StepID       string    `json:"step_id"`
	Operation    string    `json:"operation"`
	Numbers      []float64 `json:"numbers"`
	DeadlineMs   int64     `json:"deadline_ms"`
}

// OperationResponse representa uma resposta de operação via RabbitMQ
type OperationResponse struct {
	ExpressionID string     `json:"expression_id"`
	StepID       string     `json:"step_id"`
	Result       float64    `json:"result"`
	Error        *ErrorInfo `json:"error,omitempty"`
}

// ErrorInfo representa informações de erro
type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
