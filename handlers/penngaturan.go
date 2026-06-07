package handlers

import (
	"net/http"
	"rt-manager/models"
	"text/template"
	"time"
)

type PengaturanHandler struct {
	Model *models.PengaturanModel
}

func NewPengaturanHandler(model *models.PengaturanModel) *PengaturanHandler {
	return &PengaturanHandler{
		Model: model,
	}
}

// Index menampilkan halaman pengaturan
func (h *PengaturanHandler) Index(w http.ResponseWriter, r *http.Request) {
	pengaturan, err := h.Model.Get()
	if err != nil {
		// Jika belum ada data, buat data default
		pengaturan = &models.Pengaturan{
			NamaRT:       "RT 05",
			NamaKetua:    "Bambang Wijaya",
			AlamatKantor: "Jl. Mawar No. 12, RT 05/RW 03",
			NoTelepon:    "08123456789",
			Email:        "rt05@gmail.com",
			Website:      "",
		}
	}

	data := map[string]interface{}{
		"Pengaturan": pengaturan,
		"ActivePage": "pengaturan",
	}

	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/pengaturan.html"))
	tmpl.Execute(w, data)
}

// Edit mengedit pengaturan
func (h *PengaturanHandler) Edit(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		namaRT := r.FormValue("nama_rt")
		namaKetua := r.FormValue("nama_ketua")
		alamatKantor := r.FormValue("alamat_kantor")
		noTelepon := r.FormValue("no_telepon")
		email := r.FormValue("email")
		website := r.FormValue("website")

		pengaturan := models.Pengaturan{
			NamaRT:       namaRT,
			NamaKetua:    namaKetua,
			AlamatKantor: alamatKantor,
			NoTelepon:    noTelepon,
			Email:        email,
			Website:      website,
			CreatedAt:    time.Now().Format("2006-01-02 15:04:05"),
		}

		err := h.Model.InsertOrUpdate(pengaturan)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/pengaturan", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/pengaturan", http.StatusSeeOther)
}
