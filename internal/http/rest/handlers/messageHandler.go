package handlers

import (
	"fmt"
	"log/slog"
	"myHttpServer/internal/domain"
	"myHttpServer/internal/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type MessageHandler struct {
	service domain.MessageService
	log     *slog.Logger
}

func NewMessageHandler(service domain.MessageService, log *slog.Logger) *MessageHandler {
	return &MessageHandler{service: service, log: log}
}

func (h *MessageHandler) Message() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		h.log.Debug("handle Message() called")
		id, err := h.handleID(w, r)
		if err != nil {
			http.Error(w, "Id is not defined", http.StatusBadRequest)
			return
		}

		m, err := h.service.Message(id)
		if err != nil {
			http.Error(
				w,
				fmt.Sprintf("Can't find message with requested id: %v", id),
				http.StatusNotFound,
			)
			return
		}

		if err := json.Format(w, m); err != nil {
			h.log.Error("Error while formating message to json")
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
func (h *MessageHandler) Messages() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		h.log.Debug("handle Messages() called")

		var mm *domain.MessageSlice
		mm, err := h.service.Messages()
		if err != nil {
			h.log.Error("Error requesting messages: $1", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		if err := json.Format(w, *mm); err != nil {
			h.log.Error("Error formating messages to json")
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}

func (h *MessageHandler) CreateMessage() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		h.log.Debug("handler CreateMessage() callled")

		m := domain.Message{CreatedAt: time.Now()}
		if err := json.Parse(r.Body, &m); err != nil {
			h.log.Debug("CreateMessage(): parse error" + err.Error())
			http.Error(w, "Can't parse message", http.StatusBadRequest)
			return
		}

		if err := h.service.CreateMessage(&m); err != nil {
			h.log.Debug("CreateMessage(): create error" + err.Error())
			http.Error(w, "Can't create message", http.StatusInternalServerError)
			return
		}
	}
}
func (h *MessageHandler) UpdateMessage() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		h.log.Debug("handler UpdateMessage() called")

		id, err := h.handleID(w, r)
		if err != nil {
			http.Error(w, "Id is not defined", http.StatusBadRequest)
			return
		}

		m := domain.Message{}
		if err := json.Parse(r.Body, &m); err != nil {
			h.log.Debug("CreateMessage(): parse error" + err.Error())
			http.Error(w, "Can't parse message", http.StatusBadRequest)
			return
		}
		m.ID = id

		if err := h.service.UpdateMessage(id, &m); err != nil {
			h.log.Debug("UpdateMessage(): update error, id:$1", id)
			http.Error(w, "Can't update message", http.StatusInternalServerError)
			return
		}
	}
}

func (h *MessageHandler) DeleteMessage() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		h.log.Debug("handler DeleteMessage() called")

		id, err := h.handleID(w, r)
		if err != nil {
			http.Error(w, "Id is not defined", http.StatusBadRequest)
			return
		}

		if err := h.service.DeleteMessage(id); err != nil {
			h.log.Debug("DeleteMessage(): delete error, id:$1", id)
			http.Error(w, "Can't delete message", http.StatusInternalServerError)
		}
	}
}

func (h *MessageHandler) handleID(w http.ResponseWriter, r *http.Request) (int, error) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return 0, err
	}
	return id, nil
}
