package main

import (
	"log"
	"net/http"
	"rt-manager/handlers"
	"rt-manager/models"

	"github.com/dgraph-io/badger/v4"
)

func main() {
	// Buka database Badger
	db, err := badger.Open(badger.DefaultOptions("./database"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Inisialisasi model
	wargaModel := models.NewWargaModel(db)

	// Seed data dummy jika kosong
	count, _ := wargaModel.Count()
	if count == 0 {
		dummyData := []models.Warga{
			{ID: 1, Nama: "Budi Santoso", NIK: "3273010101900001", Alamat: "Blok A-12", Status: "Tetap"},
			{ID: 2, Nama: "Siti Aminah", NIK: "3273010202900002", Alamat: "Blok B-05", Status: "Kontrak"},
			{ID: 3, Nama: "Dedi Prasetyo", NIK: "3273010303900003", Alamat: "Blok C-01", Status: "Tetap"},
			{ID: 4, Nama: "Rina Wati", NIK: "3273010404900004", Alamat: "Blok A-08", Status: "Tetap"},
			{ID: 5, Nama: "Agus Salim", NIK: "3273010505900005", Alamat: "Blok B-12", Status: "Kontrak"},
			{ID: 6, Nama: "Dewi Sartika", NIK: "3273010606900006", Alamat: "Blok C-07", Status: "Tetap"},
			{ID: 7, Nama: "Eko Prasetyo", NIK: "3273010707900007", Alamat: "Blok A-03", Status: "Tetap"},
			{ID: 8, Nama: "Fitri Handayani", NIK: "3273010808900008", Alamat: "Blok B-09", Status: "Kontrak"},
			{ID: 9, Nama: "Gunawan Wibowo", NIK: "3273010909900009", Alamat: "Blok C-12", Status: "Tetap"},
			{ID: 10, Nama: "Heni Marlina", NIK: "3273011010000010", Alamat: "Blok A-15", Status: "Tetap"},
		}
		for _, w := range dummyData {
			wargaModel.Insert(w)
		}
		log.Println("✅ Data dummy berhasil ditambahkan (10 warga)")
	}

	// Inisialisasi model iuran
	iuranModel := models.NewIuranModel(db)
	iuranHandler := handlers.NewIuranHandler(iuranModel, wargaModel)
	// Inisialisasi model ronda
	rondaModel := models.NewRondaModel(db)
	rondaHandler := handlers.NewRondaHandler(rondaModel, wargaModel)
	// Inisialisasi handler
	wargaHandler := handlers.NewWargaHandler(wargaModel)
	// Inisialisasi model pengumuman
	pengumumanModel := models.NewPengumumanModel(db)
	pengumumanHandler := handlers.NewPengumumanHandler(pengumumanModel)

	// Inisialisasi model event
	eventModel := models.NewEventModel(db)
	eventHandler := handlers.NewEventHandler(eventModel)

	// Inisialisasi model pengaturan
	pengaturanModel := models.NewPengaturanModel(db)
	pengaturanHandler := handlers.NewPengaturanHandler(pengaturanModel)

	// Routing
	http.HandleFunc("/", wargaHandler.Dashboard)
	http.HandleFunc("/tambah", wargaHandler.Tambah)
	http.HandleFunc("/edit/", wargaHandler.Edit)
	http.HandleFunc("/hapus/", wargaHandler.Hapus)
	http.HandleFunc("/api/warga", wargaHandler.APIWarga)

	http.HandleFunc("/iuran", iuranHandler.Index)
	http.HandleFunc("/iuran/tambah", iuranHandler.Tambah)
	http.HandleFunc("/iuran/edit/", iuranHandler.Edit)
	http.HandleFunc("/iuran/hapus/", iuranHandler.Hapus)
	http.HandleFunc("/iuran/lunas/", iuranHandler.TandaiLunas)
	http.HandleFunc("/iuran/export-pdf", iuranHandler.ExportPDF)

	http.HandleFunc("/ronda", rondaHandler.Index)
	http.HandleFunc("/ronda/tambah", rondaHandler.Tambah)
	http.HandleFunc("/ronda/edit/", rondaHandler.Edit)
	http.HandleFunc("/ronda/hapus/", rondaHandler.Hapus)
	http.HandleFunc("/ronda/selesai/", rondaHandler.TandaiSelesai)

	http.HandleFunc("/pengumuman", pengumumanHandler.Index)
	http.HandleFunc("/pengumuman/tambah", pengumumanHandler.Tambah)
	http.HandleFunc("/pengumuman/edit/", pengumumanHandler.Edit)
	http.HandleFunc("/pengumuman/hapus/", pengumumanHandler.Hapus)
	http.HandleFunc("/pengumuman/arsip/", pengumumanHandler.Arsipkan)

	http.HandleFunc("/warga/export-pdf", wargaHandler.ExportPDF)

	http.HandleFunc("/event", eventHandler.Index)
	http.HandleFunc("/event/tambah", eventHandler.Tambah)
	http.HandleFunc("/event/edit/", eventHandler.Edit)
	http.HandleFunc("/event/hapus/", eventHandler.Hapus)

	// ========== ROUTING PENGATURAN ==========
	http.HandleFunc("/pengaturan", pengaturanHandler.Index)
	http.HandleFunc("/pengaturan/edit", pengaturanHandler.Edit)
	// ========================================

	log.Println("🚀 Server berjalan di http://localhost:8080")
	log.Println("📁 Database: Badger (folder ./database)")
	log.Println("✨ Tekan Ctrl+C untuk berhenti")
	http.ListenAndServe(":8080", nil)

}
