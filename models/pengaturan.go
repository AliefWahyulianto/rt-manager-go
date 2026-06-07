package models

import (
	"encoding/json"

	"github.com/dgraph-io/badger/v4"
)

type Pengaturan struct {
	ID           int    `json:"id"`
	NamaRT       string `json:"nama_rt"`
	NamaKetua    string `json:"nama_ketua"`
	AlamatKantor string `json:"alamat_kantor"`
	NoTelepon    string `json:"no_telepon"`
	Email        string `json:"email"`
	Website      string `json:"website"`
	CreatedAt    string `json:"created_at"`
}

type PengaturanModel struct {
	DB *badger.DB
}

func NewPengaturanModel(db *badger.DB) *PengaturanModel {
	return &PengaturanModel{DB: db}
}

// Get mengambil data pengaturan (hanya 1 record)
func (m *PengaturanModel) Get() (*Pengaturan, error) {
	var pengaturan Pengaturan
	key := []byte("pengaturan_1")

	err := m.DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, &pengaturan)
		})
	})

	if err != nil {
		return nil, err
	}
	return &pengaturan, nil
}

// InsertOrUpdate menyimpan data pengaturan
func (m *PengaturanModel) InsertOrUpdate(p Pengaturan) error {
	p.ID = 1
	key := []byte("pengaturan_1")
	value, err := json.Marshal(p)
	if err != nil {
		return err
	}

	return m.DB.Update(func(txn *badger.Txn) error {
		return txn.Set(key, value)
	})
}
