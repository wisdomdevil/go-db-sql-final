package main

import (
	"database/sql"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

// Add - реализует добавление новой посылки в таблицу parcel
func (s ParcelStore) Add(p Parcel) (int, error) {
	res, err := s.db.Exec(`INSERT INTO parcel (client, status, address, created_at)
						   VALUES (:client, :status, :address, :created_at)`,
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("created_at", p.CreatedAt))
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// Get - реализует чтение строки из таблицы parcel по заданному number
func (s ParcelStore) Get(number int) (Parcel, error) {
	p := Parcel{}

	row := s.db.QueryRow(`SELECT number, client, status, address, created_at
						  FROM parcel WHERE number=:number`,
		sql.Named("number", number))

	if err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt); err != nil {
		return Parcel{}, err
	}

	return p, nil
}

// GetByClient - реализует чтение строк из таблицы parcel по заданному client
func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	var res []Parcel

	rows, err := s.db.Query(`SELECT number, client, status, address, created_at
							 FROM parcel WHERE client=:client`,
		sql.Named("client", client))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		p := Parcel{}
		if err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt); err != nil {
			return nil, err
		}
		res = append(res, p)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return res, nil
}

// SetStatus - реализует обновление статуса посылки в таблице parcel
func (s ParcelStore) SetStatus(number int, status string) error {
	_, err := s.db.Exec(`UPDATE parcel SET status = :status WHERE number = :number`,
		sql.Named("status", status),
		sql.Named("number", number))
	if err != nil { //  <-- через перенос выглядит читабельнее
		return err
	}

	return nil
}

// SetAddress - реализует обновление адреса посылки в таблице parcel
func (s ParcelStore) SetAddress(number int, address string) error {
	_, err := s.db.Exec(`UPDATE parcel SET address = :address
						 WHERE number = :number AND status = :status`,
		sql.Named("address", address),
		sql.Named("number", number),
		sql.Named("status", ParcelStatusRegistered))
	if err != nil {
		return err
	}

	return nil
}

// Delete - реализует удаление посылки из таблицы parcel
func (s ParcelStore) Delete(number int) error {
	_, err := s.db.Exec(`DELETE FROM parcel WHERE number = :number AND status = :status`,
		sql.Named("number", number),
		sql.Named("status", ParcelStatusRegistered))
	if err != nil {
		return err
	}

	return nil
}
