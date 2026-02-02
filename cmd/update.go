package cmd

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

type githubRelease struct {
	TagName string        `json:"tag_name"`
	Assets  []githubAsset `json:"assets"`
}

type githubAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "최신 버전으로 업데이트",
	Long:  "GitHub Releases에서 최신 버전을 다운로드하여 바이너리를 업데이트합니다.",
	RunE:  runUpdate,
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

func runUpdate(cmd *cobra.Command, _ []string) error {
	fmt.Fprintln(cmd.OutOrStdout(), "최신 버전 확인 중...")

	release, err := fetchLatestRelease()
	if err != nil {
		return fmt.Errorf("최신 버전 확인 실패: %w", err)
	}

	latestVersion := strings.TrimPrefix(release.TagName, "v")
	currentVersion := parseCurrentVersion(rootCmd.Version)

	if currentVersion != "dev" && currentVersion == latestVersion {
		fmt.Fprintf(cmd.OutOrStdout(), "이미 최신 버전입니다 (v%s)\n", latestVersion)
		return nil
	}

	asset, err := findAsset(release, latestVersion)
	if err != nil {
		return err
	}

	fmt.Fprintln(cmd.OutOrStdout(), "다운로드 중...")

	archiveData, err := downloadAsset(asset.BrowserDownloadURL)
	if err != nil {
		return fmt.Errorf("다운로드 실패: %w", err)
	}

	binaryName := "bible"
	if runtime.GOOS == "windows" {
		binaryName = "bible.exe"
	}

	var newBinary []byte
	if strings.HasSuffix(asset.Name, ".zip") {
		newBinary, err = extractFromZip(archiveData, binaryName)
	} else {
		newBinary, err = extractFromTarGz(archiveData, binaryName)
	}
	if err != nil {
		return fmt.Errorf("압축 해제 실패: %w", err)
	}

	fmt.Fprintln(cmd.OutOrStdout(), "설치 중...")

	if err := replaceBinary(newBinary); err != nil {
		return fmt.Errorf("바이너리 교체 실패: %w", err)
	}

	displayCurrent := currentVersion
	if displayCurrent == "dev" {
		displayCurrent = "dev"
	} else {
		displayCurrent = "v" + displayCurrent
	}
	fmt.Fprintf(cmd.OutOrStdout(), "%s → v%s 업데이트 완료\n", displayCurrent, latestVersion)
	return nil
}

// parseCurrentVersion extracts the semver portion from rootCmd.Version.
// "1.2.1" → "1.2.1", "1.2.1 (abc1234)" → "1.2.1", "dev" → "dev".
func parseCurrentVersion(v string) string {
	if v == "" || v == "dev" {
		return "dev"
	}
	if idx := strings.Index(v, " "); idx != -1 {
		v = v[:idx]
	}
	return strings.TrimPrefix(v, "v")
}

func fetchLatestRelease() (*githubRelease, error) {
	req, err := http.NewRequest(http.MethodGet,
		"https://api.github.com/repos/yangsijun/bible-tui/releases/latest", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "bible-tui")
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("네트워크 오류: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API 응답 오류: %s", resp.Status)
	}

	var release githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("응답 파싱 오류: %w", err)
	}
	return &release, nil
}

func findAsset(release *githubRelease, version string) (*githubAsset, error) {
	ext := ".tar.gz"
	if runtime.GOOS == "windows" {
		ext = ".zip"
	}
	wantName := fmt.Sprintf("bible_%s_%s_%s%s", version, runtime.GOOS, runtime.GOARCH, ext)

	for i := range release.Assets {
		if release.Assets[i].Name == wantName {
			return &release.Assets[i], nil
		}
	}
	return nil, fmt.Errorf("현재 플랫폼(%s/%s)에 맞는 릴리스 파일을 찾을 수 없습니다: %s",
		runtime.GOOS, runtime.GOARCH, wantName)
}

func downloadAsset(url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "bible-tui")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("네트워크 오류: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("다운로드 오류: %s", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("읽기 오류: %w", err)
	}
	return data, nil
}

func extractFromTarGz(data []byte, target string) ([]byte, error) {
	gr, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer gr.Close()

	tr := tar.NewReader(gr)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if filepath.Base(hdr.Name) == target && hdr.Typeflag == tar.TypeReg {
			return io.ReadAll(tr)
		}
	}
	return nil, fmt.Errorf("아카이브에서 %s를 찾을 수 없습니다", target)
}

func extractFromZip(data []byte, target string) ([]byte, error) {
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, err
	}

	for _, f := range zr.File {
		if filepath.Base(f.Name) == target && !f.FileInfo().IsDir() {
			rc, err := f.Open()
			if err != nil {
				return nil, err
			}
			defer rc.Close()
			return io.ReadAll(rc)
		}
	}
	return nil, fmt.Errorf("아카이브에서 %s를 찾을 수 없습니다", target)
}

func replaceBinary(newBinary []byte) error {
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("현재 바이너리 경로를 찾을 수 없습니다: %w", err)
	}
	execPath, err = filepath.EvalSymlinks(execPath)
	if err != nil {
		return fmt.Errorf("심볼릭 링크 해석 실패: %w", err)
	}

	info, err := os.Stat(execPath)
	if err != nil {
		return err
	}
	perm := info.Mode().Perm()

	oldPath := execPath + ".old"

	if err := os.Rename(execPath, oldPath); err != nil {
		return fmt.Errorf("바이너리에 쓰기 권한이 없습니다 (sudo가 필요할 수 있습니다): %w", err)
	}

	if err := os.WriteFile(execPath, newBinary, perm); err != nil {
		_ = os.Rename(oldPath, execPath)
		return fmt.Errorf("새 바이너리 쓰기 실패: %w", err)
	}

	_ = os.Remove(oldPath)
	return nil
}
