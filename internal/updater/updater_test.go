//go:build unit

package updater

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"github.com/hpcsc/kz/internal/gateway"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"io/fs"
	"os"
	"path"
	"strings"
	"testing"
	"time"
)

func TestGithubReleaseUpdater_Update(t *testing.T) {
	t.Run("return error when github latest release endpoint returns error", func(t *testing.T) {
		gw := gateway.NewFakeGithubGateway()
		gw.OnLatestRelease().Return(nil, fmt.Errorf("some error"))
		u := New("amd64", "/path/to/executable", gw)

		_, err := u.UpdateFrom("v1.0.0")

		require.Error(t, err)
		require.Contains(t, err.Error(), "some error")
	})

	t.Run("return message with no error when current version is already latest", func(t *testing.T) {
		gw := gateway.NewFakeGithubGateway()
		gw.OnLatestRelease().Return(&gateway.LatestReleaseResponse{
			TagName: "v1.0.0",
		}, nil)
		u := New("amd64", "/path/to/executable", gw)

		m, err := u.UpdateFrom("v1.0.0")

		require.NoError(t, err)
		require.Equal(t, "kz is already at latest version v1.0.0", m)
	})

	t.Run("return error when github latest endpoint returns no asset matching current cpu architecture", func(t *testing.T) {
		gw := gateway.NewFakeGithubGateway()
		gw.OnLatestRelease().Return(&gateway.LatestReleaseResponse{
			TagName: "v1.1.0",
			Assets: []gateway.LatestReleaseAsset{
				{
					Url: "http://github-gw/releases/download/v1.1.0/kz-darwin-arm64.tar.gz",
				},
			},
		}, nil)
		u := New("amd64", "/path/to/executable", gw)

		_, err := u.UpdateFrom("v1.0.0")

		require.Error(t, err)
		require.Contains(t, err.Error(), "no artifact available for release v1.1.0 and current CPU architecture amd64")
	})

	t.Run("return error when github download url returns error", func(t *testing.T) {
		gw := gateway.NewFakeGithubGateway()
		gw.OnLatestRelease().Return(&gateway.LatestReleaseResponse{
			TagName: "v1.1.0",
			Assets: []gateway.LatestReleaseAsset{
				{
					Url: "http://github-gw/releases/download/v1.1.0/kz-darwin-amd64.tar.gz",
				},
			},
		}, nil)
		gw.OnDownload(mock.Anything).Return(nil, fmt.Errorf("some error"))
		u := New("amd64", "/path/to/executable", gw)

		_, err := u.UpdateFrom("v1.0.0")

		require.Error(t, err)
		require.Contains(t, err.Error(), "some error")
	})

	t.Run("return error when downloaded archive does not contain correct binary", func(t *testing.T) {
		gw := gateway.NewFakeGithubGateway()
		archiveName := "kz-darwin-amd64.tar.gz"
		gw.OnLatestRelease().Return(&gateway.LatestReleaseResponse{
			TagName: "v1.1.0",
			Assets: []gateway.LatestReleaseAsset{
				{
					Url: "http://github-gw/releases/download/v1.1.0/" + archiveName,
				},
			},
		}, nil)

		archivePath := path.Join(os.TempDir(), archiveName)
		invalidGzippedArtifact(t, archivePath)
		archive, err := os.ReadFile(archivePath)
		require.NoError(t, err)
		defer os.Remove(archivePath)

		gw.OnDownload(mock.Anything).Return(archive, nil)

		// instead of updating current executable, create a fake executable in tmp dir and replace that instead
		currentExecutable := fakeCurrentExecutable(t)
		defer os.Remove(currentExecutable)

		u := New("amd64", currentExecutable, gw)

		_, err = u.UpdateFrom("v1.0.0")

		require.Error(t, err)
		require.Contains(t, err.Error(), "no kz binary found")
	})

	t.Run("replace current executable with latest version from github", func(t *testing.T) {
		gw := gateway.NewFakeGithubGateway()
		archiveName := "kz-darwin-amd64.tar.gz"
		gw.OnLatestRelease().Return(&gateway.LatestReleaseResponse{
			TagName: "v1.1.0",
			Assets: []gateway.LatestReleaseAsset{
				{
					Url: "http://github-gw/releases/download/v1.1.0/" + archiveName,
				},
			},
		}, nil)

		archivePath := path.Join(os.TempDir(), archiveName)
		newTestGzippedArtifact(t, archivePath)
		archive, err := os.ReadFile(archivePath)
		require.NoError(t, err)
		defer os.Remove(archivePath)

		gw.OnDownload(mock.Anything).Return(archive, nil)

		// instead of updating current executable, create a fake executable in tmp dir and replace that instead
		currentExecutable := fakeCurrentExecutable(t)
		defer os.Remove(currentExecutable)

		u := New("amd64", currentExecutable, gw)

		m, err := u.UpdateFrom("v1.0.0")

		require.NoError(t, err)
		require.Equal(t, "updated to v1.1.0", m)

		stat, err := os.Stat(currentExecutable)
		require.NoError(t, err)
		require.Equal(t, fs.FileMode(0o755), stat.Mode().Perm())

		executableContent, err := os.ReadFile(currentExecutable)
		require.NoError(t, err)
		require.Equal(t, "this is newer kz", string(executableContent))
	})
}

func fakeCurrentExecutable(t *testing.T) string {
	currentExecutable, err := os.CreateTemp(os.TempDir(), "kz-*")
	require.NoError(t, err)
	defer currentExecutable.Close()
	_, err = currentExecutable.WriteString("this is old executable")
	require.NoError(t, err)
	return currentExecutable.Name()
}

func invalidGzippedArtifact(t *testing.T, path string) {
	archive, err := os.Create(path)
	require.NoError(t, err)

	defer archive.Close()

	gw := gzip.NewWriter(archive)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()

	file1 := strings.NewReader("file 1 content")
	addFile(t, tw, "file-1", file1.Size(), file1)
	file2 := strings.NewReader("file 2 content")
	addFile(t, tw, "file-2", file2.Size(), file2)
}

func newTestGzippedArtifact(t *testing.T, path string) {
	archive, err := os.Create(path)
	require.NoError(t, err)

	defer archive.Close()

	gw := gzip.NewWriter(archive)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()

	fakeExecutable := strings.NewReader("this is newer kz")
	addFile(t, tw, "kz", fakeExecutable.Size(), fakeExecutable)
	readMe := strings.NewReader("this is README")
	addFile(t, tw, "README.md", readMe.Size(), readMe)
}

func addFile(t *testing.T, tw *tar.Writer, path string, size int64, file io.Reader) {
	header := tar.Header{
		Name:    path,
		Size:    size,
		Mode:    int64(0o644),
		ModTime: time.Now(),
	}
	err := tw.WriteHeader(&header)
	require.NoError(t, err)

	_, err = io.Copy(tw, file)
	require.NoError(t, err)
}
