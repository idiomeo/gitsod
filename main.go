// main.go
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

const (
	configFile   = "config.json"
	firstFlag    = ".flag"
	currentVer   = 1                                                             // 每次发版时手动+1
	versionURL   = "https://gitee.com/idiomeo/gitsod/raw/master/install/VERSION" //VERSION文件的远程地址，用于检测是否需要更新
	downloadBase = "https://gitee.com/idiomeo/gitsod/tree/master/install"        //二进制文件的远程存储地址
)

var configURLs = []string{
	"https://gitee.com/idiomeo/gitsod-config/raw/master/config.json",
	"https://codeberg.org/idiomeo/gitsod/raw/branch/master/config.json",
}

type Config struct {
	ClonePrefix    string   `json:"clone_prefix"`
	DownloadPrefix []string `json:"download_prefixes"`
	MirrorSite     string   `json:"mirror_site"`
}

// ---------- 彩色 ----------
var (
	colorRed   = "\033[31m"
	colorGreen = "\033[32m"
	colorReset = "\033[0m"
)

func red(msg string)   { fmt.Println(colorRed + msg + colorReset) }
func green(msg string) { fmt.Println(colorGreen + msg + colorReset) }

// ---------- 主入口 ----------
func main() {
	if len(os.Args) == 1 {
		openMirror()
		return
	}
	switch strings.ToLower(os.Args[1]) {
	case "help", "-h":
		help()
	case "update":
		updateConfig()
	case "clone":
		if len(os.Args) < 3 {
			red("用法: gitsod clone <repo> [--path [dir]]")
			return
		}
		gitClone(strings.Join(os.Args[2:], " "))
	case "download", "-d":
		if len(os.Args) < 3 {
			red("用法: gitsod download <url> [--path [dir]]")
			return
		}
		download(strings.Join(os.Args[2:], " "))
	case "open":
		openMirror()
	default:
		openMirror()
	}
}

// ---------- 解析 --path ----------
func parsePath(raw string) (dir string, url string) {
	fields := strings.Fields(raw)
	dir, _ = os.Getwd() // 默认当前目录
	for i := 0; i < len(fields); {
		if fields[i] == "--path" {
			if i+1 < len(fields) && !strings.HasPrefix(fields[i+1], "-") {
				dir = fields[i+1]
				fields = append(fields[:i], fields[i+2:]...)
			} else {
				fmt.Print("请输入下载目录（绝对路径）：")
				fmt.Scanln(&dir)
				dir = strings.TrimSpace(dir)
				fields = append(fields[:i], fields[i+1:]...)
			}
			continue
		}
		i++
	}
	url = strings.Join(fields, " ")
	return
}

// ---------- 配置 ----------
func loadConfig() (*Config, error) {
	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("找不到 %s，请执行 gitsod update", configFile)
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("解析 %s 失败: %v", configFile, err)
	}
	return &cfg, nil
}

// ---------- 更新（含自升级） ----------
func updateConfig() {
	// 1. 自升级检查
	resp, err := http.Get(versionURL)
	if err == nil && resp.StatusCode == 200 {
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		remoteStr := strings.TrimSpace(string(body))
		remote, _ := strconv.Atoi(remoteStr)
		if remote > currentVer {
			green("检测到新版本，开始自更新…")
			var installer string
			switch runtime.GOOS {
			case "windows":
				installer = downloadBase + "/install.bat"
				exec.Command("powershell", "-Command",
					"iwr -useb "+installer+" | iex").Run()
			default:
				installer = downloadBase + "/install.sh"
				exec.Command("bash", "-c",
					"curl -fsSL "+installer+" | bash").Run()
			}
			green("自更新完成，请重启终端或重开 gitsod")
			return
		}
	}

	// 2. 更新配置文件
	if !commandExists("git") {
		red("本程序依赖 git，请先安装 git")
		return
	}
	for _, url := range configURLs {
		fmt.Println("尝试从 " + url + " 拉取配置…")
		resp, err := http.Get(url)
		if err == nil && resp.StatusCode == 200 {
			defer resp.Body.Close()
			out, _ := os.Create(configFile)
			io.Copy(out, resp.Body)
			out.Close()
			green("已更新 config.json")
			return
		}
	}
	red("所有镜像源均失败，请检查网络")
	red("若网络正常而镜像源均拉取失败，则联系开发者：idiomeo@foxmail.com")
}

// ---------- clone ----------
func gitClone(raw string) {
	targetDir, url := parsePath(raw)
	if !strings.HasPrefix(url, "http") {
		url = "https://" + url
	}
	cfg, _ := loadConfig()
	target := cfg.ClonePrefix + strings.TrimPrefix(url, "https://")

	if _, err := os.Stat(firstFlag); os.IsNotExist(err) {
		green("首次 clone 需缓存镜像，请稍等…")
		_ = os.WriteFile(firstFlag, nil, 0644)
	}
	_ = os.Chdir(targetDir)
	runCmd("git", "clone", target)
}

// ---------- download ----------
func download(raw string) {
	targetDir, url := parsePath(raw)
	if !strings.HasPrefix(url, "http") {
		url = "https://" + url
	}
	cfg, _ := loadConfig()
	target := cfg.DownloadPrefix[0] + "/" + url

	if !commandExists("wget") && !commandExists("curl") {
		red("未检测到 wget 或 curl")
		fmt.Println("可手动下载：", target)
		return
	}
	_ = os.Chdir(targetDir)
	if commandExists("wget") {
		runCmd("wget", target)
	} else {
		runCmd("curl", "-L", "-O", target)
	}
}

// ---------- open ----------
func openMirror() {
	cfg, _ := loadConfig()
	url := cfg.MirrorSite
	switch runtime.GOOS {
	case "windows":
		runCmd("cmd", "/c", "start", url)
	case "darwin":
		runCmd("open", url)
	case "linux":
		runCmd("xdg-open", url)
	default:
		fmt.Println("请手动打开：", url)
	}
}

// ---------- 执行命令 ----------
func runCmd(name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	cmd.Run()
}
func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// ---------- help ----------
func help() {
	fmt.Println(`gitsod - GitHub 加速小工具
用法:
  gitsod                    打开镜像站
  gitsod open               同上
  gitsod clone <repo> [--path [dir]]
  gitsod download <url> [--path [dir]]
  gitsod -d <url> [--path [dir]]
  gitsod update             更新配置 / 自升级
  gitsod help | -h          显示帮助`)
}
