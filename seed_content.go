package main

import (
	_ "embed"
	"encoding/json"
	"log"

	"portfolio-backend/config"
	"portfolio-backend/models"
)

// seed_content.json dibuat dari data frontend (jalankan: node --no-warnings scripts/gen-seed.mjs),
// jadi data awal di database sama persis dengan data cadangan di website.
//
//go:embed seed_content.json
var seedContentJSON []byte

type seedContent struct {
	Experiences []models.Experience `json:"experiences"`
	Tools       []models.Tool       `json:"tools"`
}

// seedContentData menguraikan seed_content.json. Kembalian kosong bila file rusak.
func seedContentData() seedContent {
	var content seedContent
	if err := json.Unmarshal(seedContentJSON, &content); err != nil {
		log.Println("seed_content.json tidak valid:", err)
	}
	return content
}

// SeedContent mengisi tabel Pengalaman dan Tools SEKALI, hanya bila tabelnya belum pernah
// berisi apa pun. Dihitung dengan Unscoped (termasuk baris yang sudah dihapus), supaya
// data yang sengaja dihapus admin tidak muncul lagi setiap backend dinyalakan ulang.
func SeedContent() {
	content := seedContentData()

	var experienceCount int64
	config.DB.Unscoped().Model(&models.Experience{}).Count(&experienceCount)
	if experienceCount == 0 {
		log.Println("Seeding pengalaman...")
		for _, e := range content.Experiences {
			e := e
			e.Normalize()
			if err := config.DB.Create(&e).Error; err != nil {
				log.Println("Gagal seed pengalaman:", err)
			}
		}
	}

	var toolCount int64
	config.DB.Unscoped().Model(&models.Tool{}).Count(&toolCount)
	if toolCount == 0 {
		log.Println("Seeding tools...")
		for _, t := range content.Tools {
			t := t
			t.Normalize()
			if err := config.DB.Create(&t).Error; err != nil {
				log.Println("Gagal seed tool:", err)
			}
		}
	}
}