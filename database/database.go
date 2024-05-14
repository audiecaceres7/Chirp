package database

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"sync"
)

type DB struct {
	Path        string
	Mutex       *sync.RWMutex
}

type DB_struct struct {
	Chirps map[int]Chirp `json:"chirps"`
}

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

const (
	database_name = "database.json"
)

func (db *DB) CreateDB() error {
    db_struct := DB_struct{
        Chirps: make(map[int]Chirp),
    }
    return db.WriteDB(db_struct)
}

func (db *DB) EnsureDB() error {
    _, err := os.ReadFile(db.Path)
	if errors.Is(err, os.ErrNotExist) {
		return db.CreateDB()
	}
	return err
}

func NewDB(path string) (*DB, error) {
    db := &DB{
		Path:  "./database.json",
		Mutex: &sync.RWMutex{},
	}
    err := db.EnsureDB()
    return db, err
}

func (db *DB) CreateChirp(body string) (Chirp, error) {
    db_struct, err := db.LoadDB()
    if err != nil {
        return Chirp{},  nil
    }

    id := len(db_struct.Chirps) + 1
	chirp := Chirp{
		ID:   id,
		Body: body,
	}

	db_struct.Chirps[id] = chirp

	err = db.WriteDB(db_struct)
	if err != nil {
		return Chirp{}, nil
	}

	return chirp, nil
}

func (db *DB) GetChirps() ([]Chirp, error) {
	var all_chirps []Chirp
	db_struct, err := db.LoadDB()
	if err != nil {
		return nil, err
	}

	for _, chirp := range db_struct.Chirps {
		all_chirps = append(all_chirps, chirp)
	}

	return all_chirps, nil
}

func (db *DB) LoadDB() (DB_struct, error) {
	db_struct := DB_struct{}

	data, err := os.ReadFile(db.Path)
	if err != nil {
		return DB_struct{}, err
	}

	err = json.Unmarshal(data, &db_struct)
	if err != nil {
		log.Printf("Couldn't Unmarshal data: %v", err)
		return DB_struct{}, err
	}

	return db_struct, nil
}

func (db *DB) WriteDB(db_struct DB_struct) error {
	data, err := json.Marshal(db_struct)
	if err != nil {
		log.Printf("Couldn't marshal data: %v", err)
		return err
	}

	err = os.WriteFile(db.Path, data, 0600)
	if err != nil {
		log.Printf("Couldn't Write Data: %v", err)
		return err
	}

	return nil
}
