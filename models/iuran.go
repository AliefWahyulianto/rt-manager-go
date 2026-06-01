package models

import (
	"encoding/json"
	"strconv"

	"github.com/dgraph-io/badger/v4"
)

type Iuran struct {
	ID        int    `json:"id"`
	WargaID   int    `json:"warga_id"`
	NamaWarga string `json:"nama_warga"`
	Bulan     string `json:"bulan"`
	Tahun     int    `json:"tahun"`
	Nominal   int    `json:"nominal"`
	Status    string `json:"status"` // "Lunas" atau "Belum"
}

type IuranModel struct {
	DB *badger.DB
}

func NewIuranModel(db *badger.DB) *IuranModel {
	return &IuranModel{DB: db}
}

// GetAll mengambil semua iuran
func (m *IuranModel) GetAll() ([]Iuran, error) {
	var iuranList []Iuran

	err := m.DB.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.Prefix = []byte("iuran_")
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			err := item.Value(func(val []byte) error {
				var i Iuran
				if err := json.Unmarshal(val, &i); err != nil {
					return err
				}
				iuranList = append(iuranList, i)
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})

	return iuranList, err
}

// GetByWargaID mengambil iuran berdasarkan ID warga
func (m *IuranModel) GetByWargaID(wargaID int) ([]Iuran, error) {
	all, err := m.GetAll()
	if err != nil {
		return nil, err
	}

	var result []Iuran
	for _, i := range all {
		if i.WargaID == wargaID {
			result = append(result, i)
		}
	}
	return result, nil
}

// Insert menambah iuran baru
func (m *IuranModel) Insert(i Iuran) error {
	// Generate ID baru
	all, _ := m.GetAll()
	i.ID = len(all) + 1

	key := []byte("iuran_" + strconv.Itoa(i.ID))
	value, err := json.Marshal(i)
	if err != nil {
		return err
	}

	return m.DB.Update(func(txn *badger.Txn) error {
		return txn.Set(key, value)
	})
}

// Update mengubah status iuran
func (m *IuranModel) Update(i Iuran) error {
	key := []byte("iuran_" + strconv.Itoa(i.ID))
	value, err := json.Marshal(i)
	if err != nil {
		return err
	}

	return m.DB.Update(func(txn *badger.Txn) error {
		return txn.Set(key, value)
	})
}

// Delete menghapus iuran
func (m *IuranModel) Delete(id int) error {
	key := []byte("iuran_" + strconv.Itoa(id))
	return m.DB.Update(func(txn *badger.Txn) error {
		return txn.Delete(key)
	})
}

// TotalTerkumpul menghitung total nominal iuran yang sudah lunas
func (m *IuranModel) TotalTerkumpul() (int, error) {
	all, err := m.GetAll()
	if err != nil {
		return 0, err
	}

	total := 0
	for _, i := range all {
		if i.Status == "Lunas" {
			total += i.Nominal
		}
	}
	return total, nil
}

// GetByID mengambil satu iuran berdasarkan ID
func (m *IuranModel) GetByID(id int) (*Iuran, error) {
	var iuran Iuran
	key := []byte("iuran_" + strconv.Itoa(id))

	err := m.DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, &iuran)
		})
	})

	if err != nil {
		return nil, err
	}
	return &iuran, nil
}
