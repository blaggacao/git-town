package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDirectoryCloner(t *testing.T) {
	repo := CreateTestGitTownRepo(t)
	cloner, err := NewDirectoryCloner(repo.Dir)
	assert.Nil(t, err)
	assert.NotNil(t, cloner)
	newDir := filepath.Join(createTempDir(t), "foo")
	err = cloner.CreateCopy(newDir)
	assert.Nil(t, err)
	// check if some files were cloned
	info, err := os.Stat(filepath.Join(newDir, ".gitconfig"))
	assert.Nil(t, err)
	assert.False(t, info.IsDir())
}
