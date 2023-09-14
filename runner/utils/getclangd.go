package utils

import (
	"errors"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/bitly/go-simplejson"
)

const CHECK_CLANGD_INSTALL_PATH = "/userdata/User/globalStorage/llvm-vs-code-extensions.vscode-clangd/install/"

func shellExecutor(shellCmd []string) error {
	cmd := exec.Command(shellCmd[0], shellCmd[1:]...)
	if err := cmd.Start(); err != nil {
		return errors.New("StartShellFailed")
	}

	if err := cmd.Wait(); err != nil {
		return errors.New("ExecShellFailed")
	}

	return nil
}

func findClangDFromTargetDir(targetDirPath string) string {
	var clangdPath string

	if info, err := os.Stat(targetDirPath); err != nil || !info.IsDir() {
		return clangdPath
	}

	var recurFind func(rootDir string)

	recurFind = func(rootDir string) {
		subFiles, err := os.ReadDir(rootDir)

		if err != nil {
			return
		}

		subFolders := []string{}
		for _, subFile := range subFiles {
			subFileInfo, err := os.Stat(path.Join(rootDir, subFile.Name()))

			if err != nil {
				continue
			}

			if !subFileInfo.IsDir() && strings.ToLower(subFile.Name()) == "clangd" {
				clangdPath = path.Join(rootDir, subFile.Name())
				return
			}

			if subFileInfo.IsDir() {
				subFolders = append(subFolders, subFile.Name())
			}
		}

		for _, subFolder := range subFolders {
			recurFind(path.Join(rootDir, subFolder))
		}
	}

	recurFind(targetDirPath)

	return clangdPath
}

func GetClangd() error {

	if info, err := os.Stat("/usr/bin/clangd"); err == nil && !info.IsDir() {
		return nil
	}

	resp, err := http.Get("https://api.github.com/repos/clangd/clangd/releases")

	if err != nil || resp.StatusCode != 200 {
		return errors.New("FailedOnReleaseApi")
	}

	jd, err := simplejson.NewFromReader(resp.Body)

	if err != nil {
		return errors.New("FailedOnUnMarshall")
	}

	releaseList, err := jd.Array()

	if err != nil {
		return errors.New("FailedOnTypeAssert")
	}

	if len(releaseList) == 0 {
		return errors.New("EmptyReleaseList")
	}

	var browserDownloadUrl string
	var tagName string

	for _, releaseType := range releaseList {
		release, ok := releaseType.(map[string]interface{})

		if !ok {
			continue
		}

		prereleaseType, ok := release["prerelease"]

		if !ok {
			continue
		}

		prerelease, ok := prereleaseType.(bool)

		if !ok {
			continue
		}

		if !prerelease {
			nameType, ok := release["name"]

			if !ok {
				continue
			}

			name, ok := nameType.(string)
			tagName = name

			if !ok {
				continue
			}

			assetsType, ok := release["assets"]

			if !ok {
				continue
			}

			assets, ok := assetsType.([]interface{})

			if !ok {
				continue
			}

			for _, assetType := range assets {
				asset, ok := assetType.(map[string]interface{})

				if !ok {
					continue
				}

				assetNameType, ok := asset["name"]

				if !ok {
					continue
				}

				assetName, ok := assetNameType.(string)

				if !ok {
					continue
				}

				if strings.Contains(strings.ToLower(assetName), "linux") {
					burlType, ok := asset["browser_download_url"]

					if !ok {
						continue
					}

					burl, ok := burlType.(string)

					if !ok {
						continue
					}

					browserDownloadUrl = burl

					break
				}
			}

			if len(browserDownloadUrl) > 0 {
				break
			}

		}
	}

	if len(browserDownloadUrl) == 0 || len(tagName) == 0 {
		return errors.New("LinuxAssetNotFound")
	}

	clangdInstallPath := CHECK_CLANGD_INSTALL_PATH + tagName

	info, err := os.Stat(clangdInstallPath)

	createClangD := func() error {
		commands := [][]string{
			{"wget", browserDownloadUrl, "-O", "/progress/clangd.zip"},
			{"unzip", "/progress/clangd.zip", "-d", clangdInstallPath},
			{"rm", "/progress/clangd.zip"},
		}

		for _, command := range commands {
			if err := shellExecutor(command); err != nil {
				return err
			}
		}

		clangdExecutablePath := findClangDFromTargetDir(CHECK_CLANGD_INSTALL_PATH)

		if len(clangdExecutablePath) == 0 {
			return errors.New("ClangDDownloadFailed")
		}

		if err := shellExecutor([]string{"ln", "-sf", clangdExecutablePath, "/usr/bin/clangd"}); err != nil {
			return err
		}

		return nil
	}

	if err != nil || !info.IsDir() {
		if err := os.MkdirAll(clangdInstallPath, os.ModePerm); err != nil {
			return errors.New("CreateClangdInstallDirFailed")
		}

		if err := createClangD(); err != nil {
			return err
		}

	} else if info.IsDir() {
		
		if clangdExecutablePath := findClangDFromTargetDir(CHECK_CLANGD_INSTALL_PATH); len(clangdExecutablePath) > 0 {
			if err := shellExecutor([]string{"ln", "-sf", clangdExecutablePath, "/usr/bin/clangd"}); err == nil {
				return nil
			}
		}

		if err := os.RemoveAll(clangdInstallPath); err != nil {
			return errors.New("ClearOldAssetsFailed")
		}

		if err := createClangD(); err != nil {
			return err
		}
	}

	return nil
}
