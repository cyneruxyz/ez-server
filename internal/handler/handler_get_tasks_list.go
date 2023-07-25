package handler

import (
	"encoding/json"
	"ex-server/internal/action"
	"net/http"
)

func (h *Handler) GetTasksList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	getTasksListAct := action.NewGetTasksList(h.TaskRepo)
	tasks, err := getTasksListAct.Do(h.db)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tasks)
}