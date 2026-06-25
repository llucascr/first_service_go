package handler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/llucascr/first_service_go/middleware"
	"github.com/llucascr/first_service_go/model"
	"github.com/llucascr/first_service_go/response"
	"github.com/llucascr/first_service_go/service"
)

type MetaContentHandler struct {
	service *service.MetaContentService
}

func NewMetaContentHandler(srv *service.MetaContentService) *MetaContentHandler {
	return &MetaContentHandler{
		service: srv,
	}
}

func (h *MetaContentHandler) mountHandler(r *mux.Router) {
	r.HandleFunc("/metacontent", h.createMetaContent).Methods(http.MethodPost)
	r.HandleFunc("/metacontent/list", h.listMetaContentFromNotebook).Methods(http.MethodGet)
	r.HandleFunc("/metacontent", h.getMetaContentByID).Methods(http.MethodGet)
	r.HandleFunc("/metacontent", h.updateMetaContent).Methods(http.MethodPut)
	r.HandleFunc("/metacontent", h.deleteMetaContent).Methods(http.MethodDelete)
}

func (h *MetaContentHandler) createMetaContent(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := r.Context()

	w.Header().Set("Content-Type", "application/json")

	user, ok := middleware.UserFromContext(ctx)
	if !ok {
		response.Error(w, model.ErrInvalidCredential)
		return
	}

	notebookID, err := uuid.Parse(r.Header.Get("notebook_id"))
	if err != nil {
		http.Error(w, "Erro ao fazer o parse do uuid: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var request model.MetaContentRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "JSON inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	request.UserID = user.UserID    //userID do usuário autenticado (context)
	request.NotebookID = notebookID //notebookID injetado do Header no DTO
	resp, err := h.service.Create(ctx, request)
	if err != nil {
		http.Error(w, "Erro ao salvar meta content: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Falha ao codificar resposta: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *MetaContentHandler) listMetaContentFromNotebook(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := r.Context()

	w.Header().Set("Content-Type", "application/json")

	notebookID, err := uuid.Parse(r.Header.Get("notebook_id"))
	if err != nil {
		http.Error(w, "Erro ao fazer o parse do uuid: "+err.Error(), http.StatusInternalServerError)
		return
	}

	request := model.ListMetaContentsFromNotebookDTO{
		NotebookID: notebookID,
	}
	list_contents, err := h.service.ListMetaContentsFromNotebook(ctx, request)
	if err != nil {
		http.Error(w, "Erro ao procuar lista de Meta Contents do notebook: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(list_contents); err != nil {
		http.Error(w, "Falha ao codificar resposta: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *MetaContentHandler) getMetaContentByID(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := r.Context()

	w.Header().Set("Content-Type", "application/json")

	contentID, err := uuid.Parse(r.Header.Get("content_id"))
	if err != nil {
		http.Error(w, "Erro ao fazer o parse do uuid: "+err.Error(), http.StatusInternalServerError)
		return
	}

	content_found, err := h.service.GetMetaContentByID(ctx, contentID)
	if err != nil {
		http.Error(w, "Erro ao procurar meta content: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(content_found); err != nil {
		http.Error(w, "Falha ao codificar resposta: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *MetaContentHandler) updateMetaContent(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := r.Context()

	w.Header().Set("Content-Type", "application/json")

	contentID, err := uuid.Parse(r.Header.Get("content_id"))
	if err != nil {
		http.Error(w, "Erro ao fazer o parse do uuid: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var request model.UpdateMetaContentDTO
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "JSON inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.service.Update(ctx, contentID, request)
	if err != nil {
		http.Error(w, "Erro ao atualizar meta content: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Falha ao codificar resposta: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *MetaContentHandler) deleteMetaContent(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := r.Context()

	w.Header().Set("Content-Type", "application/json")

	contentID, err := uuid.Parse(r.Header.Get("content_id"))
	if err != nil {
		http.Error(w, "Erro ao fazer o parse do uuid: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := h.service.Delete(ctx, contentID); err != nil {
		http.Error(w, "Erro ao deletar meta content: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
