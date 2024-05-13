package database

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
)

type DB struct {
	Path  string
	Mutex *sync.RWMutex
}

type DB_struct struct {
	Chirps map[string]Chirp `json:"chirps"`
}

type Chirp struct {
	Id   int
	Body string
}

const (
    database_name = "database.json"
)

func (db *DB) EnsureDB() error {
    _, err := os.Open(database_name)
    if err != os.ErrNotExist {
        err = os.WriteFile(database_name, nil, 667)
        return err
    } 
    return nil
}


func NewDB(path string) (*DB, error) {      
    err := os.WriteFile(database_name, nil, 667)
    if err != nil {
        return &DB{}, err
    }

    return &DB {
        Path: "./database/database.json",
        Mutex: &sync.RWMutex{},
    }, nil
}

func (db *DB) CreateChirp(body string) (Chirp, error) {
    db.Mutex.Lock() 
    defer db.Mutex.Unlock()

    chirp_id := 0
    chirp  := Chirp {
        Id: chirp_id,
        Body: body,
    }

    db_struct, err := db.loadDB()
    if err != nil {
        return Chirp{}, err 
    }

    db_struct.Chirps[fmt.Sprint(chirp_id)] = chirp
    chirp_id++

    err = db.WriteDB(db_struct) 
    if err != nil {
        return Chirp{}, nil
    }

    return chirp, nil
}

func (db *DB) GetChirps() ([]Chirp, error) {
    db.Mutex.Lock()
    defer db.Mutex.Unlock()

    var all_chirps []Chirp
    db_struct, err := db.loadDB()
    if err != nil {
        return nil, err 
    }

    for _, chirp := range db_struct.Chirps {
        all_chirps = append(all_chirps, chirp) 
    }

    return all_chirps, nil
}


func (db *DB) loadDB() (DB_struct, error) {
    db.Mutex.Lock()
    defer db.Mutex.Unlock()
    data, err := os.ReadFile(db.Path)
    if err != nil {
        return DB_struct{}, err
    }

    db_struct := DB_struct{}

    err = json.Unmarshal(data, db_struct) 
    if err != nil {
        log.Printf("Couldn't Unmarshal data: %v", err)
        return DB_struct{}, err
    }

    return db_struct, nil
}

func (db *DB) WriteDB(db_struct DB_struct) error {
    db.Mutex.Lock()
    defer db.Mutex.Unlock()

    file, err := os.Open(db.Path)
    if err != nil {
        log.Printf("Couldn't open file: %v", err)
        return err
    }

    data, err := json.Marshal(db_struct)
    if err != nil {
        log.Printf("Couldn't marshal data: %v", err)
        return err
    }

    _, err = file.Write(data)
    if err != nil {
        log.Printf("Couldn't Write Data: %v", err)
        return err
    }

    return nil
}


