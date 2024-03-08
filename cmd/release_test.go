package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_Release_ShouldCreateResources_FromScratch(t *testing.T) {
	// do not create the workspace -- just specify the path
	workspacePath := uuid.NewString()
	// ensure that we'll clean up
	defer os.RemoveAll(workspacePath)

	// create and execute release
	releaseId, err := createRelease(workspacePath, "foo.txt")
	assert.NoError(t, err)

	// the workspace should now be present
	assert.DirExists(t, workspacePath)
	// the workspace should contain the release (extracted bundle under a versioned directory)
	assert.DirExists(t, path.Join(workspacePath, releaseId))
	// the release should have the correct contents
	assert.FileExists(t, path.Join(workspacePath, releaseId, "foo.txt"))
	// the workspace should contain the "current" release link
	assert.FileExists(t, path.Join(workspacePath, "current"))
	// the "current" release link should point to the release folder
	assert.FileExists(t, path.Join(workspacePath, "current", "foo.txt"))
}

func Test_Release_ShouldCleanUp_WhenBundleDoesNotExist(t *testing.T) {
	workspacePath := uuid.NewString()
	// ensure that we'll clean up
	defer os.RemoveAll(workspacePath)

	// DO NOT create this bundle
	bundlePath := fmt.Sprintf("%s.zip", uuid.NewString())

	// prepare and execute command
	RootCmd.SetArgs([]string{"release", "-w", workspacePath, bundlePath})
	// FIRE!
	RootCmd.Execute()

	// the workspace will not be cleared up
	assert.DirExists(t, workspacePath)
	// the workspace directory should be empty
	entries, err := ioutil.ReadDir(workspacePath)
	assert.NoError(t, err)
	assert.Empty(t, entries)
	// the workspace should NOT contain the "current" release link
	assert.NoFileExists(t, path.Join(workspacePath, "current"))
}

func Test_Release_ShouldUpdateCurrent_WhenPreviousReleaseExists(t *testing.T) {
	workspacePath := uuid.NewString()
	defer os.RemoveAll(workspacePath)

	// === create the first release ===
	// create and execute release
	releaseId, err := createRelease(workspacePath, "foo.txt")
	assert.NoError(t, err)
	assert.FileExists(t, path.Join(workspacePath, releaseId, "foo.txt"))
	assert.FileExists(t, path.Join(workspacePath, "current", "foo.txt"))

	time.Sleep(10 * time.Millisecond)

	// === create the second release ===
	releaseId, err = createRelease(workspacePath, "bar.txt")
	assert.NoError(t, err)
	assert.FileExists(t, path.Join(workspacePath, releaseId, "bar.txt"))
	assert.FileExists(t, path.Join(workspacePath, "current", "bar.txt"))
}
