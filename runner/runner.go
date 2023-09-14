package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path"
	"runner/utils"
	"strings"

	"github.com/bitly/go-simplejson"
)

var coderServerPass string

func init() {
	for _, envStr := range os.Environ() {
		envStrSplit := strings.Split(envStr, "=")
		if envStrSplit[0] == "CODE_PASSWORD" {
			coderServerPass = envStrSplit[1]
			break
		}
	}

	if len(coderServerPass) == 0 {
		fmt.Println("[Runner][Warning] Cannot read password from environment variables, use default instead")
		coderServerPass = "admin"
	}
}

func checkCreateDir(dirName, dirPath string) error {
	fmt.Printf("[Runner][Info]    Start to check code-server %s directory\n", dirName)

	if info, err := os.Stat(dirPath); err != nil || !info.IsDir() {
		fmt.Printf("[Runner][Info]    Creating code-server %s directory...\n", dirName)
		err := os.MkdirAll(dirPath, os.ModePerm)

		if err != nil {
			fmt.Printf("[Runner][Error]   Cannot create code-server %s directory\n", dirName)
			return err
		}
	}

	return nil
}

func checkWriteFile(fileName, filePath, fileContent string) error {
	fmt.Printf("[Runner][Info]    Start to check code-server %s file exists\n", fileName)

	if info, err := os.Stat(filePath); err == nil && !info.IsDir() {
		fmt.Printf("[Runner][Info]    Found code-server %s file\n", fileName)
		return nil
	}

	fmt.Printf("[Runner][Info]    Start to write code-server %s file\n", fileName)
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)

	if err != nil {
		fmt.Printf("[Runner][Error]   Cannot open write stream to code-server %s file\n", fileName)
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	if _, err := writer.WriteString(fileContent); err != nil {
		fmt.Printf("[Runner][Error]   Write to code-server %s file failed\n", fileName)
		return err
	}

	if err := writer.Flush(); err != nil {
		fmt.Printf("[Runner][Error]   Do flush to code-server %s file failed\n", fileName)
		return err
	}

	return nil
}

const USER_DATA_DIR string = "/userdata"
const EXT_DATA_DIR string = "/extensions"

func main() {
	fmt.Println("[Runner][Info]    Preparing EmCode running environment...")

	fmt.Println("[Runner][Info]    Start to write code-server config.yaml file")

	home, err := os.UserHomeDir()

	if err != nil {
		fmt.Println("[Runner][Error]   Cannot read user home directory, aborting...")
		os.Exit(-1)
	}

	configStoreDirectory := path.Join(home, ".config", "code-server")

	if err := checkCreateDir(".config", configStoreDirectory); err != nil {
		fmt.Println("[Runner][Error]   Check and create code-server directory .config failed, aborting...")
		os.Exit(-1)
	}

	yamlConfigFilePath := path.Join(configStoreDirectory, "config.yaml")
	yamlContent := []string{
		"bind-addr: 0.0.0.0:8080",
		"auth: password",
		"password: " + coderServerPass,
		"cert: false",
	}

	if err := checkWriteFile("config.yaml", yamlConfigFilePath, strings.Join(yamlContent, "\n")); err != nil {
		fmt.Println("[Runner][Error]   Check and create code-server file config.yaml failed, aborting...")
		os.Exit(-1)
	}

	settingJsonParentDir := path.Join(USER_DATA_DIR, "User")

	if err := checkCreateDir("userdata", USER_DATA_DIR); err != nil {
		fmt.Println("[Runner][Error]   Check and create code-server directory userdata failed, aborting...")
		os.Exit(-1)
	}

	if err := checkCreateDir("setting.json directory", settingJsonParentDir); err != nil {
		fmt.Println("[Runner][Error]   Check and create code-server setting.json directory failed, aborting...")
		os.Exit(-1)
	}

	if err := checkCreateDir("extenstion data", EXT_DATA_DIR); err != nil {
		fmt.Println("[Runner][Error]   Check and create code-server directory extenstion data failed, aborting...")
		os.Exit(-1)
	}

	sj := simplejson.New()

	sj.Set("workbench.colorTheme", "Visual Studio Dark")
	sj.Set("editor.formatOnSave", true)
	sj.Set("[cpp]", map[string]string{"editor.defaultFormatter": "llvm-vs-code-extensions.vscode-clangd"})
	sj.Set("[javascript]", map[string]string{"editor.defaultFormatter": "esbenp.prettier-vscode"})
	sj.Set("[html]", map[string]string{"editor.defaultFormatter": "esbenp.prettier-vscode"})
	sj.Set("[css]", map[string]string{"editor.defaultFormatter": "esbenp.prettier-vscode"})
	sj.Set("[typescript]", map[string]string{"editor.defaultFormatter": "esbenp.prettier-vscode"})
	sj.Set("[javascriptreact]", map[string]string{"editor.defaultFormatter": "esbenp.prettier-vscode"})
	sj.Set("[python]", map[string]string{"editor.defaultFormatter": "ms-python.python"})

	defaultSettingsBytes, err := sj.MarshalJSON()

	if err != nil {
		fmt.Println("[Runner][Error]   Cannot create code-server default setting string, aborting...")
		os.Exit(-1)
	}

	settingJsonPath := path.Join(settingJsonParentDir, "settings.json")

	if err := checkWriteFile("settings.json", settingJsonPath, string(defaultSettingsBytes)); err != nil {
		fmt.Println("[Runner][Error]   Check and create code-server file settings.json failed, aborting...")
		os.Exit(-1)
	}

	extensions := []string{
		"esbenp.prettier-vscode",
		"llvm-vs-code-extensions.vscode-clangd",
		"ms-python.python",
	}

	fmt.Println("[Runner][Info]    Start to install extensions...")

	isInstallSucceed := make(chan bool)

	go func() {
		for _, ext := range extensions {
			cmd := exec.Command("code-server", "--user-data-dir", USER_DATA_DIR, "--extensions-dir", EXT_DATA_DIR, "--install-extension", ext)

			procStdout, pipeErr := cmd.StdoutPipe()
			cmd.Stderr = cmd.Stdout

			if err := cmd.Start(); err != nil {
				fmt.Printf("[Runner][Error]   Cannot start extension %s installation process\n", ext)
				isInstallSucceed <- false
				return
			}

			if pipeErr == nil {
				go func() {
					reader := bufio.NewReader(procStdout)
					for {
						line, lineErr := reader.ReadString('\n')

						if lineErr != nil {
							break
						}

						fmt.Printf("[Coder][Info]     %s", line)
					}
				}()
			}

			if err := cmd.Wait(); err != nil {
				fmt.Printf("[Runner][Error]   Cannot finishize extension %s installation process\n", ext)
				isInstallSucceed <- false
				return
			}
		}
		fmt.Println("[Runner][Info]    Extensions all installed")
		isInstallSucceed <- true
	}()

	if !(<-isInstallSucceed) {
		fmt.Println("[Runner][Error]   install extension failed, aborting...")
		os.Exit(-1)
	}

	fmt.Println("[Runner][Info]    Starting install clangd...")
	if err := utils.GetClangd(); err != nil {
		fmt.Printf("[Runner][Error]   install clangd failed, reason: %s, aborting...\n", err.Error())
		os.Exit(-1)
	}

	fmt.Println("[Runner][Info]    All process done! Starting code-server instance...")

	serverStopSignal := make(chan bool)
	go func() {
		cmd := exec.Command("code-server", "/src", "--user-data-dir", USER_DATA_DIR, "--extensions-dir", EXT_DATA_DIR)

		if err := cmd.Start(); err != nil {
			fmt.Println("[Runner][Error]    Cannot start code-server instance, aborting...")
			os.Exit(-1)
		}

		cmd.Wait()

		serverStopSignal <- true
	}()

	fmt.Println("[Runner][Info]    code-server started, processing normal requests...")
	<-serverStopSignal
	fmt.Println("[Runner][Info]    code-server terminated")
}
