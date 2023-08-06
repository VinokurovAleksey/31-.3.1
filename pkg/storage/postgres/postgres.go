package postgres

import (
	"GoNews/pkg/storage"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// База данных Postgres
type Database struct {
	db *sql.DB
}

// Создает новый экземпляр структуры
func NewDatabase(connection string) (*Database, error) {
	db, err := sql.Open("postgres", connection)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping the database: %v", err)
	}

	return &Database{db: db}, nil
}

// Закрывает соединение с базой данных
func (d *Database) Close() error {
	return d.db.Close()
}

// Извлекает все сообщения
func (d *Database) Posts() ([]storage.Post, error) {
	rows, err := d.db.Query("SELECT id, title, content FROM posts")
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve posts: %v", err)
	}
	defer rows.Close()

	var posts []storage.Post
	for rows.Next() {
		var post storage.Post
		err := rows.Scan(&post.ID, &post.Title, &post.Content)
		if err != nil {
			return nil, fmt.Errorf("failed to scan post: %v", err)
		}
		posts = append(posts, post)
	}

	return posts, nil
}

// Добавляет новый пост
func (d *Database) AddPost(post storage.Post) error {
	_, err := d.db.Exec("INSERT INTO posts (id, title, content) VALUES ($1, $2, $3)", post.ID, post.Title, post.Content)
	if err != nil {
		return fmt.Errorf("failed to add post: %v", err)
	}
	return nil
}

// Обновляет существующий пост
func (d *Database) UpdatePost(post storage.Post) error {
	_, err := d.db.Exec("UPDATE posts SET title = $1, content = $2 WHERE id = $3", post.Title, post.Content, post.ID)
	if err != nil {
		return fmt.Errorf("failed to update post: %v", err)
	}
	return nil
}

// Удаляет запись
func (d *Database) DeletePost(post storage.Post) error {
	_, err := d.db.Exec("DELETE FROM posts WHERE id = $1", post.ID)
	if err != nil {
		return fmt.Errorf("failed to delete post: %v", err)
	}
	return nil
}
