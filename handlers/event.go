package handlers

import (
	"net/http"
	"rt-manager/models"
	"strconv"
	"strings"
	"text/template"
	"time"
)

type EventHandler struct {
	Model *models.EventModel
}

func NewEventHandler(model *models.EventModel) *EventHandler {
	return &EventHandler{
		Model: model,
	}
}

// Index menampilkan daftar event
func (h *EventHandler) Index(w http.ResponseWriter, r *http.Request) {
	eventList, err := h.Model.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Balik urutan (yang terbaru di atas)
	for i, j := 0, len(eventList)-1; i < j; i, j = i+1, j-1 {
		eventList[i], eventList[j] = eventList[j], eventList[i]
	}

	funcMap := template.FuncMap{
		"add": func(a, b int) int { return a + b },
		"formatRupiah": func(nominal int) string {
			if nominal == 0 {
				return "Rp 0"
			}
			return "Rp " + strconv.Itoa(nominal)
		},
		"formatTanggal": func(tanggal string) string {
			t, _ := time.Parse("2006-01-02", tanggal)
			return t.Format("02 Jan 2006")
		},
		"getStatusBadge": func(status string) string {
			switch status {
			case "Akan Datang":
				return "bg-blue-500/20 text-blue-400"
			case "Berlangsung":
				return "bg-green-500/20 text-green-400"
			case "Selesai":
				return "bg-gray-500/20 text-gray-400"
			default:
				return "bg-gray-500/20 text-gray-400"
			}
		},
	}

	data := map[string]interface{}{
		"EventList":  eventList,
		"ActivePage": "event",
	}

	tmpl := template.Must(template.New("layout.html").Funcs(funcMap).ParseFiles("templates/layout.html", "templates/event.html"))
	tmpl.Execute(w, data)
}

// Tambah menampilkan form tambah event
func (h *EventHandler) Tambah(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		namaEvent := r.FormValue("nama_event")
		tanggal := r.FormValue("tanggal")
		lokasi := r.FormValue("lokasi")
		deskripsi := r.FormValue("deskripsi")
		anggaran, _ := strconv.Atoi(r.FormValue("anggaran"))

		eventBaru := models.Event{
			NamaEvent: namaEvent,
			Tanggal:   tanggal,
			Lokasi:    lokasi,
			Deskripsi: deskripsi,
			Anggaran:  anggaran,
		}

		err := h.Model.Insert(eventBaru)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/event", http.StatusSeeOther)
		return
	}

	funcMap := template.FuncMap{
		"now": func() string {
			return time.Now().Format("2006-01-02")
		},
	}

	data := map[string]interface{}{
		"ActivePage": "event",
	}

	tmpl := template.Must(template.New("layout.html").Funcs(funcMap).ParseFiles("templates/layout.html", "templates/event_tambah.html"))
	tmpl.Execute(w, data)
}

// Edit mengedit event
func (h *EventHandler) Edit(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	idStr := pathParts[len(pathParts)-1]
	id, _ := strconv.Atoi(idStr)

	if r.Method == "POST" {
		namaEvent := r.FormValue("nama_event")
		tanggal := r.FormValue("tanggal")
		lokasi := r.FormValue("lokasi")
		deskripsi := r.FormValue("deskripsi")
		anggaran, _ := strconv.Atoi(r.FormValue("anggaran"))

		event, _ := h.Model.GetByID(id)
		event.NamaEvent = namaEvent
		event.Tanggal = tanggal
		event.Lokasi = lokasi
		event.Deskripsi = deskripsi
		event.Anggaran = anggaran

		h.Model.Update(*event)
		http.Redirect(w, r, "/event", http.StatusSeeOther)
		return
	}

	event, err := h.Model.GetByID(id)
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
		"Event":      event,
		"ActivePage": "event",
	}

	tmpl := template.Must(template.New("layout.html").Funcs(funcMap).ParseFiles("templates/layout.html", "templates/event_edit.html"))
	tmpl.Execute(w, data)
}

// Hapus menghapus event
func (h *EventHandler) Hapus(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	idStr := pathParts[len(pathParts)-1]
	id, _ := strconv.Atoi(idStr)

	err := h.Model.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/event", http.StatusSeeOther)
}
