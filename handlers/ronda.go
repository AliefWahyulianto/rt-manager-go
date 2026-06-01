package handlers

import (
	"net/http"
	"rt-manager/models"
	"strconv"
	"strings"
	"text/template"
	"time"
)

type RondaHandler struct {
	RondaModel *models.RondaModel
	WargaModel *models.WargaModel
}

func NewRondaHandler(rondaModel *models.RondaModel, wargaModel *models.WargaModel) *RondaHandler {
	return &RondaHandler{
		RondaModel: rondaModel,
		WargaModel: wargaModel,
	}
}

// Index menampilkan daftar jadwal ronda
func (h *RondaHandler) Index(w http.ResponseWriter, r *http.Request) {
	rondaList, err := h.RondaModel.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Balik urutan (yang terbaru di atas)
	for i, j := 0, len(rondaList)-1; i < j; i, j = i+1, j-1 {
		rondaList[i], rondaList[j] = rondaList[j], rondaList[i]
	}

	funcMap := template.FuncMap{
		"add": func(a, b int) int { return a + b },
		"formatTanggal": func(tanggal string) string {
			t, _ := time.Parse("2006-01-02", tanggal)
			return t.Format("02 Jan 2006")
		},
	}

	data := map[string]interface{}{
		"RondaList":  rondaList,
		"ActivePage": "ronda",
	}

	tmpl := template.Must(template.New("layout.html").Funcs(funcMap).ParseFiles("templates/layout.html", "templates/ronda.html"))
	tmpl.Execute(w, data)
}

// Tambah menampilkan form tambah jadwal
func (h *RondaHandler) Tambah(w http.ResponseWriter, r *http.Request) {
	wargaList, err := h.WargaModel.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Method == "POST" {
		wargaID, _ := strconv.Atoi(r.FormValue("warga_id"))
		tanggal := r.FormValue("tanggal")
		shift := r.FormValue("shift")
		catatan := r.FormValue("catatan")

		// Ambil nama warga
		warga, _ := h.WargaModel.GetByID(wargaID)
		namaWarga := ""
		if warga != nil {
			namaWarga = warga.Nama
		}

		rondaBaru := models.Ronda{
			WargaID:   wargaID,
			NamaWarga: namaWarga,
			Tanggal:   tanggal,
			Shift:     shift,
			Status:    "Terjadwal",
			Catatan:   catatan,
		}

		err := h.RondaModel.Insert(rondaBaru)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/ronda", http.StatusSeeOther)
		return
	}

	data := map[string]interface{}{
		"WargaList":  wargaList,
		"ActivePage": "ronda",
	}
	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/ronda_tambah.html"))
	tmpl.Execute(w, data)
}

// TandaiSelesai mengubah status jadwal menjadi Selesai
func (h *RondaHandler) TandaiSelesai(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	idStr := pathParts[len(pathParts)-1]
	id, _ := strconv.Atoi(idStr)

	err := h.RondaModel.UpdateStatus(id, "Selesai")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/ronda", http.StatusSeeOther)
}

// Edit mengedit jadwal
func (h *RondaHandler) Edit(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	idStr := pathParts[len(pathParts)-1]
	id, _ := strconv.Atoi(idStr)

	wargaList, _ := h.WargaModel.GetAll()

	if r.Method == "POST" {
		wargaID, _ := strconv.Atoi(r.FormValue("warga_id"))
		tanggal := r.FormValue("tanggal")
		shift := r.FormValue("shift")
		status := r.FormValue("status")
		catatan := r.FormValue("catatan")

		warga, _ := h.WargaModel.GetByID(wargaID)
		namaWarga := ""
		if warga != nil {
			namaWarga = warga.Nama
		}

		rondaUpdate, _ := h.RondaModel.GetByID(id)
		rondaUpdate.WargaID = wargaID
		rondaUpdate.NamaWarga = namaWarga
		rondaUpdate.Tanggal = tanggal
		rondaUpdate.Shift = shift
		rondaUpdate.Status = status
		rondaUpdate.Catatan = catatan

		h.RondaModel.Update(*rondaUpdate)
		http.Redirect(w, r, "/ronda", http.StatusSeeOther)
		return
	}

	ronda, err := h.RondaModel.GetByID(id)
	if err != nil {
		http.Error(w, "Data tidak ditemukan", http.StatusNotFound)
		return
	}

	data := map[string]interface{}{
		"Ronda":      ronda,
		"WargaList":  wargaList,
		"ActivePage": "ronda",
	}
	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/ronda_edit.html"))
	tmpl.Execute(w, data)
}

// Hapus menghapus jadwal
func (h *RondaHandler) Hapus(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	idStr := pathParts[len(pathParts)-1]
	id, _ := strconv.Atoi(idStr)

	err := h.RondaModel.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/ronda", http.StatusSeeOther)
}
