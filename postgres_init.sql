CREATE TABLE IF NOT EXISTS tasks (
	id INTEGER GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
	title VARCHAR(255) NOT NULL,
	description TEXT NOT NULL,
	completed BOOLEAN DEFAULT FALSE,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);