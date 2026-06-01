package handlers

import (
	"encoding/json"
	"net/http"
	"rt-manager/models"
	"strconv"
	"strings"
	"text/template"
)

type WargaHandler struct {
	Model *models.WargaModel
}

func NewWargaHandler(model *models.WargaModel) *WargaHandler {
	return &WargaHandler{
		Model: model,
	}
}

// Dashboard menampilkan halaman utama
func (h *WargaHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	keyword := r.URL.Query().Get("search")
	var wargaList []models.Warga
	var err error

	if keyword != "" {
		wargaList, err = h.Model.Search(keyword)
	} else {
		wargaList, err = h.Model.GetAll()
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	totalWarga, _ := h.Model.Count()

	funcMap := template.FuncMap{
		"add": func(a, b int) int { return a + b },
		"slice": func(s string, start, end int) string {
			if start >= len(s) {
				return ""
			}
			if end > len(s) {
				end = len(s)
			}
			return s[start:end]
		},
	}

	data := map[string]interface{}{
		"WargaList":  wargaList,
		"TotalWarga": totalWarga,
		"Keyword":    keyword,
		"ActivePage": "dashboard",
	}

	tmpl := template.Must(template.New("layout.html").Funcs(funcMap).ParseFiles("templates/layout.html", "templates/dashboard.html"))
	tmpl.Execute(w, data)
}

// Tambah menampilkan form tambah dan memproses submit
func (h *WargaHandler) Tambah(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		nama := r.FormValue("nama")
		nik := r.FormValue("nik")
		alamat := r.FormValue("alamat")
		status := r.FormValue("status")

		if status == "" {
			status = "Tetap"
		}

		wargaBaru := models.Warga{
			Nama:   nama,
			NIK:    nik,
			Alamat: alamat,
			Status: status,
		}

		err := h.Model.Insert(wargaBaru)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	data := map[string]interface{}{
		"ActivePage": "tambah",
	}
	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/tambah.html"))
	tmpl.Execute(w, data)
}

// Edit menampilkan form edit dan memproses submit
func (h *WargaHandler) Edit(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	idStr := pathParts[len(pathParts)-1]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID tidak valid", http.StatusBadRequest)
		return
	}

	if r.Method == "POST" {
		nama := r.FormValue("nama")
		nik := r.FormValue("nik")
		alamat := r.FormValue("alamat")
		status := r.FormValue("status")

		wargaUpdate := models.Warga{
			ID:     id,
			Nama:   nama,
			NIK:    nik,
			Alamat: alamat,
			Status: status,
		}

		err := h.Model.Update(wargaUpdate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	warga, err := h.Model.GetByID(id)
	if err != nil {
		http.Error(w, "Data tidak ditemukan", http.StatusNotFound)
		return
	}

	data := map[string]interface{}{
		"Warga":      warga,
		"ActivePage": "edit",
	}
	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/edit.html"))
	tmpl.Execute(w, data)
}

// Hapus menghapus data warga
func (h *WargaHandler) Hapus(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	idStr := pathParts[len(pathParts)-1]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID tidak valid", http.StatusBadRequest)
		return
	}

	err = h.Model.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// APIWarga untuk pencarian via AJAX
func (h *WargaHandler) APIWarga(w http.ResponseWriter, r *http.Request) {
	keyword := r.URL.Query().Get("search")
	var wargaList []models.Warga
	var err error

	if keyword != "" {
		wargaList, err = h.Model.Search(keyword)
	} else {
		wargaList, err = h.Model.GetAll()
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(wargaList)
}