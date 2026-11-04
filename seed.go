package main

import (
	"log"
	"os"
	"portfolio-backend/config"
	"portfolio-backend/models"

	"golang.org/x/crypto/bcrypt"
)

func adminCredentials() (string, string, bool) {
	username := os.Getenv("ADMIN_USERNAME")
	password := os.Getenv("ADMIN_PASSWORD")
	if username == "" || password == "" {
		return "", "", false
	}

	return username, password, true
}

// defaultProjects adalah data awal. Isinya sama dengan data/profile.ts di frontend,
// ditambah field tampilan (category, accent, image_fit, sort_order).
func defaultProjects() []models.Project {
	return []models.Project{
		{
			Title:       "Sistem Pemesanan Online Pitstop Vape",
			Stack:       "PHP, MySQL, HTML, CSS",
			Description: "Website pemesanan online dan katalog produk vape dengan fitur manajemen produk, pemesanan pelanggan, serta pengelolaan data berbasis MySQL. Antarmuka responsif dan mudah digunakan.",
			Link:        "https://pitstopvape.netlify.app/",
			Image:       "/Projects/pitstop-vape.png",
			Category:    "Web Development", Accent: "blue", ImageFit: "cover", SortOrder: 1,
		},
		{
			Title:       "Desain UI Website Top Up MLBB",
			Stack:       "Figma",
			Description: "Prototype UI website top up Mobile Legends dengan fitur pembelian diamond, pilihan nominal, metode pembayaran, dan riwayat transaksi. Alur navigasi intuitif dari pemilihan produk hingga konfirmasi.",
			Link:        "https://www.figma.com/proto/fRhtF3jEL9VQ5EYerpOSFW/NOIRPEDIA?node-id=92-2961",
			Image:       "/Projects/mlbb-topup.png",
			Category:    "UI/UX Prototype", Accent: "cyan", ImageFit: "cover", SortOrder: 2,
		},
		{
			Title:       "Nexa Projects — Company Profile & Jasa Web",
			Stack:       "HTML, CSS, JavaScript",
			Description: "Membangun website company profile untuk layanan jasa pembuatan website, dilengkapi halaman portofolio dan informasi layanan.",
			Link:        "https://nexaprojects.netlify.app/",
			Image:       "/Projects/nexa.png",
			Category:    "Web Development", Accent: "violet", ImageFit: "contain", SortOrder: 3,
		},
		{
			Title:       "Desain UI Sistem Pemesanan Lapangan Olahraga",
			Stack:       "Figma",
			Description: "Prototype UI aplikasi booking lapangan olahraga, mencakup alur reservasi, pemilihan jadwal, dan halaman konfirmasi.",
			Link:        "https://www.figma.com/proto/JOpUtmbQg0ZCQw0VHWKBQt/SIBOLANG-JIR",
			Image:       "/Projects/sibolang.png",
			Category:    "UI/UX System", Accent: "orange", ImageFit: "contain", SortOrder: 4,
		},
		{
			Title:       "Simulasi Jaringan Arena Esports",
			Stack:       "GNS3",
			Description: "Merancang dan mensimulasikan topologi jaringan arena esports dengan 3 zona dinamis (gaming, streaming, broadcast) dan 1 area pengawasan CCTV.",
			Link:        "",
			Image:       "",
			Category:    "Infrastructure", Accent: "slate", ImageFit: "contain", SortOrder: 5,
		},
	}
}

// backfillProjectView mengisi field tampilan pada project LAMA (dibuat sebelum kolom
// category/sort_order ada). Hanya mengisi yang kosong, jadi hasil edit lewat admin
// tidak pernah ditimpa. Dicocokkan lewat judul.
func backfillProjectView() {
	for _, d := range defaultProjects() {
		var p models.Project
		if err := config.DB.Where("title = ?", d.Title).First(&p).Error; err != nil {
			continue
		}
		updates := map[string]interface{}{}
		// Baris "lama" = belum pernah diatur lewat Admin (kategori kosong DAN urutan 0).
		// AutoMigrate memberi baris lama accent=blue dan image_fit=cover sebagai default,
		// jadi untuk baris lama kedua nilai itu ikut diisi ulang. Baris yang sudah diubah
		// admin tidak masuk kondisi ini, sehingga pilihannya tidak pernah ditimpa.
		legacy := p.Category == "" && p.SortOrder == 0
		if p.Category == "" {
			updates["category"] = d.Category
		}
		if p.SortOrder == 0 {
			updates["sort_order"] = d.SortOrder
		}
		if legacy {
			updates["accent"] = d.Accent
			updates["image_fit"] = d.ImageFit
		}
		if len(updates) > 0 {
			config.DB.Model(&p).Updates(updates)
			log.Println("Backfill tampilan project:", p.Title)
		}
	}
}

func SeedDatabase() {
	// 1. Seed Projects
	var projectCount int64
	config.DB.Model(&models.Project{}).Count(&projectCount)

	if projectCount == 0 {
		log.Println("Seeding default projects...")
		for _, p := range defaultProjects() {
			p := p
			if err := config.DB.Create(&p).Error; err != nil {
				log.Println("Gagal seed project:", err)
			}
		}
		log.Println("Seeding projects selesai!")
	} else {
		backfillProjectView()
	}

	// 2. Seed Admin User
	var userCount int64
	config.DB.Model(&models.User{}).Count(&userCount)

	if userCount == 0 {
		username, password, ok := adminCredentials()
		if !ok {
			log.Println("Admin seed dilewati: ADMIN_USERNAME dan ADMIN_PASSWORD wajib diisi")
			return
		}

		log.Println("Seeding default admin user...")
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Println("Gagal hash password admin:", err)
			return
		}
		admin := models.User{
			Username: username,
			Password: string(hashedPassword),
		}
		if err := config.DB.Create(&admin).Error; err != nil {
			log.Println("Gagal seed admin user:", err)
		} else {
			log.Println("Default admin user created!")
		}
	}
}
