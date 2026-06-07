package handlers

import (
	"net/http"
	"rt-manager/models"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/jung-kurt/gofpdf/v2"
)

type IuranHandler struct {
	IuranModel *models.IuranModel
	WargaModel *models.WargaModel
}

func NewIuranHandler(iuranModel *models.IuranModel, wargaModel *models.WargaModel) *IuranHandler {
	return &IuranHandler{
		IuranModel: iuranModel,
		WargaModel: wargaModel,
	}
}

// Index menampilkan daftar iuran
func (h *IuranHandler) Index(w http.ResponseWriter, r *http.Request) {
	iuranList, err := h.IuranModel.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	totalTerkumpul, _ := h.IuranModel.TotalTerkumpul()

	funcMap := template.FuncMap{
		"formatRupiah": func(nominal int) string {
			return "Rp " + strconv.Itoa(nominal)
		},
		"add": func(a, b int) int { return a + b },
	}

	data := map[string]interface{}{
		"IuranList":      iuranList,
		"TotalTerkumpul": totalTerkumpul,
		"ActivePage":     "iuran",
	}

	tmpl := template.Must(template.New("layout.html").Funcs(funcMap).ParseFiles("templates/layout.html", "templates/iuran.html"))
	tmpl.Execute(w, data)
}

// Tambah menampilkan form tambah iuran
func (h *IuranHandler) Tambah(w http.ResponseWriter, r *http.Request) {
	wargaList, err := h.WargaModel.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Method == "POST" {
		wargaID, _ := strconv.Atoi(r.FormValue("warga_id"))
		bulan := r.FormValue("bulan")
		tahun, _ := strconv.Atoi(r.FormValue("tahun"))
		nominal, _ := strconv.Atoi(r.FormValue("nominal"))
		status := r.FormValue("status")

		// Ambil nama warga
		warga, _ := h.WargaModel.GetByID(wargaID)
		namaWarga := ""
		if warga != nil {
			namaWarga = warga.Nama
		}

		iuranBaru := models.Iuran{
			WargaID:   wargaID,
			NamaWarga: namaWarga,
			Bulan:     bulan,
			Tahun:     tahun,
			Nominal:   nominal,
			Status:    status,
		}

		err := h.IuranModel.Insert(iuranBaru)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/iuran", http.StatusSeeOther)
		return
	}

	data := map[string]interface{}{
		"WargaList":  wargaList,
		"ActivePage": "iuran",
	}
	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/iuran_tambah.html"))
	tmpl.Execute(w, data)
}

// TandaiLunas mengubah status iuran menjadi Lunas
func (h *IuranHandler) TandaiLunas(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	idStr := pathParts[len(pathParts)-1]
	id, _ := strconv.Atoi(idStr)

	iuran, err := h.IuranModel.GetByID(id)
	if err != nil {
		http.Error(w, "Data tidak ditemukan", http.StatusNotFound)
		return
	}

	iuran.Status = "Lunas"
	err = h.IuranModel.Update(*iuran)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/iuran", http.StatusSeeOther)
}

// Edit mengedit iuran
func (h *IuranHandler) Edit(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	idStr := pathParts[len(pathParts)-1]
	id, _ := strconv.Atoi(idStr)

	if r.Method == "POST" {
		status := r.FormValue("status")
		iuran, _ := h.IuranModel.GetByID(id)
		iuran.Status = status
		h.IuranModel.Update(*iuran)
		http.Redirect(w, r, "/iuran", http.StatusSeeOther)
		return
	}

	iuran, err := h.IuranModel.GetByID(id)
	if err != nil {
		http.Error(w, "Data tidak ditemukan", http.StatusNotFound)
		return
	}

	data := map[string]interface{}{
		"Iuran":      iuran,
		"ActivePage": "iuran",
	}
	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/iuran_edit.html"))
	tmpl.Execute(w, data)
}

// Hapus menghapus iuran
func (h *IuranHandler) Hapus(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	idStr := pathParts[len(pathParts)-1]
	id, _ := strconv.Atoi(idStr)

	err := h.IuranModel.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/iuran", http.StatusSeeOther)
}

// GetByID mengambil satu iuran (untuk internal use)
func (h *IuranHandler) GetByID(id int) (*models.Iuran, error) {
	return h.IuranModel.GetByID(id)
}

// ExportPDF mengexport data iuran ke PDF
func (h *IuranHandler) ExportPDF(w http.ResponseWriter, r *http.Request) {
	// Ambil semua data iuran
	iuranList, err := h.IuranModel.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Hitung total nominal iuran (lunas)
	totalTerkumpul, _ := h.IuranModel.TotalTerkumpul()

	// Hitung total nominal seluruh iuran (lunas + belum)
	totalKeseluruhan := 0
	for _, i := range iuranList {
		totalKeseluruhan += i.Nominal
	}

	// Buat PDF baru
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// ========== HEADER ==========
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Laporan Iuran Warga RT 05")
	pdf.Ln(12)

	pdf.SetFont("Arial", "", 10)
	pdf.Cell(40, 10, "Tanggal Cetak: "+time.Now().Format("02 Jan 2006 15:04:05"))
	pdf.Ln(15)

	// ========== HEADER TABEL ==========
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(15, 8, "No")
	pdf.Cell(55, 8, "Nama Warga")
	pdf.Cell(30, 8, "Bulan")
	pdf.Cell(20, 8, "Tahun")
	pdf.Cell(35, 8, "Nominal")
	pdf.Cell(25, 8, "Status")
	pdf.Ln(8)

	// ========== DATA TABEL ==========
	pdf.SetFont("Arial", "", 9)
	for i, iuran := range iuranList {
		pdf.Cell(15, 7, strconv.Itoa(i+1))
		pdf.Cell(55, 7, iuran.NamaWarga)
		pdf.Cell(30, 7, iuran.Bulan)
		pdf.Cell(20, 7, strconv.Itoa(iuran.Tahun))
		pdf.Cell(35, 7, "Rp "+strconv.Itoa(iuran.Nominal))

		// Status dengan warna? PDF basic gak support warna gampang, kita tulis teks aja
		if iuran.Status == "Lunas" {
			pdf.SetTextColor(0, 150, 0) // hijau
			pdf.Cell(25, 7, "Lunas")
			pdf.SetTextColor(0, 0, 0) // reset ke hitam
		} else {
			pdf.SetTextColor(200, 0, 0) // merah
			pdf.Cell(25, 7, "Belum")
			pdf.SetTextColor(0, 0, 0)
		}
		pdf.Ln(7)
	}

	// ========== FOOTER TOTAL ==========
	pdf.Ln(10)
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(40, 10, "Total Iuran Terkumpul: Rp "+strconv.Itoa(totalTerkumpul))
	pdf.Ln(7)
	pdf.Cell(40, 10, "Total Keseluruhan Iuran: Rp "+strconv.Itoa(totalKeseluruhan))
	pdf.Ln(7)
	pdf.Cell(40, 10, "Jumlah Data Iuran: "+strconv.Itoa(len(iuranList)))

	// Set header response
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=laporan_iuran.pdf")

	// Output PDF
	err = pdf.Output(w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
