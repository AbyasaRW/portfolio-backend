package models_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"portfolio-backend/models"
)

func TestExperienceNormalize_FiltersCategoriesAndDedupes(t *testing.T) {
	e := models.Experience{Company: " MNCTV ", Categories: " Broadcast , tech,BROADCAST, ngawur ,, operations ", Tone: " CYAN "}
	e.Normalize()

	require.Equal(t, "MNCTV", e.Company)
	require.Equal(t, "broadcast,tech,operations", e.Categories)
	require.Equal(t, "cyan", e.Tone)
}

func TestExperienceNormalize_DefaultsAndLists(t *testing.T) {
	e := models.Experience{Company: "x", Tone: "merah", Skills: " Go ,, Fiber , ", Points: "satu\r\n\r\n  dua  \n"}
	e.Normalize()

	require.Equal(t, "blue", e.Tone)
	require.Equal(t, "Go, Fiber", e.Skills)
	require.Equal(t, "satu\ndua", e.Points)
	require.Equal(t, "", e.Categories)
}

func TestExperienceProblem(t *testing.T) {
	require.NotEmpty(t, (&models.Experience{}).Problem())
	require.Empty(t, (&models.Experience{Company: "ok"}).Problem())
}

func TestToolNormalize(t *testing.T) {
	tool := models.Tool{Name: "  Next.js ", Icon: " /tools/nextjs.svg ", Group: " Backend "}
	tool.Normalize()
	require.Equal(t, "Next.js", tool.Name)
	require.Equal(t, "/tools/nextjs.svg", tool.Icon)
	require.Equal(t, "backend", tool.Group)

	bad := models.Tool{Name: "x", Group: "lainnya"}
	bad.Normalize()
	require.Equal(t, "frontend", bad.Group)
}

func TestToolProblem(t *testing.T) {
	// Controller selalu memanggil Normalize() sebelum Problem(); tes mengikuti urutan itu.
	blank := models.Tool{Name: "  "}
	blank.Normalize()
	require.NotEmpty(t, blank.Problem())

	ok := models.Tool{Name: " Go "}
	ok.Normalize()
	require.Empty(t, ok.Problem())
}