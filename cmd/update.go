package cmd

import (
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

func init() {
	if runtime.GOOS == "windows" {
		// Windows currently does not work.
		return
	}

	RootCmd.AddCommand(updateCmd)
	updateCmd.Flags().BoolVarP(&optUpdateCmdYes, "yes", "y", false, "Confirm update without prompt, if available")
}

var optUpdateCmdYes bool
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update lvl to the latest version",
	Long: `Update lvl to the latest version
lvl will automatically download the latest version and replace the installed executable.`,

	RunE: func(cmd *cobra.Command, args []string) error {
		apiClient := http.Client{}

		releases, err := getLvlReleases(&apiClient)
		if err != nil {
			return fmt.Errorf("failed to check for releases on GitHub: %v", err)
		}

		fmt.Printf("Current version: %s\n", strings.TrimSpace(version))

		// I assume these are ordered as most recent, right?
		release := releases[0]

		if release.TagName == strings.TrimSpace(version) {
			fmt.Println("Up to date, no update needed!")
			return nil
		}

		fmt.Printf("New version of lvl found: %s\n", release.TagName)

		if !optUpdateCmdYes {
			if !confirmPrompt("Do you want to continue with this update?") {
				return nil
			}
		}

		fmt.Println("Updating...")

		assets, err := getReleaseAssets(&apiClient, release.AssetsUrl)
		if err != nil {
			return fmt.Errorf("failed to check for release assets on GitHub: %v", err)
		}

		assetName := getAssetFileName()

		var ourAsset GitHubReleaseAsset
		found := false
		for _, asset := range assets {
			if asset.Name == assetName {
				ourAsset = asset
				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf("unable to find version of new release for this operating system")
		}

		execPath, err := getExecutablePath()
		if err != nil {
			return fmt.Errorf("unable to find location of lvl executable: %v", err)
		}

		newExecPath := execPath + ".new"
		err = downloadNewExecutable(&apiClient, ourAsset.Url, newExecPath)
		if err != nil {
			return fmt.Errorf("failed to download new release asset: %v", err)
		}

		if runtime.GOOS == "windows" {
			panic("Windows updates currently unimplemented")
		} else {
			err := os.Rename(newExecPath, execPath)
			if err != nil {
				return fmt.Errorf("failed to move new executable into place: %v", err)
			}
		}

		fmt.Println("Update success!")

		return nil
	},
}

func downloadNewExecutable(apiClient *http.Client, url string, destination string) error {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	setGitHubApiHeaders(request)
	request.Header.Set("Accept", "application/octet-stream")

	resp, err := apiClient.Do(request)
	if err != nil {
		return err
	}

	err = checkHttpStatus(resp)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	file, err := os.OpenFile(destination, os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func checkHttpStatus(response *http.Response) error {
	if response.StatusCode < 200 || response.StatusCode > 299 {
		return fmt.Errorf("bad HTTP response code: %d", response.StatusCode)
	}

	return nil
}

func getExecutablePath() (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", err
	}

	execPath, err = filepath.EvalSymlinks(execPath)
	if err != nil {
		return "", err
	}

	return execPath, nil
}

func getLvlReleases(client *http.Client) ([]GitHubRelease, error) {
	request, err := http.NewRequest("GET", "https://api.github.com/repos/level27/lvl/releases", nil)
	if err != nil {
		return nil, err
	}

	setGitHubApiHeaders(request)

	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	err = checkHttpStatus(resp)
	if err != nil {
		return nil, err
	}

	return readJson[[]GitHubRelease](resp)
}

func getReleaseAssets(client *http.Client, url string) ([]GitHubReleaseAsset, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	setGitHubApiHeaders(request)

	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	err = checkHttpStatus(resp)
	if err != nil {
		return nil, err
	}

	return readJson[[]GitHubReleaseAsset](resp)
}

func setGitHubApiHeaders(request *http.Request) {
	request.Header.Set("User-Agent", getUserAgent())
	request.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	request.Header.Set("Accept", "application/vnd.github+json")
}

func readJson[T any](response *http.Response) (T, error) {
	var result T

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return result, err
	}

	return result, nil
}

// Get the appropriate release asset file name for our
func getAssetFileName() string {
	arch := runtime.GOARCH
	platform := runtime.GOOS

	if platform == "windows" {
		return fmt.Sprintf("lvl-windows-%s.exe", arch)
	}

	if platform == "linux" {
		return fmt.Sprintf("lvl-linux-%s", arch)
	}

	if platform == "darwin" {
		return fmt.Sprintf("lvl-darwin-%s", arch)
	}

	panic("unknown platform for auto update!")
}

type GitHubRelease struct {
	TagName   string `json:"tag_name"`
	AssetsUrl string `json:"assets_url"`
}

type GitHubReleaseAsset struct {
	Url  string `json:"url"`
	Name string `json:"name"`
}
