package handlers

import (
	"encoding/json"
	"net/http"
	"rt-manager/models"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/jung-kurt/gofpdf/v2"
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

	// ========== TAMBAHKAN INI ==========
	// Ambil pengumuman aktif untuk ditampilkan di dashboard
	pengumumanModel := models.NewPengumumanModel(h.Model.DB)
	pengumumanAktif, _ := pengumumanModel.GetAktif()

	// Ambil total iuran terkumpul
	iuranModel := models.NewIuranModel(h.Model.DB)
	totalIuran, _ := iuranModel.TotalTerkumpul()
	// ===================================

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
		"truncate": func(s string, n int) string {
			if len(s) > n {
				return s[:n] + "..."
			}
			return s
		},
		"first": func(n int, list []models.Pengumuman) []models.Pengumuman { // <-- TAMBAHKAN INI
			if n > len(list) {
				return list
			}
			return list[:n]
		},
		"formatTanggal": func(tanggal string) string { // <-- TAMBAHKAN JUGA INI
			t, _ := time.Parse("2006-01-02", tanggal)
			return t.Format("02 Jan 2006")
		},
	}

	data := map[string]interface{}{
		"WargaList":       wargaList,
		"TotalWarga":      totalWarga,
		"Keyword":         keyword,
		"ActivePage":      "dashboard",
		"PengumumanAktif": pengumumanAktif, // TAMBAHKAN
		"TotalIuran":      totalIuran,      // TAMBAHKAN (untuk statistik)
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

// ExportPDF mengexport data warga ke PDF
func (h *WargaHandler) ExportPDF(w http.ResponseWriter, r *http.Request) {
	// Ambil semua data warga
	wargaList, err := h.Model.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Buat PDF baru
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)

	// Header
	pdf.Cell(40, 10, "Laporan Data Warga RT 05")
	pdf.Ln(12)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(40, 10, "Tanggal Cetak: "+time.Now().Format("02 Jan 2006 15:04:05"))
	pdf.Ln(15)

	// Header Tabel
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(15, 8, "No")
	pdf.Cell(60, 8, "Nama")
	pdf.Cell(50, 8, "NIK")
	pdf.Cell(50, 8, "Alamat")
	pdf.Cell(30, 8, "Status")
	pdf.Ln(8)

	// Data Tabel
	pdf.SetFont("Arial", "", 9)
	for i, w := range wargaList {
		pdf.Cell(15, 7, strconv.Itoa(i+1))
		pdf.Cell(60, 7, w.Nama)
		pdf.Cell(50, 7, w.NIK)
		pdf.Cell(50, 7, w.Alamat)
		pdf.Cell(30, 7, w.Status)
		pdf.Ln(7)
	}

	// Footer Total
	pdf.Ln(10)
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(40, 10, "Total Warga: "+strconv.Itoa(len(wargaList)))

	// Set header response
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=laporan_warga.pdf")

	// Output PDF
	err = pdf.Output(w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
