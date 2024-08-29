INSERT INTO tasks (title, description, status)
VALUES ($1, $2, $3)
    RETURNING *;

SELECT * FROM tasks
WHERE id = $1;

SELECT * FROM tasks
ORDER BY created_at;

UPDATE tasks
SET title = $2, description = $3, status = $4, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
    RETURNING *;

DELETE FROM tasks
WHERE id = $1;

INSERT INTO authors (name, email)
VALUES ($1, $2)
    RETURNING *;

SELECT * FROM authors
WHERE id = $1;

SELECT * FROM authors
ORDER BY name;

INSERT INTO task_authors (task_id, author_id)
VALUES ($1, $2)
    ON CONFLICT DO NOTHING;

SELECT a.* FROM authors a
                    JOIN task_authors ta ON a.id = ta.author_id
WHERE ta.task_id = $1;
