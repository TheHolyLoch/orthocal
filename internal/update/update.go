// Orthocal - Developed by dgm (dgm@tuta.com)
// orthocal/internal/update/update.go

package update

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

type UpdateResult struct {
	Source        string
	TargetPath    string
	TempPath      string
	BackupPath    string
	BackupCreated bool
	BytesWritten  int64
}

func IsHTTPSource(source string) bool {
	parsed, err := url.Parse(source)
	if err != nil {
		return false
	}

	return parsed.Scheme == "http" || parsed.Scheme == "https"
}

func UpdateDatabase(targetPath string, source string) (UpdateResult, error) {
	if source == "" {
		return UpdateResult{}, errors.New("update source is required")
	}

	if err := reject_unsupported_scheme(source); err != nil {
		return UpdateResult{}, err
	}

	if !IsHTTPSource(source) {
		same, err := same_path(source, targetPath)
		if err != nil {
			return UpdateResult{}, err
		}

		if same {
			return UpdateResult{}, errors.New("source path must not match target path")
		}
	}

	targetDir := filepath.Dir(targetPath)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return UpdateResult{}, err
	}

	tempFile, err := os.CreateTemp(targetDir, ".orthocal-*.db")
	if err != nil {
		return UpdateResult{}, err
	}

	tempPath := tempFile.Name()
	tempFile.Close()

	result := UpdateResult{
		Source:     source,
		TargetPath: targetPath,
		TempPath:   tempPath,
		BackupPath: targetPath + ".bak",
	}

	defer func() {
		if result.TempPath != "" {
			os.Remove(result.TempPath)
		}
	}()

	if IsHTTPSource(source) {
		result.BytesWritten, err = downloadHTTPSource(source, tempPath)
	} else {
		result.BytesWritten, err = copyLocalSource(source, tempPath)
	}
	if err != nil {
		return UpdateResult{}, err
	}

	if result.BytesWritten == 0 {
		return UpdateResult{}, errors.New("source produced an empty database file")
	}

	if err := ValidateDatabase(tempPath); err != nil {
		return UpdateResult{}, err
	}

	if _, err := os.Stat(targetPath); err == nil {
		if err := os.Remove(result.BackupPath); err != nil && !os.IsNotExist(err) {
			return UpdateResult{}, err
		}

		if err := os.Rename(targetPath, result.BackupPath); err != nil {
			return UpdateResult{}, err
		}

		result.BackupCreated = true
	} else if err != nil && !os.IsNotExist(err) {
		return UpdateResult{}, err
	}

	if err := os.Rename(tempPath, targetPath); err != nil {
		return UpdateResult{}, err
	}

	result.TempPath = ""
	return result, nil
}

func copyLocalSource(source string, dest string) (int64, error) {
	sourceInfo, err := os.Stat(source)
	if err != nil {
		return 0, err
	}

	if sourceInfo.IsDir() {
		return 0, errors.New("source path is a directory")
	}

	if sourceInfo.Size() == 0 {
		return 0, errors.New("source file is empty")
	}

	input, err := os.Open(source)
	if err != nil {
		return 0, err
	}
	defer input.Close()

	output, err := os.Create(dest)
	if err != nil {
		return 0, err
	}
	defer output.Close()

	written, err := io.Copy(output, input)
	if err != nil {
		return written, err
	}

	if err := output.Sync(); err != nil {
		return written, err
	}

	return written, nil
}

func downloadHTTPSource(source string, dest string) (int64, error) {
	client := http.Client{
		Timeout: 60 * time.Second,
	}

	request, err := http.NewRequest(http.MethodGet, source, nil)
	if err != nil {
		return 0, err
	}
	request.Header.Set("User-Agent", "orthocal/1.0")

	response, err := client.Do(request)
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("download failed with HTTP status %d", response.StatusCode)
	}

	output, err := os.Create(dest)
	if err != nil {
		return 0, err
	}
	defer output.Close()

	written, err := io.Copy(output, response.Body)
	if err != nil {
		return written, err
	}

	if written == 0 {
		return 0, errors.New("downloaded file is empty")
	}

	if err := output.Sync(); err != nil {
		return written, err
	}

	return written, nil
}

func reject_unsupported_scheme(source string) error {
	parsed, err := url.Parse(source)
	if err != nil {
		return err
	}

	if parsed.Scheme == "" || parsed.Scheme == "http" || parsed.Scheme == "https" {
		return nil
	}

	return fmt.Errorf("unsupported update source scheme: %s", parsed.Scheme)
}

func same_path(source string, target string) (bool, error) {
	sourceAbs, err := filepath.Abs(source)
	if err != nil {
		return false, err
	}

	targetAbs, err := filepath.Abs(target)
	if err != nil {
		return false, err
	}

	return sourceAbs == targetAbs, nil
}
