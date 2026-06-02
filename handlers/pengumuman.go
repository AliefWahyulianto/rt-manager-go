package handlers

import (
	"net/http"
	"rt-manager/models"
	"strconv"
	"strings"
	"text/template"
	"time"
)

type PengumumanHandler struct {
	Model *models.PengumumanModel
}

func NewPengumumanHandler(model *models.PengumumanModel) *PengumumanHandler {
	return &PengumumanHandler{
		Model: model,
	}
}

// Index menampilkan daftar pengumuman
func (h *PengumumanHandler) Index(w http.ResponseWriter, r *http.Request) {
	list, err := h.Model.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Balik urutan (yang baru di atas)
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}

	funcMap := template.FuncMap{
		"add": func(a, b int) int { return a + b },
		"formatTanggal": func(tanggal string) string {
			t, _ := time.Parse("2006-01-02", tanggal)
			return t.Format("02 Jan 2006")
		},
		"truncate": func(s string, n int) string {
			if len(s) > n {
				return s[:n] + "..."
			}
			return s
		},
	}

	data := map[string]interface{}{
		"PengumumanList": list,
		"ActivePage":     "pengumuman",
	}

	tmpl := template.Must(template.New("layout.html").Funcs(funcMap).ParseFiles("templates/layout.html", "templates/pengumuman.html"))
	tmpl.Execute(w, data)
}

// Tambah menampilkan form tambah pengumuman
func (h *PengumumanHandler) Tambah(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		judul := r.FormValue("judul")
		isi := r.FormValue("isi")
		tanggal := r.FormValue("tanggal")
		status := r.FormValue("status")

		pengumumanBaru := models.Pengumuman{
			Judul:   judul,
			Isi:     isi,
			Tanggal: tanggal,
			Status:  status,
		}

		err := h.Model.Insert(pengumumanBaru)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/pengumuman", http.StatusSeeOther)
		return
	}

	funcMap := template.FuncMap{
		"now": func() string {
			return time.Now().Format("2006-01-02")
		},
	}

	data := map[string]interface{}{
		"ActivePage": "pengumuman",
	}

	tmpl := template.Must(template.New("layout.html").Funcs(funcMap).ParseFiles("templates/layout.html", "templates/pengumuman_tambah.html"))
	tmpl.Execute(w, data)
}

// Edit mengedit pengumuman
func (h *PengumumanHandler) Edit(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	idStr := pathParts[len(pathParts)-1]
	id, _ := strconv.Atoi(idStr)

	if r.Method == "POST" {
		judul := r.FormValue("judul")
		isi := r.FormValue("isi")
		tanggal := r.FormValue("tanggal")
		status := r.FormValue("status")

		pengumuman, _ := h.Model.GetByID(id)
		pengumuman.Judul = judul
		pengumuman.Isi = isi
		pengumuman.Tanggal = tanggal
		pengumuman.Status = status

		h.Model.Update(*pengumuman)
		http.Redirect(w, r, "/pengumuman", http.StatusSeeOther)
		return
	}

	pengumuman, err := h.Model.GetByID(id)
	if err != nil {
		http.Error(w, "Data tidak ditemukan", http.StatusNotFound)
		return
	}

	funcMap := template.FuncMap{
		"now": func() string {
			return time.Now().Format("2006-01-02")
		},
	}

	data := map[string]interface{}{
		"Pengumuman": pengumuman,
		"ActivePage": "pengumuman",
	}

	tmpl := template.Must(template.New("layout.html").Funcs(funcMap).ParseFiles("templates/layout.html", "templates/pengumuman_edit.html"))
	tmpl.Execute(w, data)
}

// Hapus menghapus pengumuman
func (h *PengumumanHandler) Hapus(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	idStr := pathParts[len(pathParts)-1]
	id, _ := strconv.Atoi(idStr)

	err := h.Model.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/pengumuman", http.StatusSeeOther)
}

// Arsipkan mengubah status menjadi Arsip
func (h *PengumumanHandler) Arsipkan(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	idStr := pathParts[len(pathParts)-1]
	id, _ := strconv.Atoi(idStr)

	pengumuman, _ := h.Model.GetByID(id)
	pengumuman.Status = "Arsip"
	h.Model.Update(*pengumuman)

	http.Redirect(w, r, "/pengumuman", http.StatusSeeOther)
}
