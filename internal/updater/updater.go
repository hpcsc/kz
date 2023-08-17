package updater

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/hpcsc/kz/internal/gateway"
	"io"
	"os"
	"path"
	"strings"
)

type Updater interface {
	UpdateFrom(currentVersion string) (string, error)
}

var _ Updater = (*githubReleaseUpdater)(nil)

type githubReleaseUpdater struct {
	currentArch       string
	currentExecutable string
	githubGateway     gateway.Github
}

func New(currentArch string, currentExecutable string, githubGateway gateway.Github) Updater {
	return &githubReleaseUpdater{
		currentArch:       currentArch,
		currentExecutable: currentExecutable,
		githubGateway:     githubGateway,
	}
}

func (u *githubReleaseUpdater) UpdateFrom(currentVersion string) (string, error) {
	latestResponse, err := u.githubGateway.LatestRelease()
	if err != nil {
		return "", err
	}

	if strings.TrimPrefix(latestResponse.TagName, "v") == strings.TrimPrefix(currentVersion, "v") {
		return fmt.Sprintf("kz is already at latest version %s", latestResponse.TagName), nil
	}

	var matchingAsset *gateway.LatestReleaseAsset
	archivedSuffix := fmt.Sprintf("%s.tar.gz", u.currentArch)
	for i, a := range latestResponse.Assets {
		if strings.Contains(a.Url, archivedSuffix) {
			matchingAsset = &latestResponse.Assets[i]
			break
		}
	}

	if matchingAsset == nil {
		return "", fmt.Errorf("no artifact available for release %s and current CPU architecture %s", latestResponse.TagName, u.currentArch)
	}

	downloadResponse, err := u.githubGateway.Download(matchingAsset.Url)
	if err != nil {
		return "", err
	}

	gzipReader, err := gzip.NewReader(bytes.NewReader(downloadResponse))
	if err != nil {
		return "", fmt.Errorf("failed to create new gzip reader: %v", err)
	}

	defer gzipReader.Close()
	tarReader := tar.NewReader(gzipReader)

	for {
		hdr, err := tarReader.Next()
		if err != nil {
			if err == io.EOF {
				return "", fmt.Errorf("invalid archive from %s, no kz binary found", matchingAsset.Url)
			}

			return "", fmt.Errorf("failed to read archive file: %v", err)
		}

		if hdr.Name == "kz" {
			originalPath := path.Dir(u.currentExecutable)

			// copy content of latest binary to a temporary file, chmod and replace current executable with that temporary file
			latestBinary, err := os.CreateTemp(originalPath, fmt.Sprintf("kz-%s", latestResponse.TagName))
			if err != nil {
				return "", fmt.Errorf("failed to create temporary file: %v", err)
			}

			defer latestBinary.Close()

			_, err = io.Copy(latestBinary, tarReader)
			if err != nil {
				return "", fmt.Errorf("failed to copy file content from archive to temporary file: %v", err)
			}

			if err = latestBinary.Chmod(0o755); err != nil {
				return "", fmt.Errorf("failed to change temporary file mode: %v", err)
			}

			if err = os.Rename(latestBinary.Name(), u.currentExecutable); err != nil {
				return "", fmt.Errorf("failed to rename %s to %s: %v", latestBinary.Name(), u.currentExecutable, err)
			}

			break
		}
	}

	return fmt.Sprintf("updated to %s", latestResponse.TagName), nil
}
