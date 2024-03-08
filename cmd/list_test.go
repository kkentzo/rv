package cmd

import (
	"fmt"
	"os"
	"path"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kkentzo/rv/release"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_List_ShouldListAllReleases(t *testing.T) {
	workspacePath := uuid.NewString()
	defer os.RemoveAll(workspacePath)

	// === create the first release ===
	// create and execute release
	out, err := createRelease(workspacePath, "foo.txt", 2)
	releaseId1 := release.ReleaseFormatRe.FindString(out)
	require.NoError(t, err)
	assert.FileExists(t, path.Join(workspacePath, releaseId1, "foo.txt"))
	assert.FileExists(t, path.Join(workspacePath, release.CurrentLinkName, "foo.txt"))

	time.Sleep(10 * time.Millisecond)

	// === create the second release ===
	out, err = createRelease(workspacePath, "bar.txt", 2)
	releaseId2 := release.ReleaseFormatRe.FindString(out)
	require.NoError(t, err)
	assert.FileExists(t, path.Join(workspacePath, releaseId2, "bar.txt"))
	assert.FileExists(t, path.Join(workspacePath, release.CurrentLinkName, "bar.txt"))

	list, err := listReleases(workspacePath)
	require.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("%s\n%s\n", releaseId2, releaseId1), list)
}
