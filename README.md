# Task Management System

This project is a simple Go REST API for task management system with a background process that periodically runs every 5 minutes to check for tasks marked as completed for over 7 days and delete them from the database.

## Features
- Create, Retrieve, Update, and Delete tasks
- Background task using Goroutines running every 5 minutes to delete completed tasks older than 7 days

## Tech Stack 
- Go
- PostgreSQL
- Docker

## Setup Instructions

1. Clone the repository.
   ```
   git clone https://github.com/Peeranut-Kit/go_backend_test.git
   cd go_backend_test
   ```

2. Use a .env file provided in the repository. Please ensure that the config in .env file is not the same as any service running on your system (port number, database name).

3. Download and initialize PostgreSQL container in docker.
   If you do not have docker installed in your computer, feel free to check their website. https://www.docker.com/

   ```
   docker run --name postgresTask -e POSTGRES_PASSWORD=password -p 5432:5432 -d postgres
   ```

4. Set up PostgreSQL database and create table.
   1. Open a shell inside a running Docker container.
   ```
   docker exec -it postgresTask psql -U postgres
   ```
   2. Create a table using SQL command in postgres_init.sql file to create tasks table in the database.
  
5. Run Go API service.
   ```
   go run main.go
   ```

## Endpoints
1. GET /tasks
2. GET /tasks/{id}
3. POST /tasks
4. PUT /tasks/{id}
5. DELETE /tasks/{id}

## Background Task (Cronjob)
The background routine is implemented in service folder. The results are logged into background_task.log file in the same directory.
