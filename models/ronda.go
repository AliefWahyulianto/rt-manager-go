package models

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/dgraph-io/badger/v4"
)

type Ronda struct {
	ID        int    `json:"id"`
	WargaID   int    `json:"warga_id"`
	NamaWarga string `json:"nama_warga"`
	Tanggal   string `json:"tanggal"` // format: 2026-06-01
	Shift     string `json:"shift"`   // Malam 1 (19-21), Malam 2 (21-23), Dini hari (23-01)
	Status    string `json:"status"`  // "Terjadwal", "Selesai", "Batal"
	Catatan   string `json:"catatan"` // opsional
	CreatedAt string `json:"created_at"`
}

type RondaModel struct {
	DB *badger.DB
}

func NewRondaModel(db *badger.DB) *RondaModel {
	return &RondaModel{DB: db}
}

// GetAll mengambil semua jadwal ronda
func (m *RondaModel) GetAll() ([]Ronda, error) {
	var rondaList []Ronda

	err := m.DB.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.Prefix = []byte("ronda_")
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			err := item.Value(func(val []byte) error {
				var r Ronda
				if err := json.Unmarshal(val, &r); err != nil {
					return err
				}
				rondaList = append(rondaList, r)
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})

	return rondaList, err
}

// GetByID mengambil satu jadwal berdasarkan ID
func (m *RondaModel) GetByID(id int) (*Ronda, error) {
	var ronda Ronda
	key := []byte("ronda_" + strconv.Itoa(id))

	err := m.DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, &ronda)
		})
	})

	if err != nil {
		return nil, err
	}
	return &ronda, nil
}

// GetByTanggal mengambil jadwal berdasarkan tanggal
func (m *RondaModel) GetByTanggal(tanggal string) ([]Ronda, error) {
	all, err := m.GetAll()
	if err != nil {
		return nil, err
	}

	var result []Ronda
	for _, r := range all {
		if r.Tanggal == tanggal {
			result = append(result, r)
		}
	}
	return result, nil
}

// GetJadwalMalamIni mengambil jadwal untuk hari ini
func (m *RondaModel) GetJadwalMalamIni() ([]Ronda, error) {
	// Format tanggal sekarang: 2026-06-01
	tanggalSekarang := getCurrentDate()
	return m.GetByTanggal(tanggalSekarang)
}

// Insert menambah jadwal baru
func (m *RondaModel) Insert(r Ronda) error {
	// Generate ID baru
	all, _ := m.GetAll()
	r.ID = len(all) + 1
	r.CreatedAt = getCurrentDateTime()

	key := []byte("ronda_" + strconv.Itoa(r.ID))
	value, err := json.Marshal(r)
	if err != nil {
		return err
	}

	return m.DB.Update(func(txn *badger.Txn) error {
		return txn.Set(key, value)
	})
}

// Update mengubah jadwal
func (m *RondaModel) Update(r Ronda) error {
	key := []byte("ronda_" + strconv.Itoa(r.ID))
	value, err := json.Marshal(r)
	if err != nil {
		return err
	}

	return m.DB.Update(func(txn *badger.Txn) error {
		return txn.Set(key, value)
	})
}

// UpdateStatus mengubah status jadwal
func (m *RondaModel) UpdateStatus(id int, status string) error {
	r, err := m.GetByID(id)
	if err != nil {
		return err
	}
	r.Status = status
	return m.Update(*r)
}

// Delete menghapus jadwal
func (m *RondaModel) Delete(id int) error {
	key := []byte("ronda_" + strconv.Itoa(id))
	return m.DB.Update(func(txn *badger.Txn) error {
		return txn.Delete(key)
	})
}

// Count menghitung total jadwal
func (m *RondaModel) Count() (int, error) {
	all, err := m.GetAll()
	return len(all), err
}

// Helper functions
func getCurrentDate() string {
	// Format: 2026-06-01
	now := time.Now()
	return now.Format("2006-01-02")
}

func getCurrentDateTime() string {
	now := time.Now()
	return now.Format("2006-01-02 15:04:05")
}
