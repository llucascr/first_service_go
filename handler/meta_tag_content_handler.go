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

type MetaTagContentHandler struct {
	service *service.MetaTagContentService
}

func NewMetaTagContentHandler(srv *service.MetaTagContentService) *MetaTagContentHandler {
	return &MetaTagContentHandler{
		service: srv,
	}
}

func (h *MetaTagContentHandler) mountHandler(r *mux.Router) {
	r.HandleFunc("/metatagcontent", h.createMetaTagContent).Methods(http.MethodPost)
	r.HandleFunc("/metatagcontent/list/content", h.listMetaTagContentFromContent).Methods(http.MethodGet)
	r.HandleFunc("/metatagcontent/list/tag", h.listMetaTagContentFromTag).Methods(http.MethodGet)
	r.HandleFunc("/metatagcontent/list/notebook", h.listMetaTagContentFromNotebook).Methods(http.MethodGet)
	r.HandleFunc("/metatagcontent", h.getMetaTagContent).Methods(http.MethodGet)
	r.HandleFunc("/metatagcontent", h.updateMetaTagContent).Methods(http.MethodPut)
	r.HandleFunc("/metatagcontent", h.deleteMetaTagContent).Methods(http.MethodDelete)
}

func (h *MetaTagContentHandler) createMetaTagContent(w http.ResponseWriter, r *http.Request) {
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

	contentID, err := uuid.Parse(r.Header.Get("content_id"))
	if err != nil {
		http.Error(w, "Erro ao fazer o parse do uuid: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tagID, err := uuid.Parse(r.Header.Get("tag_id"))
	if err != nil {
		http.Error(w, "Erro ao fazer o parse do uuid: "+err.Error(), http.StatusInternalServerError)
		return
	}

	request := model.MetaTagContentRequestDTO{
		UserID:     user.UserID, //userID do usuário autenticado (context)
		NotebookID: notebookID,  //notebookID injetado do Header no DTO
		ContentID:  contentID,   //contentID injetado do Header no DTO
		TagID:      tagID,       //tagID injetado do Header no DTO
	}
	resp, err := h.service.Create(ctx, request)
	if err != nil {
		http.Error(w, "Erro ao salvar meta tag content: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Falha ao codificar resposta: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *MetaTagContentHandler) listMetaTagContentFromContent(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := r.Context()

	w.Header().Set("Content-Type", "application/json")

	contentID, err := uuid.Parse(r.Header.Get("content_id"))
	if err != nil {
		http.Error(w, "Erro ao fazer o parse do uuid: "+err.Error(), http.StatusInternalServerError)
		return
	}

	request := model.ListMetaTagContentsFromContentDTO{
		ContentID: contentID,
	}
	list_meta_tags, err := h.service.ListMetaTagContentsFromContent(ctx, request)
	if err != nil {
		http.Error(w, "Erro ao procurar lista de Meta Tag Contents do content: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(list_meta_tags); err != nil {
		http.Error(w, "Falha ao codificar resposta: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *MetaTagContentHandler) listMetaTagContentFromTag(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := r.Context()

	w.Header().Set("Content-Type", "application/json")

	tagID, err := uuid.Parse(r.Header.Get("tag_id"))
	if err != nil {
		http.Error(w, "Erro ao fazer o parse do uuid: "+err.Error(), http.StatusInternalServerError)
		return
	}

	request := model.ListMetaTagContentsFromTagDTO{
		TagID: tagID,
	}
	list_meta_tags, err := h.service.ListMetaTagContentsFromTag(ctx, request)
	if err != nil {
		http.Error(w, "Erro ao procurar lista de Meta Tag Contents da tag: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(list_meta_tags); err != nil {
		http.Error(w, "Falha ao codificar resposta: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *MetaTagContentHandler) listMetaTagContentFromNotebook(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := r.Context()

	w.Header().Set("Content-Type", "application/json")

	notebookID, err := uuid.Parse(r.Header.Get("notebook_id"))
	if err != nil {
		http.Error(w, "Erro ao fazer o parse do uuid: "+err.Error(), http.StatusInternalServerError)
		return
	}

	request := model.ListMetaTagContentsFromNotebookDTO{
		NotebookID: notebookID,
	}
	list_meta_tags, err := h.service.ListMetaTagContentsFromNotebook(ctx, request)
	if err != nil {
		http.Error(w, "Erro ao procurar lista de Meta Tag Contents do notebook: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(list_meta_tags); err != nil {
		http.Error(w, "Falha ao codificar resposta: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *MetaTagContentHandler) getMetaTagContent(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := r.Context()

	w.Header().Set("Content-Type", "application/json")

	contentID, err := uuid.Parse(r.Header.Get("content_id"))
	if err != nil {
		http.Error(w, "Erro ao fazer o parse do uuid: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tagID, err := uuid.Parse(r.Header.Get("tag_id"))
	if err != nil {
		http.Error(w, "Erro ao fazer o parse do uuid: "+err.Error(), http.StatusInternalServerError)
		return
	}

	meta_tag_found, err := h.service.GetMetaTagContent(ctx, contentID, tagID)
	if err != nil {
		http.Error(w, "Erro ao procurar meta tag content: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(meta_tag_found); err != nil {
		http.Error(w, "Falha ao codificar resposta: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *MetaTagContentHandler) updateMetaTagContent(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := r.Context()

	w.Header().Set("Content-Type", "application/json")

	contentID, err := uuid.Parse(r.Header.Get("content_id"))
	if err != nil {
		http.Error(w, "Erro ao fazer o parse do uuid: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tagID, err := uuid.Parse(r.Header.Get("tag_id"))
	if err != nil {
		http.Error(w, "Erro ao fazer o parse do uuid: "+err.Error(), http.StatusInternalServerError)
		return
	}

	notebookID, err := uuid.Parse(r.Header.Get("notebook_id"))
	if err != nil {
		http.Error(w, "Erro ao fazer o parse do uuid: "+err.Error(), http.StatusInternalServerError)
		return
	}

	request := model.UpdateMetaTagContentDTO{
		NotebookID: notebookID, //notebookID injetado do Header no DTO
	}
	resp, err := h.service.Update(ctx, contentID, tagID, request)
	if err != nil {
		http.Error(w, "Erro ao atualizar meta tag content: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Falha ao codificar resposta: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *MetaTagContentHandler) deleteMetaTagContent(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := r.Context()

	w.Header().Set("Content-Type", "application/json")

	contentID, err := uuid.Parse(r.Header.Get("content_id"))
	if err != nil {
		http.Error(w, "Erro ao fazer o parse do uuid: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tagID, err := uuid.Parse(r.Header.Get("tag_id"))
	if err != nil {
		http.Error(w, "Erro ao fazer o parse do uuid: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := h.service.Delete(ctx, contentID, tagID); err != nil {
		http.Error(w, "Erro ao deletar meta tag content: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
