package models

import (
	"strings"

	"gorm.io/gorm"
)

// Experience adalah satu kartu di section "Career Experience".
// Nama kolom JSON sengaja mengikuti data frontend (data/experience.ts + profile.ts).
type Experience struct {
	gorm.Model
	Company    string `json:"company"`
	Role       string `json:"role"`                         // judul di sisi depan kartu
	FullRole   string `json:"full_role"`                    // judul lengkap di sisi belakang kartu
	Place      string `json:"place"`                        // lokasi
	Period     string `json:"period"`                       // mis. "Agustus - November 2023"
	Categories string `json:"categories"`                   // filter, pisahkan koma: broadcast,tech,operations
	Tone       string `json:"tone" gorm:"default:blue"`     // warna inisial: blue|cyan|orange|violet
	Initials   string `json:"initials"`                     // cadangan bila logo tidak ada
	Logo       string `json:"logo"`                         // path/URL logo
	Summary    string `json:"summary"`                      // ringkasan depan; **tebal** didukung
	Skills     string `json:"skills"`                       // pisahkan koma
	Points     string `json:"points"`                       // tugas, satu per baris
	SortOrder  int    `json:"sort_order"`                   // makin kecil makin awal
}

// Pilihan yang valid.
var (
	ValidTones      = map[string]bool{"blue": true, "cyan": true, "orange": true, "violet": true}
	ValidCategories = map[string]bool{"broadcast": true, "tech": true, "operations": true}
)

// Normalize merapikan isian: spasi dibuang, kategori disaring, nilai tidak valid diganti default.
func (e *Experience) Normalize() {
	e.Company = strings.TrimSpace(e.Company)
	e.Role = strings.TrimSpace(e.Role)
	e.FullRole = strings.TrimSpace(e.FullRole)
	e.Place = strings.TrimSpace(e.Place)
	e.Period = strings.TrimSpace(e.Period)
	e.Initials = strings.TrimSpace(e.Initials)
	e.Logo = strings.TrimSpace(e.Logo)
	e.Summary = strings.TrimSpace(e.Summary)

	e.Tone = strings.ToLower(strings.TrimSpace(e.Tone))
	if !ValidTones[e.Tone] {
		e.Tone = "blue"
	}

	// Hanya kategori yang dikenal, tanpa duplikat, urutan input dipertahankan.
	seen := map[string]bool{}
	var cats []string
	for _, c := range strings.Split(e.Categories, ",") {
		c = strings.ToLower(strings.TrimSpace(c))
		if ValidCategories[c] && !seen[c] {
			seen[c] = true
			cats = append(cats, c)
		}
	}
	e.Categories = strings.Join(cats, ",")

	e.Skills = joinClean(e.Skills, ",", ", ")
	e.Points = joinClean(e.Points, "\n", "\n")
}

// joinClean memecah teks, membuang item kosong dan spasi, lalu menyambung lagi.
func joinClean(value, sep, joiner string) string {
	var out []string
	for _, part := range strings.Split(strings.ReplaceAll(value, "\r\n", "\n"), sep) {
		if part = strings.TrimSpace(part); part != "" {
			out = append(out, part)
		}
	}
	return strings.Join(out, joiner)
}

// Problem mengembalikan pesan kesalahan bila isian belum valid, atau "" bila sudah valid.
func (e *Experience) Problem() string {
	if e.Company == "" {
		return "Company tidak boleh kosong"
	}
	return ""
}

// Base memberi akses ke ID/created_at supaya controller bisa melindunginya saat update.
func (e *Experience) Base() *gorm.Model { return &e.Model }
