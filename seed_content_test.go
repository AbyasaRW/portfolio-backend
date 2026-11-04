package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// Data seed harus lengkap dan valid. Kalau frontend/data berubah, jalankan ulang
// `node --no-warnings scripts/gen-seed.mjs` supaya tes ini tetap sesuai.
func TestSeedContent_ExperiencesComplete(t *testing.T) {
	content := seedContentData()

	require.Len(t, content.Experiences, 4)
	for i, e := range content.Experiences {
		require.NotEmpty(t, e.Company, "experience %d", i)
		require.NotEmpty(t, e.Role, "experience %d", i)
		require.NotEmpty(t, e.FullRole, "experience %d", i)
		require.NotEmpty(t, e.Period, "experience %d", i)
		require.NotEmpty(t, e.Points, "experience %d", i)
		require.NotEmpty(t, e.Categories, "experience %d", i)
		require.Equal(t, i+1, e.SortOrder)
		// sudah dalam bentuk ternormalisasi
		before := e
		e.Normalize()
		require.Equal(t, before.Categories, e.Categories)
		require.Equal(t, before.Tone, e.Tone)
	}
	require.Equal(t, "MNCTV", content.Experiences[0].Company)
	require.True(t, strings.Contains(content.Experiences[1].Categories, "tech"))
}

func TestSeedContent_ToolsGroupedAndValid(t *testing.T) {
	content := seedContentData()

	require.Len(t, content.Tools, 18)
	groups := map[string]int{}
	for _, tool := range content.Tools {
		require.NotEmpty(t, tool.Name)
		require.True(t, strings.HasPrefix(tool.Icon, "/tools/"), tool.Name)
		groups[tool.Group]++
	}
	require.Equal(t, map[string]int{"frontend": 5, "backend": 4, "design": 4, "dev": 5}, groups)
}