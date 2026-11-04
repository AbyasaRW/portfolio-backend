package models_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"portfolio-backend/models"
)

func TestNormalizeTrimsAndKeepsValidValues(t *testing.T) {
	p := models.Project{
		Title:    "  Judul  ",
		Stack:    " Go, Fiber ",
		Category: " Web Development ",
		Accent:   "  Violet ",
		ImageFit: " CONTAIN ",
	}

	p.Normalize()

	require.Equal(t, "Judul", p.Title)
	require.Equal(t, "Go, Fiber", p.Stack)
	require.Equal(t, "Web Development", p.Category)
	require.Equal(t, "violet", p.Accent)
	require.Equal(t, "contain", p.ImageFit)
}

func TestNormalizeFallsBackToDefaults(t *testing.T) {
	p := models.Project{Title: "x", Accent: "hotpink", ImageFit: "stretch"}

	p.Normalize()

	require.Equal(t, "blue", p.Accent)
	require.Equal(t, "cover", p.ImageFit)
}

func TestNormalizeEmptyValuesGetDefaults(t *testing.T) {
	p := models.Project{Title: "x"}

	p.Normalize()

	require.Equal(t, "blue", p.Accent)
	require.Equal(t, "cover", p.ImageFit)
}
