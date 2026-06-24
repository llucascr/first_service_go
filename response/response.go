// Package response centraliza a escrita de respostas HTTP (sucesso e erro).
// É a camada de borda: traduz erros de domínio em status HTTP e garante que
// detalhes internos nunca vazem para o cliente.
package response

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/llucascr/first_service_go/model"
)

// JSON escreve data como JSON com o status informado.
func JSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data == nil {
		return
	}
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("[response] falha ao codificar resposta: %v", err)
	}
}

// Error traduz err em uma resposta HTTP. A precedência é:
//  1. *model.AppError  -> usa Status e Message próprios;
//  2. sentinel de domínio -> status mapeado + a própria mensagem do sentinel;
//  3. qualquer outro erro -> 500, loga o erro interno e expõe mensagem genérica.
func Error(w http.ResponseWriter, err error) {
	var appErr *model.AppError
	if errors.As(err, &appErr) {
		if appErr.Err != nil {
			log.Printf("[response] %d: %v", appErr.Status, appErr.Err)
		}
		writeError(w, appErr.Status, appErr.Message)
		return
	}

	switch {
	case errors.Is(err, model.ErrUserExists):
		writeError(w, http.StatusConflict, err.Error())
	case errors.Is(err, model.ErrInvalidCredential):
		writeError(w, http.StatusUnauthorized, err.Error())
	case errors.Is(err, model.ErrNotFound):
		writeError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, model.ErrInvalidBody):
		writeError(w, http.StatusBadRequest, err.Error())
	default:
		log.Printf("[response] erro interno não tratado: %v", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	JSON(w, status, map[string]string{"error": message})
}
