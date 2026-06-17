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

type TagHandler struct {
	service *service.TagService
}

func NewTagHandler(srv *service.TagService) *TagHandler {
	return &TagHandler{
		service: srv,
	}
}

func (h *TagHandler) MountTagHandler(r *mux.Router) {
	r.HandleFunc("/tag", h.CreateTag).Methods(http.MethodPost)
	r.HandleFunc("/tag/list", h.ListTagFromUser).Methods(http.MethodGet)
	r.HandleFunc("/tag", h.GetTagByID).Methods(http.MethodGet)
	r.HandleFunc("/tag", h.UpdateTag).Methods(http.MethodPut)
	r.HandleFunc("/tag", h.DeleteTag).Methods(http.MethodDelete)
}

func (h *TagHandler) CreateTag(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := context.TODO()

	w.Header().Set("Content-Type", "application/json")

	userID, err := uuid.Parse(r.Header.Get("user_id"))
	if err != nil {
		http.Error(w, "Erro ao fazer o parse do uuid: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var request model.TagRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "JSON inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	request.UserID = userID //userID injetado do Header no DTO
	resp, err := h.service.Create(ctx, request)
	if err != nil {
		http.Error(w, "Erro ao salvar tag: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Falha ao codificar resposta: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *TagHandler) ListTagFromUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := context.TODO()

	w.Header().Set("Content-Type", "application/json")

	userID, err := uuid.Parse(r.Header.Get("user_id"))
	if err != nil {
		http.Error(w, "Erro ao fazer o parse do uuid: "+err.Error(), http.StatusInternalServerError)
		return
	}

	request := model.ListTagsFromUserDTO{
		UserID: userID,
	}
	list_tags, err := h.service.ListTagsFromUser(ctx, request)
	if err != nil {
		http.Error(w, "Erro ao procuar lista de Tags do user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(list_tags); err != nil {
		http.Error(w, "Falha ao codificar resposta: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *TagHandler) GetTagByID(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := context.TODO()

	w.Header().Set("Content-Type", "application/json")

	tagID, err := uuid.Parse(r.Header.Get("tag_id"))
	if err != nil {
		http.Error(w, "Erro ao fazer o parse do uuid: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tag_found, err := h.service.GetTagByID(ctx, tagID)
	if err != nil {
		http.Error(w, "Erro ao procurar tag: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(tag_found); err != nil {
		http.Error(w, "Falha ao codificar resposta: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *TagHandler) UpdateTag(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := context.TODO()

	w.Header().Set("Content-Type", "application/json")

	tagID, err := uuid.Parse(r.Header.Get("tag_id"))
	if err != nil {
		http.Error(w, "Erro ao fazer o parse do uuid: "+err.Error(), http.StatusInternalServerError)
		return
	}

	userID, err := uuid.Parse(r.Header.Get("user_id"))
	if err != nil {
		http.Error(w, "Erro ao fazer o parse do uuid: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var request model.TagRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "JSON inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	request.UserID = userID
	resp, err := h.service.Update(ctx, tagID, request)
	if err != nil {
		http.Error(w, "Erro ao atualizar tag: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Falha ao codificar resposta: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *TagHandler) DeleteTag(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := context.TODO()

	w.Header().Set("Content-Type", "application/json")

	tagID, err := uuid.Parse(r.Header.Get("tag_id"))
	if err != nil {
		http.Error(w, "Erro ao fazer o parse do uuid: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := h.service.Delete(ctx, tagID); err != nil {
		http.Error(w, "Erro ao deletar tag: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
