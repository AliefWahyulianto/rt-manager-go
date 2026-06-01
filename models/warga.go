package models

import (
	"encoding/json"
	"strconv"

	"github.com/dgraph-io/badger/v4"
)

type Warga struct {
	ID     int    `json:"id"`
	Nama   string `json:"nama"`
	NIK    string `json:"nik"`
	Alamat string `json:"alamat"`
	Status string `json:"status"`
}

type WargaModel struct {
	DB *badger.DB
}

func NewWargaModel(db *badger.DB) *WargaModel {
	return &WargaModel{DB: db}
}

// GetAll mengambil semua data warga
func (m *WargaModel) GetAll() ([]Warga, error) {
	var wargaList []Warga

	err := m.DB.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.Prefix = []byte("warga_")
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			err := item.Value(func(val []byte) error {
				var w Warga
				if err := json.Unmarshal(val, &w); err != nil {
					return err
				}
				wargaList = append(wargaList, w)
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})

	return wargaList, err
}

// Search mencari warga berdasarkan nama atau alamat
func (m *WargaModel) Search(keyword string) ([]Warga, error) {
	all, err := m.GetAll()
	if err != nil {
		return nil, err
	}

	var result []Warga
	for _, w := range all {
		if contains(w.Nama, keyword) || contains(w.Alamat, keyword) {
			result = append(result, w)
		}
	}
	return result, nil
}

// GetByID mengambil satu warga berdasarkan ID
func (m *WargaModel) GetByID(id int) (*Warga, error) {
	var warga Warga
	key := []byte("warga_" + strconv.Itoa(id))

	err := m.DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, &warga)
		})
	})

	if err != nil {
		return nil, err
	}
	return &warga, nil
}

// Insert menambah warga baru
func (m *WargaModel) Insert(w Warga) error {
	// Generate ID baru
	count, _ := m.Count()
	w.ID = count + 1

	key := []byte("warga_" + strconv.Itoa(w.ID))
	value, err := json.Marshal(w)
	if err != nil {
		return err
	}

	return m.DB.Update(func(txn *badger.Txn) error {
		return txn.Set(key, value)
	})
}

// Update mengubah data warga
func (m *WargaModel) Update(w Warga) error {
	key := []byte("warga_" + strconv.Itoa(w.ID))
	value, err := json.Marshal(w)
	if err != nil {
		return err
	}

	return m.DB.Update(func(txn *badger.Txn) error {
		return txn.Set(key, value)
	})
}

// Delete menghapus warga
func (m *WargaModel) Delete(id int) error {
	key := []byte("warga_" + strconv.Itoa(id))
	return m.DB.Update(func(txn *badger.Txn) error {
		return txn.Delete(key)
	})
}

// Count menghitung total warga
func (m *WargaModel) Count() (int, error) {
	count := 0
	err := m.DB.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.Prefix = []byte("warga_")
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			count++
		}
		return nil
	})
	return count, err
}

// Helper function untuk pencarian substring
func contains(s, substr string) bool {
	return len(substr) == 0 || len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}