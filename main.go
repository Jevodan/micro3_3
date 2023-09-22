package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"github.com/streadway/amqp"
)

/*
2ой сервис:
Вычитывает из очереди кролика сообщения. Создаёт короткую ссылку и кладет в базу.
База вида: урл, короткая ссылка, ттл.
**** При запуске микросервиса ожидает сообщения в кролике и считывает их при появлении, затем добавляет в БД

Соответственно если переходишь по этой ссылке и ттл не истёк, он тебя редиректит.
**** Ссылка в виде: localhost:8001/to/?key=RySCu
**** key - короткая ссылка
*/

const (
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

var (
	db *Storage
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Не удалось подключиться к Rabbit")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Не удалось открыть канал")
	//defer ch.Close()
	queueName := "first_queue"

	//создайте потребителя (consumer). В этом примере мы будем выводить полученные сообщения в консоль:
	msgs, err := ch.Consume(
		queueName, // Имя очереди
		"",        // Consumer
		true,      // AutoAck - автоматическое подтверждение получения сообщения
		false,     // Exclusive
		false,     // NoLocal
		false,     // NoWait
		nil,       // Args
	)
	failOnError(err, "Не удалось зарегистрировать потребителя")

	go func() {
		db = InitDb()
		for d := range msgs {
			var url Url
			err := json.Unmarshal(d.Body, &url)
			if err != nil {
				fmt.Println("Ошибка декодирования JSON:", err)
				continue
			}

			link := Url{Link: url.Link, Short: shorting(), Ttl: 4200}
			fmt.Printf("Получено сообщение: %s\n", d.Body)
			db.Create(link)
		}
	}()

	fmt.Println("Ожидание сообщений. Для завершения нажмите CTRL+C")
	http.HandleFunc("/to/", redirectHandle)
	log.Fatal(http.ListenAndServe(":8001", nil))
	select {}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s : %s", msg, err)
	}
}

func redirectHandle(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	url, _ := db.GetUrl(key)
	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusFound)
	fmt.Fprintf(w, "<script>location='%s';</script>", url)
	//http.Redirect(w, r, url, http.StatusFound)
}

func shorting() string {
	b := make([]byte, 5)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
