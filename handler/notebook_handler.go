package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
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

func (h *NotebookHandler) MountNotebookHandler(r *mux.Router) {
	r.HandleFunc("/notebook", h.CreateNotebook).Methods(http.MethodPost)
	r.HandleFunc("/notebook/list", h.ListNotebookFromUser).Methods(http.MethodGet)
	r.HandleFunc("/notebook", h.GetNotebookByID).Methods(http.MethodGet)
	r.HandleFunc("/notebook", h.UpdateNotebook).Methods(http.MethodPut)
	r.HandleFunc("/notebook", h.DeleteNotebook).Methods(http.MethodDelete)
}

func (h *NotebookHandler) CreateNotebook(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := context.TODO()

	w.Header().Set("Content-Type", "application/json")

	userID, err := uuid.Parse(r.Header.Get("user_id"))
	if err != nil {
		http.Error(w, "Erro ao fazer o parse do uuid: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var request model.NotebookRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "JSON inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	request.UserID = userID //userID injetado do Header no DTO
	resp, err := h.service.Create(ctx, request)
	if err != nil {
		http.Error(w, "Erro ao salvar notebook: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Falha ao codificar resposta: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *NotebookHandler) ListNotebookFromUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := context.TODO()

	w.Header().Set("Content-Type", "application/json")

	userID, err := uuid.Parse(r.Header.Get("user_id"))
	if err != nil {
		http.Error(w, "Erro ao fazer o parse do uuid: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO: fix: corrigir os erros, mudar o jeito que estão aparecendo e os casos que aparecem
	request := model.ListNotebooksFromUserDTO{
		UserID: userID,
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

func (h *NotebookHandler) GetNotebookByID(w http.ResponseWriter, r *http.Request) {
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

func (h *NotebookHandler) UpdateNotebook(w http.ResponseWriter, r *http.Request) {
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

func (h *NotebookHandler) DeleteNotebook(w http.ResponseWriter, r *http.Request) {
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