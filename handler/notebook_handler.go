package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/llucascr/first_service_go/middleware"
	"github.com/llucascr/first_service_go/model"
	"github.com/llucascr/first_service_go/service"
)

type NotebookHandler struct {
	service *service.NoteBookService
}

func NewNotebookHandler(srv *service.NoteBookService) *NotebookHandler {
	return &NotebookHandler{
		service: srv,
	}
}

func (h *NotebookHandler) mountHandler(r *mux.Router) {
	r.HandleFunc("/notebook", h.createNotebook).Methods(http.MethodPost)
	r.HandleFunc("/notebook/list", h.listNotebookFromUser).Methods(http.MethodGet)
	r.HandleFunc("/notebook", h.getNotebookByID).Methods(http.MethodGet)
	r.HandleFunc("/notebook", h.updateNotebook).Methods(http.MethodPut)
	r.HandleFunc("/notebook", h.deleteNotebook).Methods(http.MethodDelete)
}

func (h *NotebookHandler) createNotebook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userRaw := ctx.Value(middleware.UserContextKey)

	user, ok := userRaw.(*model.User)
	if !ok || user == nil {
		http.Error(w, "user not in context", http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var request model.NotebookRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "JSON inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	request.UserID = user.UserID

	resp, err := h.service.Create(ctx, request)
	if err != nil {
		http.Error(w, "Erro ao salvar notebook: "+err.Error(), http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(resp)
}

func (h *NotebookHandler) listNotebookFromUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userRaw := ctx.Value(middleware.UserContextKey)

	w.Header().Set("Content-Type", "application/json")

	user, ok := userRaw.(*model.User)
	if !ok || user == nil {
		http.Error(w, "user not in context", http.StatusForbidden)
		return
	}

	// TODO: fix: corrigir os erros, mudar o jeito que estão aparecendo e os casos que aparecem
	request := model.ListNotebooksFromUserDTO{
		UserID: user.UserID,
	}
	list_notebooks, err := h.service.ListNotebooksFromUser(ctx, request)
	if err != nil {
		http.Error(w, "Erro ao procuar lista de Notebooks do user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(list_notebooks); err != nil {
		http.Error(w, "Falha ao codificar resposta: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *NotebookHandler) getNotebookByID(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := context.TODO()

	w.Header().Set("Content-Type", "application/json")

	notebook_id, err := uuid.Parse(r.Header.Get("notebook_id"))
	if err != nil {
		http.Error(w, "Erro ao fazer o parse do uuid: "+err.Error(), http.StatusInternalServerError)
		return
	}

	notebook_found, err := h.service.GetNotebookByID(ctx, notebook_id)
	if err != nil {
		http.Error(w, "Erro ao procurar notebook: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(notebook_found); err != nil {
		http.Error(w, "Falha ao codificar resposta: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *NotebookHandler) updateNotebook(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := context.TODO()

	w.Header().Set("Content-Type", "application/json")

	notebookID, err := uuid.Parse(r.Header.Get("notebook_id"))
	if err != nil {
		http.Error(w, "Erro ao fazer o parse do uuid: "+err.Error(), http.StatusInternalServerError)
		return
	}

	user_id, err := uuid.Parse(r.Header.Get("user_id"))

	var request model.NotebookRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "JSON inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	request.UserID = user_id
	resp, err := h.service.Update(ctx, notebookID, request)
	if err != nil {
		http.Error(w, "Erro ao atualizar notebook: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Falha ao codificar resposta: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *NotebookHandler) deleteNotebook(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := context.TODO()

	w.Header().Set("Content-Type", "application/json")

	notebookID, err := uuid.Parse(r.Header.Get("notebook_id"))
	if err != nil {
		http.Error(w, "Erro ao fazer o parse do uuid: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := h.service.Delete(ctx, notebookID); err != nil {
		http.Error(w, "Erro ao deletar notebook: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
