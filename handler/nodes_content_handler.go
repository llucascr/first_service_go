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

type NodesContentHandler struct {
	service *service.NodesContentService
}

func NewNodesContentHandler(srv *service.NodesContentService) *NodesContentHandler {
	return &NodesContentHandler{
		service: srv,
	}
}

func (h *NodesContentHandler) MountNodesContentHandler(r *mux.Router) {
	r.HandleFunc("/nodescontent", h.CreateNodesContent).Methods(http.MethodPost)
	r.HandleFunc("/nodescontent/list/notebook", h.ListNodesContentFromNotebook).Methods(http.MethodGet)
	r.HandleFunc("/nodescontent/list/content", h.ListNodesContentFromContent).Methods(http.MethodGet)
	r.HandleFunc("/nodescontent", h.GetNodesContentByID).Methods(http.MethodGet)
	r.HandleFunc("/nodescontent", h.UpdateNodesContent).Methods(http.MethodPut)
	r.HandleFunc("/nodescontent", h.DeleteNodesContent).Methods(http.MethodDelete)
}

func (h *NodesContentHandler) CreateNodesContent(w http.ResponseWriter, r *http.Request) {
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

	request := model.NodesContentRequestDTO{
		UserID:     userID,     //userID injetado do Header no DTO
		NotebookID: notebookID, //notebookID injetado do Header no DTO
		ContentID:  contentID,  //contentID injetado do Header no DTO
	}
	resp, err := h.service.Create(ctx, request)
	if err != nil {
		http.Error(w, "Erro ao salvar nodes content: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Falha ao codificar resposta: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *NodesContentHandler) ListNodesContentFromNotebook(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := context.TODO()

	w.Header().Set("Content-Type", "application/json")

	notebookID, err := uuid.Parse(r.Header.Get("notebook_id"))
	if err != nil {
		http.Error(w, "Erro ao fazer o parse do uuid: "+err.Error(), http.StatusInternalServerError)
		return
	}

	request := model.ListNodesContentsFromNotebookDTO{
		NotebookID: notebookID,
	}
	list_nodes, err := h.service.ListNodesContentsFromNotebook(ctx, request)
	if err != nil {
		http.Error(w, "Erro ao procurar lista de Nodes Contents do notebook: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(list_nodes); err != nil {
		http.Error(w, "Falha ao codificar resposta: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *NodesContentHandler) ListNodesContentFromContent(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := context.TODO()

	w.Header().Set("Content-Type", "application/json")

	contentID, err := uuid.Parse(r.Header.Get("content_id"))
	if err != nil {
		http.Error(w, "Erro ao fazer o parse do uuid: "+err.Error(), http.StatusInternalServerError)
		return
	}

	request := model.ListNodesContentsFromContentDTO{
		ContentID: contentID,
	}
	list_nodes, err := h.service.ListNodesContentsFromContent(ctx, request)
	if err != nil {
		http.Error(w, "Erro ao procurar lista de Nodes Contents do content: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(list_nodes); err != nil {
		http.Error(w, "Falha ao codificar resposta: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *NodesContentHandler) GetNodesContentByID(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := context.TODO()

	w.Header().Set("Content-Type", "application/json")

	nodeID, err := uuid.Parse(r.Header.Get("node_id"))
	if err != nil {
		http.Error(w, "Erro ao fazer o parse do uuid: "+err.Error(), http.StatusInternalServerError)
		return
	}

	node_found, err := h.service.GetNodesContentByID(ctx, nodeID)
	if err != nil {
		http.Error(w, "Erro ao procurar nodes content: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(node_found); err != nil {
		http.Error(w, "Falha ao codificar resposta: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *NodesContentHandler) UpdateNodesContent(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := context.TODO()

	w.Header().Set("Content-Type", "application/json")

	nodeID, err := uuid.Parse(r.Header.Get("node_id"))
	if err != nil {
		http.Error(w, "Erro ao fazer o parse do uuid: "+err.Error(), http.StatusInternalServerError)
		return
	}

	contentID, err := uuid.Parse(r.Header.Get("content_id"))
	if err != nil {
		http.Error(w, "Erro ao fazer o parse do uuid: "+err.Error(), http.StatusInternalServerError)
		return
	}

	request := model.UpdateNodesContentDTO{
		ContentID: contentID, //contentID injetado do Header no DTO
	}
	resp, err := h.service.Update(ctx, nodeID, request)
	if err != nil {
		http.Error(w, "Erro ao atualizar nodes content: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Falha ao codificar resposta: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *NodesContentHandler) DeleteNodesContent(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := context.TODO()

	w.Header().Set("Content-Type", "application/json")

	nodeID, err := uuid.Parse(r.Header.Get("node_id"))
	if err != nil {
		http.Error(w, "Erro ao fazer o parse do uuid: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := h.service.Delete(ctx, nodeID); err != nil {
		http.Error(w, "Erro ao deletar nodes content: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
