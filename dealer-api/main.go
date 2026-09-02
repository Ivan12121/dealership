package dealer_api

import (
	"log"
	"net/http"
	"dealer-api/handlers"
)

func main() {
	http.HandleFunc("auth/login", handlers.LoginHandler)
	http.HandleFunc("auth/register", handlers.RegisterHandler)

	log.Println("dealer-api запущен на порту 8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Ошибка запуска сервера: ", err)
	}
}