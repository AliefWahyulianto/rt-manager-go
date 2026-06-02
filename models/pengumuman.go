package models

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/dgraph-io/badger/v4"
)

type Pengumuman struct {
	ID        int    `json:"id"`
	Judul     string `json:"judul"`
	Isi       string `json:"isi"`
	Tanggal   string `json:"tanggal"`
	Status    string `json:"status"` // "Aktif" atau "Arsip"
	CreatedAt string `json:"created_at"`
}

type PengumumanModel struct {
	DB *badger.DB
}

func NewPengumumanModel(db *badger.DB) *PengumumanModel {
	return &PengumumanModel{DB: db}
}

// GetAll mengambil semua pengumuman
func (m *PengumumanModel) GetAll() ([]Pengumuman, error) {
	var list []Pengumuman

	err := m.DB.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.Prefix = []byte("pengumuman_")
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			err := item.Value(func(val []byte) error {
				var p Pengumuman
				if err := json.Unmarshal(val, &p); err != nil {
					return err
				}
				list = append(list, p)
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})

	return list, err
}

// GetAktif mengambil pengumuman yang statusnya Aktif
func (m *PengumumanModel) GetAktif() ([]Pengumuman, error) {
	all, err := m.GetAll()
	if err != nil {
		return nil, err
	}

	var result []Pengumuman
	for _, p := range all {
		if p.Status == "Aktif" {
			result = append(result, p)
		}
	}
	return result, nil
}

// GetByID mengambil satu pengumuman berdasarkan ID
func (m *PengumumanModel) GetByID(id int) (*Pengumuman, error) {
	var p Pengumuman
	key := []byte("pengumuman_" + strconv.Itoa(id))

	err := m.DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, &p)
		})
	})

	if err != nil {
		return nil, err
	}
	return &p, nil
}

// Insert menambah pengumuman baru
func (m *PengumumanModel) Insert(p Pengumuman) error {
	all, _ := m.GetAll()
	p.ID = len(all) + 1
	p.CreatedAt = time.Now().Format("2006-01-02 15:04:05")

	// Format tanggal untuk display
	if p.Tanggal == "" {
		p.Tanggal = time.Now().Format("2006-01-02")
	}
	if p.Status == "" {
		p.Status = "Aktif"
	}

	key := []byte("pengumuman_" + strconv.Itoa(p.ID))
	value, err := json.Marshal(p)
	if err != nil {
		return err
	}

	return m.DB.Update(func(txn *badger.Txn) error {
		return txn.Set(key, value)
	})
}

// Update mengubah pengumuman
func (m *PengumumanModel) Update(p Pengumuman) error {
	key := []byte("pengumuman_" + strconv.Itoa(p.ID))
	value, err := json.Marshal(p)
	if err != nil {
		return err
	}

	return m.DB.Update(func(txn *badger.Txn) error {
		return txn.Set(key, value)
	})
}

// Delete menghapus pengumuman
func (m *PengumumanModel) Delete(id int) error {
	key := []byte("pengumuman_" + strconv.Itoa(id))
	return m.DB.Update(func(txn *badger.Txn) error {
		return txn.Delete(key)
	})
}

// Count menghitung total pengumuman
func (m *PengumumanModel) Count() (int, error) {
	all, err := m.GetAll()
	return len(all), err
}
