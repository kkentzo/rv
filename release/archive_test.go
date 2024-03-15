package release

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Tarball_Decompression(t *testing.T) {
	uid, gid, err := resolveUser("")
	require.NoError(t, err)

	target := uuid.NewString()
	defer os.RemoveAll(target)

	require.NoError(t, decompressArchive("test/foo.tar.gz", target, uid, gid))
	assert.FileExists(t, path.Join(target, "foo/bar.txt"))
}

func Test_Zip_Decompression(t *testing.T) {
	uid, gid, err := resolveUser("")
	require.NoError(t, err)

	target := uuid.NewString()
	defer os.RemoveAll(target)

	require.NoError(t, decompressArchive("test/foo.zip", target, uid, gid))
	assert.FileExists(t, path.Join(target, "foo/bar.txt"))
}

func Test_UnsupportedArchive(t *testing.T) {
	uid, gid, err := resolveUser("")
	require.NoError(t, err)

	source := fmt.Sprintf("%s.txt", uuid.NewString())
	require.NoError(t, os.WriteFile(source, []byte("hello"), 0777))
	defer os.Remove(source)

	assert.ErrorContains(t, decompressArchive(source, "", uid, gid), "unsupported")
}
