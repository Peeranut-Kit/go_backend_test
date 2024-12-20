package service

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/Peeranut-Kit/go_backend_test/repo"
)

func BackgroundTask(r repo.TaskRepositoryInterface) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	// Run a loop to handle each tick
	for range ticker.C {
		cleanupOldTasks(r)
	}
}

func cleanupOldTasks(r repo.TaskRepositoryInterface) {
	// delete completed tasks older than 7 days
	tasks, err := r.GetOldFinishedTasks()
	if err != nil {
		log.Println("Error fetching old finished tasks:", err)
		return
	}

	// open logging file
	logFile, err := os.OpenFile("service/background_task.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal("Error opening log file:", err)
	}
	defer logFile.Close()

	for _, task := range tasks {
		log.Printf("Deleting old finished task ID: %d\n", task.ID)

		// Log the task to the file
		taskByte, err := json.Marshal(task)
		if err != nil {
			log.Println("Error marshalling task:", err)
			return
		}
		logFile.Write(taskByte)

		// delete the task from database
		err = r.DeleteTask(int(task.ID))
		if err != nil {
			log.Println("Error deleting task:", err)
			return
		}
	}
}
