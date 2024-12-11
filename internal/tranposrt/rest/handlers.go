package rest

import "net/http"

// Регистрация нового пользователя
func Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

// Авторизация пользователя
func Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

// Добавление заказа
func PostOrders() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

// Получение списка заказов
func GetOrders() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

// Получение баланса
func GetBalance() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

// Запрос на вывод бонусов
func WithdrawRequest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

// Получение списка выводов
func GetWithdrawls() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
