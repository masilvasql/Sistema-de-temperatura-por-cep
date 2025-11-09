package tests

import (
    "os"
    "strings"
    "testing"

    "github.com/masilvasql/sistema-de-temperatura-por-cep/pkg"
    "github.com/stretchr/testify/assert"
)

func TestIsValidZipCode(t *testing.T) {
    t.Run("valid 8 digits", func(t *testing.T) {
        assert.True(t, pkg.IsValidZipCode("12345678"))
    })

    t.Run("invalid short", func(t *testing.T) {
        assert.False(t, pkg.IsValidZipCode("1234"))
    })

    t.Run("empty", func(t *testing.T) {
        assert.False(t, pkg.IsValidZipCode(""))
    })
}

func TestGetRootPath(t *testing.T) {
    root := pkg.GetRootPath()
    assert.NotEmpty(t, root, "root path should not be empty")
    // root path should end with OS path separator
    assert.True(t, strings.HasSuffix(root, string(os.PathSeparator)), "root path should end with path separator")
}
