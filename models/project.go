package models

import (
	"strings"

	"gorm.io/gorm"
)

// Project adalah satu karya yang ditampilkan di carousel "Featured Projects".
// Nama kolom JSON sengaja sama dengan data/profile.ts di frontend.
type Project struct {
	gorm.Model
	Title       string `json:"title"`
	Stack       string `json:"stack"`
	Description string `json:"description"`
	Link        string `json:"link"`
	Image       string `json:"image"`

	// Tampilan di kartu frontend.
	Category  string `json:"category"`                       // label kategori, mis. "Web Development"
	Accent    string `json:"accent" gorm:"default:blue"`     // warna tag: blue|cyan|violet|orange|slate
	ImageFit  string `json:"image_fit" gorm:"default:cover"` // cover = penuh, contain = utuh
	SortOrder int    `json:"sort_order"`                     // makin kecil makin awal
}

// Pilihan yang valid untuk Accent dan ImageFit.
var (
	ValidAccents   = map[string]bool{"blue": true, "cyan": true, "violet": true, "orange": true, "slate": true}
	ValidImageFits = map[string]bool{"cover": true, "contain": true}
)

// Normalize merapikan isian: spasi dibuang dan nilai tidak valid diganti default.
func (p *Project) Normalize() {
	p.Title = strings.TrimSpace(p.Title)
	p.Stack = strings.TrimSpace(p.Stack)
	p.Description = strings.TrimSpace(p.Description)
	p.Link = strings.TrimSpace(p.Link)
	p.Image = strings.TrimSpace(p.Image)
	p.Category = strings.TrimSpace(p.Category)

	p.Accent = strings.ToLower(strings.TrimSpace(p.Accent))
	if !ValidAccents[p.Accent] {
		p.Accent = "blue"
	}
	p.ImageFit = strings.ToLower(strings.TrimSpace(p.ImageFit))
	if !ValidImageFits[p.ImageFit] {
		p.ImageFit = "cover"
	}
}
