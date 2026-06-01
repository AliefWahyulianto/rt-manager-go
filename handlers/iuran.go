package handlers

import (
	"net/http"
	"rt-manager/models"
	"strconv"
	"strings"
	"text/template"
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
