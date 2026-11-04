package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAdminCredentialsFromEnvironment(t *testing.T) {
	t.Setenv("ADMIN_USERNAME", "portfolio-admin")
	t.Setenv("ADMIN_PASSWORD", "strong-password")

	username, password, ok := adminCredentials()

	require.True(t, ok)
	require.Equal(t, "portfolio-admin", username)
	require.Equal(t, "strong-password", password)
}

func TestAdminCredentialsRequireBothValues(t *testing.T) {
	t.Setenv("ADMIN_USERNAME", "portfolio-admin")
	t.Setenv("ADMIN_PASSWORD", "")

	username, password, ok := adminCredentials()

	require.False(t, ok)
	require.Empty(t, username)
	require.Empty(t, password)
}
