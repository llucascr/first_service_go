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

type MetaTagContentHandler struct {
	service *service.MetaTagContentService
}

func NewMetaTagContentHandler(srv *service.MetaTagContentService) *MetaTagContentHandler {
	return &MetaTagContentHandler{
		service: srv,
	}
}

func (h *MetaTagContentHandler) MountMetaTagContentHandler(r *mux.Router) {
	r.HandleFunc("/metatagcontent", h.CreateMetaTagContent).Methods(http.MethodPost)
	r.HandleFunc("/metatagcontent/list/content", h.ListMetaTagContentFromContent).Methods(http.MethodGet)
	r.HandleFunc("/metatagcontent/list/tag", h.ListMetaTagContentFromTag).Methods(http.MethodGet)
	r.HandleFunc("/metatagcontent/list/notebook", h.ListMetaTagContentFromNotebook).Methods(http.MethodGet)
	r.HandleFunc("/metatagcontent", h.GetMetaTagContent).Methods(http.MethodGet)
	r.HandleFunc("/metatagcontent", h.UpdateMetaTagContent).Methods(http.MethodPut)
	r.HandleFunc("/metatagcontent", h.DeleteMetaTagContent).Methods(http.MethodDelete)
}

func (h *MetaTagContentHandler) CreateMetaTagContent(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := context.TODO()

	w.Header().Set("Content-Type", "application/json")

	userID, err := uuid.Parse(r.Header.Get("user_id"))
	if err != nil {
		http.Error(w, "Erro ao fazer o parse do uuid: "+err.Error(), http.StatusInternalServerError)
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
		UserID:     userID,     //userID injetado do Header no DTO
		NotebookID: notebookID, //notebookID injetado do Header no DTO
		ContentID:  contentID,  //contentID injetado do Header no DTO
		TagID:      tagID,      //tagID injetado do Header no DTO
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

func (h *MetaTagContentHandler) ListMetaTagContentFromContent(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := context.TODO()

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

func (h *MetaTagContentHandler) ListMetaTagContentFromTag(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := context.TODO()

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

func (h *MetaTagContentHandler) ListMetaTagContentFromNotebook(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := context.TODO()

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

func (h *MetaTagContentHandler) GetMetaTagContent(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := context.TODO()

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

func (h *MetaTagContentHandler) UpdateMetaTagContent(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := context.TODO()

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

func (h *MetaTagContentHandler) DeleteMetaTagContent(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := context.TODO()

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