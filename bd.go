package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type Storage struct {
	db *sql.DB
}

func InitDb() *Storage {
	db, err := sql.Open("mysql", "root:serpent@/tasks")
	if err != nil {
		fmt.Println("Ошибка:", err)
	} else {
		fmt.Println("Соединение с БД")
	}
	storage := Storage{
		db: db,
	}
	return &storage
}

// Создание новой задачи
func (s *Storage) Create(u Url) error {
	res, err := s.db.Prepare("INSERT INTO `url`(id, url, short, ttl) VALUES (NULL, ?, ?, ?);")
	if err != nil {
		fmt.Println(err)
		return err
	}
	_, err = res.Exec(u.Link, u.Short, u.Ttl)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return err
}

// Получение ссылки по короткой ссылке
func (s *Storage) GetUrl(token string) (string, error) {
	var url string
	err := s.db.QueryRow("SELECT url FROM url WHERE short= ?", token).Scan(&url)

	if err != nil {
		fmt.Println(err)
		return url, err
	}

	return url, nil
}
