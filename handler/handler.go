package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/Peeranut-Kit/go_backend_test/repo"
	"github.com/Peeranut-Kit/go_backend_test/utils"
)

type TaskHandler struct {
	DB *sql.DB
	TaskRepo *repo.PostgresDB
}

func (h TaskHandler) TasksHandler(w http.ResponseWriter, r *http.Request) {
	// Route for /tasks
	switch r.Method {
	case http.MethodGet:
		// GET /tasks
		taskList, err := h.TaskRepo.GetTasks()
		if err != nil {
			log.Println("Error getting tasks:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		taskListJSON, err := json.Marshal(taskList)
		if err != nil {
			log.Println("Error marshalling task list:", err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(taskListJSON)
		if err != nil {
			log.Println("Error writing response:", err)
		}
		return
	case http.MethodPost:
		// POST /tasks
		var task utils.Task
		err := json.NewDecoder(r.Body).Decode(&task)
		if err != nil {
			log.Println("Error decoding request body:", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		createdTask, err := h.TaskRepo.CreateTask(task)
		if err != nil {
			log.Println("Error creating task:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(createdTask)
		if err != nil {
			log.Println("Error encoding created task:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func (h TaskHandler) TaskHandlerByID(w http.ResponseWriter, r *http.Request) {
	// Route for /tasks/{id}
	urlPathSegments := strings.Split(r.URL.Path, "/tasks/")
	if len(urlPathSegments) > 2 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	taskId, err := strconv.Atoi(urlPathSegments[len(urlPathSegments)-1])
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodGet:
		// GET /tasks/{id}
		task, err := h.TaskRepo.GetTaskById(taskId)
		if err != nil {
			if err == repo.ErrTaskNotFound {
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(map[string]string{
					"message": "Task not found",
				})
			} else {
				log.Println("Error fetching task:", err)
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}

		taskJSON, err := json.Marshal(task)
		if err != nil {
			log.Println("Error marshalling task:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(taskJSON)
		if err != nil {
			log.Println("Error writing response:", err)
		}
		return
	case http.MethodPut:
		// PUT /tasks/{id}
		var task utils.Task
		err := json.NewDecoder(r.Body).Decode(&task)
		if err != nil {
			log.Println("Error decoding request body:", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		updatedTask, err := h.TaskRepo.UpdateTask(taskId, task)
		if err != nil {
			if err == repo.ErrTaskNotFound {
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(map[string]string{
					"message": "Task not found",
				})
			} else {
				log.Println("Error updating task:", err)
				w.WriteHeader(http.StatusInternalServerError) // 500 Internal Server Error for other issues
				json.NewEncoder(w).Encode(map[string]string{
					"message": "Error deleting task",
				})
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(updatedTask)
		if err != nil {
			log.Println("Error encoding response:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	case http.MethodDelete:
		// DELETE /tasks/{id}
		err := h.TaskRepo.DeleteTask(taskId)
		if err != nil {
			if err == repo.ErrTaskNotFound {
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(map[string]string{
					"message": "Task not found",
				})
			} else {
				log.Println("Error deleting task:", err)
				w.WriteHeader(http.StatusInternalServerError) // 500 Internal Server Error for other issues
				json.NewEncoder(w).Encode(map[string]string{
					"message": "Error deleting task",
				})
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(map[string]string{
			"message": "Task deleted successfully",
		})
		if err != nil {
			log.Println("Error encoding response:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
