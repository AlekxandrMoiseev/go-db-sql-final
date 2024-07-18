package main

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3" // Импортируйте SQLite драйвер
	"github.com/stretchr/testify/require"
)

var (
	// randSource источник псевдо случайных чисел.
	// Для повышения уникальности в качестве seed
	// используется текущее время в unix формате (в виде числа)
	randSource = rand.NewSource(time.Now().UnixNano())
	// randRange использует randSource для генерации случайных чисел
	randRange = rand.New(randSource)
)

// getTestParcel возвращает тестовую посылку
func getTestParcel() Parcel {
	return Parcel{
		Client:    1000,
		Status:    ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

// setupDatabase инициализирует базу данных и создает таблицу parcel
func setupDatabase(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:") // Используем in-memory базу данных для тестов
	require.NoError(t, err)

	createTableSQL := `CREATE TABLE IF NOT EXISTS parcel (
        number INTEGER PRIMARY KEY AUTOINCREMENT,
        client INTEGER NOT NULL,
        status TEXT NOT NULL,
        address TEXT NOT NULL,
        created_at TEXT NOT NULL
    );`

	_, err = db.Exec(createTableSQL)
	require.NoError(t, err)

	return db
}

// TestAddGetDelete проверяет добавление, получение и удаление посылки
func TestAddGetDelete(t *testing.T) {
	db := setupDatabase(t)
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotZero(t, id)

	// Установите поле Number для сравнения
	parcel.Number = id

	// get
	addedParcel, err := store.Get(id)
	require.NoError(t, err)
	require.Equal(t, parcel, addedParcel)

	// delete
	err = store.Delete(id)
	require.NoError(t, err)

	// check deletion
	_, err = store.Get(id)
	require.Error(t, err)
}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {
	db := setupDatabase(t)
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotZero(t, id)

	// Установите поле Number для сравнения
	parcel.Number = id

	// set address
	newAddress := "new test address"
	err = store.SetAddress(id, newAddress)
	require.NoError(t, err)

	// Обновите поле Address для сравнения
	parcel.Address = newAddress

	// check
	updatedParcel, err := store.Get(id)
	require.NoError(t, err)
	require.Equal(t, parcel, updatedParcel)
}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {
	db := setupDatabase(t)
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotZero(t, id)

	// Установите поле Number для сравнения
	parcel.Number = id

	// set status
	newStatus := "sent"
	err = store.SetStatus(id, newStatus)
	require.NoError(t, err)

	// Обновите поле Status для сравнения
	parcel.Status = newStatus

	// check
	updatedParcel, err := store.Get(id)
	require.NoError(t, err)
	require.Equal(t, parcel, updatedParcel)
}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {
	db := setupDatabase(t)
	defer db.Close()

	store := NewParcelStore(db)
	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	parcelMap := map[int]Parcel{}
	client := randRange.Intn(10_000_000)
	for i := range parcels {
		parcels[i].Client = client
	}

	// add
	for _, parcel := range parcels {
		id, err := store.Add(parcel)
		require.NoError(t, err)
		require.NotZero(t, id)
		parcel.Number = id
		parcelMap[id] = parcel
	}

	// get by client
	storedParcels, err := store.GetByClient(client)
	require.NoError(t, err)
	require.Len(t, storedParcels, len(parcels))

	// check
	for _, parcel := range storedParcels {
		expectedParcel, exists := parcelMap[parcel.Number]
		require.True(t, exists)
		require.Equal(t, expectedParcel, parcel)
	}
}
