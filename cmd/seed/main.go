package main

import (
	"course-planner-api/config"
	"course-planner-api/internal/models"
	"log"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func main() {
	db := config.LoadDatabase()

	// Seed admin
	if err := seedUserIfNotExists(db, "admin@example.com", &models.User{
		Name:   "Admin Satu",
		Email:  "admin@example.com",
		Role:   "admin",
		NIM:    "",
		NIDN:   "",
	}); err != nil {
		log.Fatalf("failed to seed admin: %v", err)
	}

	if err := seedUserIfNotExists(db, "admin2@example.com", &models.User{
		Name:   "Admin Dua",
		Email:  "admin2@example.com",
		Role:   "admin",
		NIM:    "",
		NIDN:   "",
	}); err != nil {
		log.Fatalf("failed to seed admin kedua: %v", err)
	}

	// Seed dosen
	dosen := &models.User{
		Name:  "Dosen Satu",
		Email: "dosen@example.com",
		Role:  "dosen",
		NIM:   "",
		NIDN:  "1234567890",
	}
	if err := seedUserIfNotExists(db, dosen.Email, dosen); err != nil {
		log.Fatalf("failed to seed dosen: %v", err)
	}

	dosen2 := &models.User{
		Name:  "Dosen Dua",
		Email: "dosen2@example.com",
		Role:  "dosen",
		NIM:   "",
		NIDN:  "0987654321",
	}
	if err := seedUserIfNotExists(db, dosen2.Email, dosen2); err != nil {
		log.Fatalf("failed to seed dosen kedua: %v", err)
	}

	// Reload dosen to ensure we have its ID
	var dosenUser models.User
	if err := db.Where("email = ?", dosen.Email).First(&dosenUser).Error; err != nil {
		log.Fatalf("failed to load dosen after seed: %v", err)
	}

	// Seed mahasiswa with dosen_pa_id
	mahasiswa := &models.User{
		Name:      "Mahasiswa Satu",
		Email:     "mahasiswa@example.com",
		Role:      "mahasiswa",
		NIM:       "2024000001",
		NIDN:      "",
		DosenPAID: &dosenUser.ID,
	}
	if err := seedUserIfNotExists(db, mahasiswa.Email, mahasiswa); err != nil {
		log.Fatalf("failed to seed mahasiswa: %v", err)
	}

	mahasiswa2 := &models.User{
		Name:      "Mahasiswa Dua",
		Email:     "mahasiswa2@example.com",
		Role:      "mahasiswa",
		NIM:       "2024000002",
		NIDN:      "",
		DosenPAID: &dosenUser.ID,
	}
	if err := seedUserIfNotExists(db, mahasiswa2.Email, mahasiswa2); err != nil {
		log.Fatalf("failed to seed mahasiswa kedua: %v", err)
	}

	// Seed courses
	if err := seedCourseIfNotExists(db, "JSI60214", &models.Course{
		Kode: "JSI60214",
		Nama: "Aplikasi Berbasis Layanan",
		SKS:  3,
	}); err != nil {
		log.Fatalf("failed to seed course ABL: %v", err)
	}

	if err := seedCourseIfNotExists(db, "JSI60204", &models.Course{
		Kode: "JSI60204",
		Nama: "Tata Kelola",
		SKS:  3,
	}); err != nil {
		log.Fatalf("failed to seed course Tata Kelola: %v", err)
	}

	// Seed rooms
	if err := seedRoomIfNotExists(db, "H1.1", &models.Room{
		Nama: "H1.1",
	}); err != nil {
		log.Fatalf("failed to seed room H1.1: %v", err)
	}

	if err := seedRoomIfNotExists(db, "H1.2", &models.Room{
		Nama: "H1.2",
	}); err != nil {
		log.Fatalf("failed to seed room H1.2: %v", err)
	}

	if err := seedRoomIfNotExists(db, "H1.3", &models.Room{
		Nama: "H1.3",
	}); err != nil {
		log.Fatalf("failed to seed room H1.3: %v", err)
	}

	// Seed one class (static IDs from request)
	if err := seedClassIfNotExists(db); err != nil {
		log.Fatalf("failed to seed class: %v", err)
	}

	log.Println("Seeding users selesai.")
}

func seedUserIfNotExists(db *gorm.DB, email string, user *models.User) error {
	var existing models.User
	if err := db.Where("email = ?", email).First(&existing).Error; err == nil {
		// sudah ada, skip
		return nil
	}

	// default password untuk semua akun seed
	const defaultPassword = "password123"
	hashed, err := bcrypt.GenerateFromPassword([]byte(defaultPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashed)

	return db.Create(user).Error
}

func seedCourseIfNotExists(db *gorm.DB, kode string, course *models.Course) error {
	var existing models.Course
	if err := db.Where("kode = ?", kode).First(&existing).Error; err == nil {
		// sudah ada, skip
		return nil
	}
	return db.Create(course).Error
}

func seedRoomIfNotExists(db *gorm.DB, nama string, room *models.Room) error {
	var existing models.Room
	if err := db.Where("nama = ?", nama).First(&existing).Error; err == nil {
		// sudah ada, skip
		return nil
	}
	return db.Create(room).Error
}

func seedClassIfNotExists(db *gorm.DB) error {
	courseID := uuid.MustParse("1c1cf54c-e380-4369-a038-5d2bcf0926c0")
	dosenID := uuid.MustParse("cba38bab-3e52-4f06-9bfe-112ae81e32cf")
	roomID := uuid.MustParse("dc72adcf-c666-4753-994a-a26f6e0718d3")

	// Cek apakah sudah ada kelas dengan kombinasi ini
	var existing models.Class
	if err := db.Where("course_id = ? AND dosen_id = ? AND nama_kelas = ? AND hari = ? AND room_id = ?",
		courseID, dosenID, "B", "Senin", roomID).First(&existing).Error; err == nil {
		return nil
	}

	jamMulai, err := time.Parse("15:04", "08:00")
	if err != nil {
		return err
	}
	jamSelesai, err := time.Parse("15:04", "10:00")
	if err != nil {
		return err
	}

	class := &models.Class{
		CourseID:          courseID,
		DosenID:           dosenID,
		NamaKelas:         "B",
		Hari:              "Senin",
		JamMulai:          jamMulai,
		JamSelesai:        jamSelesai,
		RoomID:            roomID,
		Kuota:             30,
		SemesterPenawaran: "ganjil",
	}

	return db.Create(class).Error
}
