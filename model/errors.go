package model

import "errors"

// Erros sentinel de domínio. Mensagens em inglês, minúsculas e sem pontuação
// final, seguindo a convenção do Go para error strings. O mapeamento para
// status HTTP fica centralizado em response.Error.
var (
	ErrUserExists        = errors.New("user already exists")  // -> 409
	ErrInvalidCredential = errors.New("invalid credential")   // -> 401
	ErrNotFound          = errors.New("resource not found")   // -> 404
	ErrInvalidBody       = errors.New("invalid request body") // -> 400
)

// AppError carrega um status HTTP e uma mensagem segura para o cliente quando o
// caso não se encaixa num sentinel fixo (ex.: erro de validação com detalhe
// dinâmico). Err guarda o erro interno para log, nunca exposto ao cliente.
type AppError struct {
	Status  int
	Message string
	Err     error
}

func (e *AppError) Error() string { return e.Message }

func (e *AppError) Unwrap() error { return e.Err }

func NewAppError(status int, message string, err error) *AppError {
	return &AppError{Status: status, Message: message, Err: err}
}
