package models

import (
	"strings"

	"gorm.io/gorm"
)

// Tool adalah satu logo teknologi di section "Tools & Capabilities".
type Tool struct {
	gorm.Model
	Group     string `json:"group" gorm:"column:tool_group;default:frontend"` // frontend|backend|design|dev
	Name      string `json:"name"`
	Icon      string `json:"icon"` // path/URL logo (file di public/tools/ frontend)
	SortOrder int    `json:"sort_order"`
}

// "group" adalah kata kunci SQL, jadi kolomnya diberi nama tool_group.
var ValidToolGroups = map[string]bool{"frontend": true, "backend": true, "design": true, "dev": true}

// Normalize merapikan isian dan mengganti kelompok yang tidak dikenal dengan "frontend".
func (t *Tool) Normalize() {
	t.Name = strings.TrimSpace(t.Name)
	t.Icon = strings.TrimSpace(t.Icon)
	t.Group = strings.ToLower(strings.TrimSpace(t.Group))
	if !ValidToolGroups[t.Group] {
		t.Group = "frontend"
	}
}

// Problem mengembalikan pesan kesalahan bila isian belum valid, atau "" bila sudah valid.
func (t *Tool) Problem() string {
	if t.Name == "" {
		return "Name tidak boleh kosong"
	}
	return ""
}

// Base memberi akses ke ID/created_at supaya controller bisa melindunginya saat update.
func (t *Tool) Base() *gorm.Model { return &t.Model }
