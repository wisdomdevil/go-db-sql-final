package main

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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

// TestAddGetDelete проверяет добавление, получение и удаление посылки
func TestAddGetDelete(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	parcelID, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotEmpty(t, parcelID)

	// get
	gotParcel, err := store.Get(parcelID)
	require.NoError(t, err)
	gotParcel.Number = parcel.Number
	assert.Equal(t, parcel, gotParcel)

	// delete
	err = store.Delete(parcelID)
	require.NoError(t, err)
	_, err = store.Get(parcelID)
	require.ErrorIs(t, err, sql.ErrNoRows)
}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	parcelID, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotEmpty(t, parcelID)

	// set address
	newAddress := "new test address"
	err = store.SetAddress(parcelID, newAddress)
	require.NoError(t, err)

	// check
	gotParcel, err := store.Get(parcelID)
	require.NoError(t, err)
	assert.Equal(t, newAddress, gotParcel.Address)
}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	parcelID, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotEmpty(t, parcelID)

	// set status
	err = store.SetStatus(parcelID, ParcelStatusSent)
	require.NoError(t, err)

	// check
	gotParcel, err := store.Get(parcelID)
	require.NoError(t, err)
	assert.Equal(t, ParcelStatusSent, gotParcel.Status)
}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)

	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	parcelMap := map[int]Parcel{}

	// задаём всем посылкам один и тот же идентификатор клиента
	client := randRange.Intn(10_000_000)
	parcels[0].Client = client
	parcels[1].Client = client
	parcels[2].Client = client

	// add
	for i := 0; i < len(parcels); i++ {
		id, err := store.Add(parcels[i])
		require.NoError(t, err)
		assert.Greater(t, id, 0)

		// обновляем идентификатор добавленной у посылки
		parcels[i].Number = id

		// сохраняем добавленную посылку в структуру map, чтобы её можно было легко достать по идентификатору посылки
		parcelMap[id] = parcels[i]
	}

	// get by client
	storedParcels, err := store.GetByClient(client)
	require.NoError(t, err)
	assert.ElementsMatch(t, parcels, storedParcels)

	// check
	for _, parcel := range storedParcels {
		assert.Contains(t, parcelMap, parcel.Number)
		assert.Equal(t, parcelMap[parcel.Number], parcel)
	}
}
