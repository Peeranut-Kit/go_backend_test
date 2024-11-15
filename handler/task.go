package handler

import (
	"log"
	"strconv"

	"github.com/Peeranut-Kit/go_backend_test/repo"
	"github.com/Peeranut-Kit/go_backend_test/utils"
	"github.com/gofiber/fiber/v2"
)

type TaskHandlerInterface interface {
	GetTasksHandler(c *fiber.Ctx) error
	GetTaskHandler(c *fiber.Ctx) error
	PostTaskHandler(c *fiber.Ctx) error
	PutTaskHandler(c *fiber.Ctx) error
	DeleteTaskHandler(c *fiber.Ctx) error
}

// Primary adapter
type HttpTaskHandler struct {
	TaskRepo repo.TaskRepositoryInterface
}

// Initiate primary adapter
func NewHttpTaskHandler(repo repo.TaskRepositoryInterface) *HttpTaskHandler {
	return &HttpTaskHandler{TaskRepo: repo}
}

func (h *HttpTaskHandler) GetTasksHandler(c *fiber.Ctx) error {
	tasks, err := h.TaskRepo.GetTasks()
	if err != nil {
		log.Println("Error getting tasks:", err)
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(tasks)
}

func (h *HttpTaskHandler) GetTaskHandler(c *fiber.Ctx) error {
	taskId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	task, err := h.TaskRepo.GetTaskById(taskId)
	if err != nil {
		if err == utils.ErrNotFound {
			return c.Status(fiber.StatusNotFound).SendString(err.Error())
		} else {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
	}

	return c.JSON(task)
}

func (h *HttpTaskHandler) PostTaskHandler(c *fiber.Ctx) error {
	task := &utils.Task{}
	// or task := new(utils.Task) because BodyParser() expects a pointer to a struct, not the struct itself.
	if err := c.BodyParser(task); err != nil {
		log.Println("Error decoding request body:", err)
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	// Core Logic
	userIdString := c.Locals("user_id").(string)
	userIDInt, err := strconv.Atoi(userIdString)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString(err.Error())
	}
	task.UserID = userIDInt

	createdTask, err := h.TaskRepo.CreateTask(task)
	if err != nil {
		log.Println("Error creating task:", err)
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(fiber.Map{
		"message":     "Create Task Successful",
		"createdTask": createdTask,
	})
}

func (h *HttpTaskHandler) PutTaskHandler(c *fiber.Ctx) error {
	taskId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	task := new(utils.Task)
	if err := c.BodyParser(task); err != nil {
		log.Println("Error decoding request body:", err)
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	updatedTask, err := h.TaskRepo.UpdateTask(taskId, task)
	if err != nil {
		if err == utils.ErrNotFound {
			return c.Status(fiber.StatusNotFound).SendString(err.Error())
		} else {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
	}

	return c.JSON(fiber.Map{
		"message":     "Update Task Successful",
		"updatedTask": updatedTask,
	})
}

func (h *HttpTaskHandler) DeleteTaskHandler(c *fiber.Ctx) error {
	taskId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	err = h.TaskRepo.DeleteTask(taskId)
	if err != nil {
		if err == utils.ErrNotFound {
			return c.Status(fiber.StatusNotFound).SendString(err.Error())
		} else {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
	}

	return c.SendStatus(fiber.StatusNoContent)
}
