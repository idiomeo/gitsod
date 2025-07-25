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
	"strings"
)

const (
	configFile = "config.json"
	firstFlag  = ".flag"
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

// ---------- 彩色打印 ----------
var (
	colorRed   = "\033[31m"
	colorGreen = "\033[32m"
	colorReset = "\033[0m"
)

func red(msg string)   { fmt.Println(colorRed + msg + colorReset) }
func green(msg string) { fmt.Println(colorGreen + msg + colorReset) }

// ---------- 命令行参数处理 ----------
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
			red("用法: gitsod clone <GitHub仓库URL>")
			return
		}
		gitClone(strings.Join(os.Args[2:], " "))
	case "download", "-d":
		if len(os.Args) < 3 {
			red("用法: gitsod download <GitHub文件URL>")
			return
		}
		download(strings.Join(os.Args[2:], " "))
	case "open":
		openMirror()
	default:
		openMirror()
	}
}

// ---------- 加载配置 ----------
func loadConfig() (*Config, error) {
	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("找不到 %s，请执行 gitsod update 拉取配置文件", configFile)
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("解析 %s 失败: %v", configFile, err)
	}
	return &cfg, nil
}

// ---------- 更新配置 ----------
func updateConfig() {
	if !commandExists("git") {
		red("本程序依赖 git，请先安装 git 后再试。")
		return
	}

	for _, url := range configURLs {
		fmt.Println("尝试从 " + url + " 拉取配置...")
		resp, err := http.Get(url)
		if err == nil && resp.StatusCode == 200 {
			defer resp.Body.Close()
			out, err := os.Create(configFile)
			if err != nil {
				red("无法写入配置文件: " + err.Error())
				return
			}
			defer out.Close()
			if _, err := io.Copy(out, resp.Body); err != nil {
				red("写入配置文件失败: " + err.Error())
				return
			}
			green("已更新 config.json")
			return
		}
	}
	red("所有镜像源均拉取失败，请检查网络或稍后重试。")
}

// ---------- 克隆 ----------
func gitClone(rawURL string) {
	cfg, err := loadConfig()
	if err != nil {
		red(err.Error())
		return
	}
	if !strings.HasPrefix(rawURL, "http") {
		rawURL = "https://" + rawURL
	}
	target := cfg.ClonePrefix + strings.TrimPrefix(rawURL, "https://")

	if _, err := os.Stat(firstFlag); os.IsNotExist(err) {
		green("首次 clone 需从 GitHub 上缓存镜像地址，请稍等（以后克隆即可从缓存中直接拉取）")
		_ = os.WriteFile(firstFlag, nil, 0644)
	}

	runCmd("git", "clone", target)
}

// ---------- 下载 ----------
func download(rawURL string) {
	cfg, err := loadConfig()
	if err != nil {
		red(err.Error())
		return
	}
	if !strings.HasPrefix(rawURL, "http") {
		rawURL = "https://" + rawURL
	}
	target := cfg.DownloadPrefix[0] + "/" + rawURL
	if !commandExists("wget") && !commandExists("curl") {
		red("系统未检测到 wget 或 curl，请安装其中一个后再试。")
		fmt.Println("可直接使用浏览器下载：", target)
		return
	}
	if commandExists("wget") {
		runCmd("wget", target)
	} else {
		runCmd("curl", "-L", "-O", target)
	}
}

// ---------- 打开镜像 ----------
func openMirror() {
	cfg, err := loadConfig()
	if err != nil {
		red(err.Error())
		return
	}
	url := cfg.MirrorSite
	switch runtime.GOOS {
	case "windows":
		runCmd("cmd", "/c", "start", url)
	case "darwin":
		runCmd("open", url)
	case "linux":
		runCmd("xdg-open", url)
	default:
		fmt.Println("请手动打开浏览器访问:", url)
	}
}

// ---------- 执行外部命令 ----------
func runCmd(name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		red("执行失败: " + err.Error())
	}
}

func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// ---------- 帮助 ----------
func help() {
	fmt.Println(strings.TrimSpace(`
gitsod - GitHub 加速小工具
用法:
  gitsod                    打开镜像站
  gitsod open               同上
  gitsod clone <url>        从镜像克隆仓库
  gitsod download <url>     从镜像下载文件
  gitsod -d <url>           同上
  gitsod update             更新配置文件
  gitsod help | -h          显示此帮助

`))
}
