Main .go

package main

import (
    "database/sql"
    "fmt"
    "time"

    _ "modernc.org/sqlite"
)

const (
    ParcelStatusRegistered = "registered"
    ParcelStatusSent       = "sent"
    ParcelStatusDelivered  = "delivered"
)

type Parcel struct {
    Number    int
    Client    int
    Status    string
    Address   string
    CreatedAt string
}

type ParcelStore interface {
    Add(Parcel) (int, error)
    Get(int) (Parcel, error)
    GetByClient(int) ([]Parcel, error)
    SetStatus(int, string) error
    SetAddress(int, string) error
    Delete(int) error
}

type SQLiteStore struct {
    db *sql.DB
}

func NewParcelStore(db *sql.DB) *SQLiteStore {
    return &SQLiteStore{db: db}
}

func (s *SQLiteStore) Add(p Parcel) (int, error) {
    query := `INSERT INTO parcels (client, status, address, created_at) VALUES (?, ?, ?, ?)`
    result, err := s.db.Exec(query, p.Client, p.Status, p.Address, p.CreatedAt)
    if err != nil {
        return 0, err
    }
    id, err := result.LastInsertId()
    if err != nil {
        return 0, err
    }
    return int(id), nil
}

func (s *SQLiteStore) Get(number int) (Parcel, error) {
    query := `SELECT number, client, status, address, created_at FROM parcels WHERE number = ?`
    row := s.db.QueryRow(query, number)
    var p Parcel
    err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
    if err != nil {
        return Parcel{}, err
    }
    return p, nil
}

func (s *SQLiteStore) GetByClient(client int) ([]Parcel, error) {
    query := `SELECT number, client, status, address, created_at FROM parcels WHERE client = ?`
    rows, err := s.db.Query(query, client)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var parcels []Parcel
    for rows.Next() {
        var p Parcel
        err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
        if err != nil {
            return nil, err
        }
        parcels = append(parcels, p)
    }
    return parcels, nil
}

func (s *SQLiteStore) SetStatus(number int, status string) error {
    query := `UPDATE parcels SET status = ? WHERE number = ?`
    _, err := s.db.Exec(query, status, number)
    return err
}

func (s *SQLiteStore) SetAddress(number int, address string) error {
    query := `UPDATE parcels SET address = ? WHERE number = ?`
    _, err := s.db.Exec(query, address, number)
    return err
}

func (s *SQLiteStore) Delete(number int) error {
    query := `DELETE FROM parcels WHERE number = ?`
    _, err := s.db.Exec(query, number)
    return err
}

type ParcelService struct {
    store ParcelStore
}

func NewParcelService(store ParcelStore) ParcelService {
    return ParcelService{store: store}
}

func (s ParcelService) Register(client int, address string) (Parcel, error) {
    var p Parcel
    p = Parcel{
        Client:    client,
        Status:    ParcelStatusRegistered,
        Address:   address,
        CreatedAt: time.Now().UTC().Format(time.RFC3339),
    }

    id, err := s.store.Add(p)
    if err != nil {
        return p, err
    }

    p.Number = id

    fmt.Printf("Новая посылка № %d на адрес %s от клиента с идентификатором %d зарегистрирована %s\n",
        p.Number, p.Address, p.Client, p.CreatedAt)

    return p, nil
}

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

func (s ParcelService) ChangeAddress(number int, address string) error {
    return s.store.SetAddress(number, address)
}

func (s ParcelService) Delete(number int) error {
    return s.store.Delete(number)
}

func main() {
<<<<<<< HEAD
	// настройте подключение к БД
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	store := NewParcelStore(db)
	service := NewParcelService(store)
=======
    db, err := sql.Open("sqlite3", "tracker.db")
    if err != nil {
        fmt.Println(err)
        return
    }
    defer db.Close()

    store := NewParcelStore(db)
    service := NewParcelService(store)
>>>>>>> 9a4e42ef888a45772eb8b234b1f14c71399ff43f

    client := 1
    address := "Псков, д. Пушкина, ул. Колотушкина, д. 5"
    p, err := service.Register(client, address)
    if err != nil {
        fmt.Println(err)
        return
    }

    newAddress := "Саратов, д. Верхние Зори, ул. Козлова, д. 25"
    err = service.ChangeAddress(p.Number, newAddress)
    if err != nil {
        fmt.Println(err)
        return
    }

    err = service.NextStatus(p.Number)
    if err != nil {
        fmt.Println(err)
        return
    }

    err = service.PrintClientParcels(client)
    if err != nil {
        fmt.Println(err)
        return
    }

    err = service.Delete(p.Number)
    if err != nil {
        fmt.Println(err)
        return
    }

    err = service.PrintClientParcels(client)
    if err != nil {
        fmt.Println(err)
        return
    }

    p, err = service.Register(client, address)
    if err != nil {
        fmt.Println(err)
        return
    }

    err = service.Delete(p.Number)
    if err != nil {
        fmt.Println(err)
        return
    }

    err = service.PrintClientParcels(client)
    if err != nil {
        fmt.Println(err)
        return
    }
}