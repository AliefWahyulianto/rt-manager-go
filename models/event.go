package models

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/dgraph-io/badger/v4"
)

type Event struct {
	ID        int    `json:"id"`
	NamaEvent string `json:"nama_event"`
	Tanggal   string `json:"tanggal"` // format: 2026-06-01
	Lokasi    string `json:"lokasi"`
	Deskripsi string `json:"deskripsi"`
	Anggaran  int    `json:"anggaran"` // dalam Rupiah
	Status    string `json:"status"`   // "Akan Datang", "Berlangsung", "Selesai"
	CreatedAt string `json:"created_at"`
}

type EventModel struct {
	DB *badger.DB
}

func NewEventModel(db *badger.DB) *EventModel {
	return &EventModel{DB: db}
}

// GetAll mengambil semua event
func (m *EventModel) GetAll() ([]Event, error) {
	var eventList []Event

	err := m.DB.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.Prefix = []byte("event_")
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			err := item.Value(func(val []byte) error {
				var e Event
				if err := json.Unmarshal(val, &e); err != nil {
					return err
				}
				eventList = append(eventList, e)
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})

	return eventList, err
}

// GetByID mengambil satu event berdasarkan ID
func (m *EventModel) GetByID(id int) (*Event, error) {
	var event Event
	key := []byte("event_" + strconv.Itoa(id))

	err := m.DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, &event)
		})
	})

	if err != nil {
		return nil, err
	}
	return &event, nil
}

// GetEventAkanDatang mengambil event dengan status "Akan Datang" atau tanggal > sekarang
func (m *EventModel) GetEventAkanDatang() ([]Event, error) {
	all, err := m.GetAll()
	if err != nil {
		return nil, err
	}

	today := time.Now().Format("2006-01-02")
	var result []Event
	for _, e := range all {
		if e.Tanggal >= today && e.Status != "Selesai" {
			result = append(result, e)
		}
	}
	return result, nil
}

// Insert menambah event baru
func (m *EventModel) Insert(e Event) error {
	all, _ := m.GetAll()
	e.ID = len(all) + 1
	e.CreatedAt = time.Now().Format("2006-01-02 15:04:05")

	// Set status otomatis berdasarkan tanggal
	today := time.Now().Format("2006-01-02")
	if e.Tanggal < today {
		e.Status = "Selesai"
	} else if e.Tanggal == today {
		e.Status = "Berlangsung"
	} else {
		e.Status = "Akan Datang"
	}

	key := []byte("event_" + strconv.Itoa(e.ID))
	value, err := json.Marshal(e)
	if err != nil {
		return err
	}

	return m.DB.Update(func(txn *badger.Txn) error {
		return txn.Set(key, value)
	})
}

// Update mengubah event
func (m *EventModel) Update(e Event) error {
	// Update status otomatis
	today := time.Now().Format("2006-01-02")
	if e.Tanggal < today {
		e.Status = "Selesai"
	} else if e.Tanggal == today {
		e.Status = "Berlangsung"
	} else {
		e.Status = "Akan Datang"
	}

	key := []byte("event_" + strconv.Itoa(e.ID))
	value, err := json.Marshal(e)
	if err != nil {
		return err
	}

	return m.DB.Update(func(txn *badger.Txn) error {
		return txn.Set(key, value)
	})
}

// Delete menghapus event
func (m *EventModel) Delete(id int) error {
	key := []byte("event_" + strconv.Itoa(id))
	return m.DB.Update(func(txn *badger.Txn) error {
		return txn.Delete(key)
	})
}

// Count menghitung total event
func (m *EventModel) Count() (int, error) {
	all, err := m.GetAll()
	return len(all), err
}
