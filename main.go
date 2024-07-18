package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	ParcelStatusSent      = "sent"
	ParcelStatusDelivered = "delivered"
)

// ParcelService предоставляет методы для работы с посылками
type ParcelService struct {
	store ParcelStore
}

// NewParcelService создает новый экземпляр ParcelService с переданным хранилищем посылок
func NewParcelService(store ParcelStore) ParcelService {
	return ParcelService{store: store}
}

// Register регистрирует новую посылку и возвращает её
func (s ParcelService) Register(client int, address string) (Parcel, error) {
	parcel := Parcel{
		Client:    client,
		Status:    ParcelStatusRegistered,
		Address:   address,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}

	id, err := s.store.Add(parcel)
	if err != nil {
		return parcel, err
	}

	parcel.Number = id

	fmt.Printf("Новая посылка № %d на адрес %s от клиента с идентификатором %d зарегистрирована %s\n",
		parcel.Number, parcel.Address, parcel.Client, parcel.CreatedAt)

	return parcel, nil
}

// PrintClientParcels выводит все посылки клиента на экран
func (s ParcelService) PrintClientParcels(client int) error {
	parcels, err := s.store.GetByClient(client)
	if err != nil {
		return err
	}

	fmt.Printf("Посылки клиента %d:\n", client)
	for _, parcel := range parcels {
		fmt.Printf("Посылка № %d на адрес %s от клиента с идентификатором %d зарегистрирована %s, статус %s\n",
			parcel.Number, parcel.Address, parcel.Client, parcel.CreatedAt, parcel.Status)
	}
	fmt.Println()

	return nil
}

// NextStatus переводит посылку на следующий статус
func (s ParcelService) NextStatus(number int) error {
	parcel, err := s.store.Get(number)
	if err != nil {
		return err
	}

	var nextStatus string
	switch parcel.Status {
	case ParcelStatusRegistered:
		nextStatus = ParcelStatusSent
	case ParcelStatusSent:
		nextStatus = ParcelStatusDelivered
	case ParcelStatusDelivered:
		return nil
	}

	fmt.Printf("У посылки № %d новый статус: %s\n", number, nextStatus)

	return s.store.SetStatus(number, nextStatus)
}

// ChangeAddress изменяет адрес посылки
func (s ParcelService) ChangeAddress(number int, address string) error {
	return s.store.SetAddress(number, address)
}

// Delete удаляет посылку
func (s ParcelService) Delete(number int) error {
	return s.store.Delete(number)
}

// createTableIfNotExists создает таблицу parcel, если она не существует
func createTableIfNotExists(db *sql.DB) error {
	createTableSQL := `CREATE TABLE IF NOT EXISTS parcel (
        number INTEGER PRIMARY KEY AUTOINCREMENT,
        client INTEGER NOT NULL,
        status TEXT NOT NULL,
        address TEXT NOT NULL,
        created_at TEXT NOT NULL
    );`

	_, err := db.Exec(createTableSQL)
	return err
}

func main() {
	// Настройка подключения к базе данных
	db, err := sql.Open("sqlite3", "tracker.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	// Создание таблицы, если она не существует
	err = createTableIfNotExists(db)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Создание объекта ParcelStore
	store := NewParcelStore(db)
	service := NewParcelService(store)

	// Регистрация посылки
	client := 1
	address := "Псков, д. Пушкина, ул. Колотушкина, д. 5"
	p, err := service.Register(client, address)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Изменение адреса
	newAddress := "Саратов, д. Верхние Зори, ул. Козлова, д. 25"
	err = service.ChangeAddress(p.Number, newAddress)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Изменение статуса
	err = service.NextStatus(p.Number)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Вывод посылок клиента
	err = service.PrintClientParcels(client)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Попытка удаления отправленной посылки
	err = service.Delete(p.Number)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Вывод посылок клиента (предыдущая посылка не должна удалиться, так как её статус НЕ «зарегистрирована»)
	err = service.PrintClientParcels(client)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Регистрация новой посылки
	p, err = service.Register(client, address)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Удаление новой посылки
	err = service.Delete(p.Number)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Вывод посылок клиента (здесь не должно быть последней посылки, так как она должна была успешно удалиться)
	err = service.PrintClientParcels(client)
	if err != nil {
		fmt.Println(err)
		return
	}
}
