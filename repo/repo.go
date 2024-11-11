package repo

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/Peeranut-Kit/go_backend_test/utils"
)

var ErrTaskNotFound = errors.New("task not found")

type PostgresDB struct {
	db *sql.DB
}

func NewPostgresDB(db *sql.DB) *PostgresDB {
	return &PostgresDB{db: db}
}

/*type TaskRepository interface {
	GetTasks() ([]utils.Task, error)
	CreateTask(task utils.Task) (utils.Task, error)
	GetTaskById(id int) (utils.Task, error)
	UpdateTask(id int, task utils.Task) (utils.Task, error)
	DeleteTask(id int) error

	GetOldFinishedTasks() ([]utils.Task, error)
}*/

func (postgres *PostgresDB) GetTasks() ([]utils.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var tasks []utils.Task
	const query = `SELECT * FROM tasks`
	result, err := postgres.db.QueryContext(ctx, query)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer result.Close()
	for result.Next() {
		var task utils.Task
		result.Scan(
			&task.Id,
			&task.Title,
			&task.Description,
			&task.Completed,
			&task.CreatedAt,
		)

		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (postgres *PostgresDB) CreateTask(task utils.Task) (utils.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var createdTask utils.Task
	const query = `INSERT INTO tasks (title, description, completed) VALUES ($1, $2, $3) RETURNING id, title, description, completed, created_at`
	row := postgres.db.QueryRowContext(ctx, query, task.Title, task.Description, false)
	err := row.Scan(
		&createdTask.Id,
		&createdTask.Title,
		&createdTask.Description,
		&createdTask.Completed,
		&createdTask.CreatedAt,
	)

	if err != nil {
		log.Println(err.Error())
		return utils.Task{}, err
	}

	return createdTask, nil
}

func (postgres *PostgresDB) GetTaskById(id int) (utils.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var task utils.Task
	const query = `SELECT * FROM tasks WHERE id = $1`
	row := postgres.db.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&task.Id,
		&task.Title,
		&task.Description,
		&task.Completed,
		&task.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return utils.Task{}, ErrTaskNotFound
	} else if err != nil {
		log.Println(err.Error())
		return utils.Task{}, err
	}

	return task, nil
}

func (postgres *PostgresDB) UpdateTask(id int, task utils.Task) (utils.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var updatedTask utils.Task
	const query = `UPDATE tasks SET title = $1, description = $2, completed = $3 WHERE id = $4 RETURNING id, title, description, completed, created_at`
	row := postgres.db.QueryRowContext(ctx, query, task.Title, task.Description, task.Completed, id)
	err := row.Scan(
		&updatedTask.Id,
		&updatedTask.Title,
		&updatedTask.Description,
		&updatedTask.Completed,
		&updatedTask.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return utils.Task{}, ErrTaskNotFound
	} else if err != nil {
		log.Println(err.Error())
		return utils.Task{}, err
	}

	return updatedTask, nil
}

func (postgres *PostgresDB) DeleteTask(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	const query = `DELETE FROM tasks WHERE id = $1`
	_, err := postgres.db.ExecContext(ctx, query, id)

	if err == sql.ErrNoRows {
		return ErrTaskNotFound
	} else if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func (postgres *PostgresDB) GetOldFinishedTasks() ([]utils.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var tasks []utils.Task
	const query = `SELECT * FROM tasks WHERE completed = true AND created_at < NOW() - INTERVAL '7 days';`
	result, err := postgres.db.QueryContext(ctx, query)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer result.Close()
	for result.Next() {
		var task utils.Task
		result.Scan(
			&task.Id,
			&task.Title,
			&task.Description,
			&task.Completed,
			&task.CreatedAt,
		)

		tasks = append(tasks, task)
	}
	return tasks, nil
}
