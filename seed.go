package main

import (
	"log"
	"portfolio-backend/config"
	"portfolio-backend/models"

	"golang.org/x/crypto/bcrypt"
)

func SeedDatabase() {
	// 1. Seed Projects
	var projectCount int64
	config.DB.Model(&models.Project{}).Count(&projectCount)

	if projectCount == 0 {
		log.Println("Seeding default projects...")
		defaultProjects := []models.Project{
			{
				Title:       "Sistem Pemesanan Online Pitstop Vape",
				Stack:       "PHP, MySQL, HTML, CSS",
				Description: "Website pemesanan online dan katalog produk vape dengan fitur manajemen produk, pemesanan pelanggan, serta pengelolaan data berbasis MySQL. Antarmuka responsif dan mudah digunakan.",
				Link:        "https://pitstopvape.netlify.app/",
				Image:       "/Projects/pitstop-vape.png",
			},
			{
				Title:       "Desain UI Website Top Up MLBB",
				Stack:       "Figma",
				Description: "Prototype UI website top up Mobile Legends dengan fitur pembelian diamond, pilihan nominal, metode pembayaran, dan riwayat transaksi. Alur navigasi intuitif dari pemilihan produk hingga konfirmasi.",
				Link:        "https://www.figma.com/proto/fRhtF3jEL9VQ5EYerpOSFW/NOIRPEDIA?node-id=92-2961",
				Image:       "/Projects/mlbb-topup.png",
			},
			{
				Title:       "Nexa Projects — Company Profile & Jasa Web",
				Stack:       "HTML, CSS, JavaScript",
				Description: "Membangun website company profile untuk layanan jasa pembuatan website, dilengkapi halaman portofolio dan informasi layanan.",
				Link:        "https://nexaprojects.netlify.app/",
				Image:       "/Projects/nexa.png",
			},
			{
				Title:       "Desain UI Sistem Pemesanan Lapangan Olahraga",
				Stack:       "Figma",
				Description: "Prototype UI aplikasi booking lapangan olahraga, mencakup alur reservasi, pemilihan jadwal, dan halaman konfirmasi.",
				Link:        "https://www.figma.com/proto/JOpUtmbQg0ZCQw0VHWKBQt/SIBOLANG-JIR",
				Image:       "/Projects/sibolang.png",
			},
			{
				Title:       "Simulasi Jaringan Arena Esports",
				Stack:       "GNS3",
				Description: "Merancang dan mensimulasikan topologi jaringan arena esports dengan 3 zona dinamis (gaming, streaming, broadcast) dan 1 area pengawasan CCTV.",
				Link:        "",
				Image:       "",
			},
		}

		for _, p := range defaultProjects {
			if err := config.DB.Create(&p).Error; err != nil {
				log.Println("Gagal seed project:", err)
			}
		}
		log.Println("Seeding projects selesai!")
	}

	// 2. Seed Admin User
	var userCount int64
	config.DB.Model(&models.User{}).Count(&userCount)

	if userCount == 0 {
		log.Println("Seeding default admin user...")
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		if err != nil {
			log.Println("Gagal hash password admin:", err)
			return
		}
		admin := models.User{
			Username: "admin",
			Password: string(hashedPassword),
		}
		if err := config.DB.Create(&admin).Error; err != nil {
			log.Println("Gagal seed admin user:", err)
		} else {
			log.Println("Default admin user created! (Username: admin, Password: admin123)")
		}
	}
}
