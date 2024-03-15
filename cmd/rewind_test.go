package cmd

import (
	"os"
	"path"
	"testing"

	"github.com/google/uuid"
	"github.com/kkentzo/rv/release"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Rewind_NoTarget_ShouldResetCurrentAndDeletePreviousCurrent(t *testing.T) {
	workspacePath := uuid.NewString()
	defer os.RemoveAll(workspacePath)

	releases, err := createReleases(workspacePath, 3)
	require.Len(t, releases, 3)

	// verify the current release
	current, err := release.GetCurrent(workspacePath)
	assert.NoError(t, err)
	assert.Equal(t, releases[2], current)

	// ok, let's rewind now with no target
	_, err = rewindRelease(workspacePath, "")
	require.NoError(t, err)

	// verify the remaining releases
	assert.DirExists(t, path.Join(workspacePath, releases[0]))
	assert.DirExists(t, path.Join(workspacePath, releases[1]))
	assert.NoDirExists(t, path.Join(workspacePath, releases[2]))

	// verify the current link
	current, err = release.GetCurrent(workspacePath)
	assert.NoError(t, err)
	assert.Equal(t, releases[1], current)
}

func Test_Rewind_WithTarget_ShouldResetCurrentAndDeletePreviousCurrent(t *testing.T) {
	workspacePath := uuid.NewString()
	defer os.RemoveAll(workspacePath)

	releases, err := createReleases(workspacePath, 3)
	require.Len(t, releases, 3)

	// verify the current release
	current, err := release.GetCurrent(workspacePath)
	assert.NoError(t, err)
	assert.Equal(t, releases[2], current)

	// ok, let's rewind now
	_, err = rewindRelease(workspacePath, releases[0])
	require.NoError(t, err)

	// verify the remaining releases
	assert.DirExists(t, path.Join(workspacePath, releases[0]))
	assert.NoDirExists(t, path.Join(workspacePath, releases[1]))
	assert.NoDirExists(t, path.Join(workspacePath, releases[2]))

	// verify the current link
	current, err = release.GetCurrent(workspacePath)
	assert.NoError(t, err)
	assert.Equal(t, releases[0], current)
}

func Test_Rewind_WhenTheTargetDoesNotExist(t *testing.T) {
	workspacePath := uuid.NewString()
	defer os.RemoveAll(workspacePath)

	releases, err := createReleases(workspacePath, 3)
	require.Len(t, releases, 3)

	// verify the current release
	current, err := release.GetCurrent(workspacePath)
	assert.NoError(t, err)
	assert.Equal(t, releases[2], current)

	// ok, let's rewind now with no target
	out, err := rewindRelease(workspacePath, "a_non_existent_target")
	assert.NoError(t, err)
	assert.Contains(t, out, "a_non_existent_target not found")

	assert.DirExists(t, path.Join(workspacePath, releases[0]))
	assert.DirExists(t, path.Join(workspacePath, releases[1]))
	assert.DirExists(t, path.Join(workspacePath, releases[2]))

	// verify the current link
	current, err = release.GetCurrent(workspacePath)
	assert.NoError(t, err)
	assert.Equal(t, releases[2], current)
}

func Test_Rewind_WhenThereIsOnlyOneRelease(t *testing.T) {
	workspacePath := uuid.NewString()
	defer os.RemoveAll(workspacePath)

	releases, err := createReleases(workspacePath, 1)
	require.Len(t, releases, 1)

	// verify the current link
	current, err := release.GetCurrent(workspacePath)
	assert.NoError(t, err)
	assert.Equal(t, releases[0], current)

	// ok, let's rewind now
	out, err := rewindRelease(workspacePath, "")
	assert.NoError(t, err)
	assert.Contains(t, out, "only one release in workspace")

	assert.DirExists(t, path.Join(workspacePath, releases[0]))

	// verify the current link
	current, err = release.GetCurrent(workspacePath)
	assert.NoError(t, err)
	assert.Equal(t, releases[0], current)
}

func Test_Rewind_WhenTargetIsAlreadyCurrent(t *testing.T) {
	workspacePath := uuid.NewString()
	defer os.RemoveAll(workspacePath)

	releases, err := createReleases(workspacePath, 2)
	require.Len(t, releases, 2)

	// verify the current release
	current, err := release.GetCurrent(workspacePath)
	assert.NoError(t, err)
	assert.Equal(t, releases[1], current)

	// ok, let's rewind now with no target
	out, err := rewindRelease(workspacePath, releases[1])
	assert.NoError(t, err)
	assert.Contains(t, out, "will not rewind")

	assert.DirExists(t, path.Join(workspacePath, releases[0]))
	assert.DirExists(t, path.Join(workspacePath, releases[1]))

	// verify the current link
	current, err = release.GetCurrent(workspacePath)
	assert.NoError(t, err)
	assert.Equal(t, releases[1], current)
}
