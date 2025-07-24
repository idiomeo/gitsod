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
	configURL  = "https://gitee.com/idiomeo/gitsod-config/raw/master/config.json"
)

type Config struct {
	ClonePrefix    string   `json:"clone_prefix"`
	DownloadPrefix []string `json:"download_prefixes"`
	MirrorSite     string   `json:"mirror_site"`
}

func main() {
	if len(os.Args) == 1 {
		openMirror()
		return
	}

	switch strings.ToLower(os.Args[1]) {
	case "help":
		help()
	case "-h":
		help()
	case "update":
		updateConfig()
	case "clone":
		if len(os.Args) < 3 {
			fmt.Println("用法: gitsod clone <Github文件的URL>")
			return
		}
		gitClone(strings.Join(os.Args[2:], " "))
	case "download":
		if len(os.Args) < 3 {
			fmt.Println("用法: gitsod download <Github文件的URL>")
			fmt.Println("(可以精简为 gitsod -d <Github文件的URL>)")
			return
		}
		download(strings.Join(os.Args[2:], " "))
	case "-d":
		if len(os.Args) < 3 {
			fmt.Println("用法: gitsod d <Github文件的URL>")
			return
		}
		download(strings.Join(os.Args[2:], " "))
	case "open":
		openMirror()
	default:
		openMirror()
	}
}

func loadConfig() (*Config, error) {
	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("找不到 %s，请先执行 gitsod update 拉取配置", configFile)
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("解析 %s 失败: %v", configFile, err)
	}
	return &cfg, nil
}

func updateConfig() {
	if !commandExists("git") {
		fmt.Println("本程序依赖于git，检测到系统未安装 git，请先安装 git 后再试。")
		return
	}

	fmt.Println("正在拉取最新配置...")
	resp, err := http.Get(configURL)
	if err != nil || resp.StatusCode != 200 {
		fmt.Println("拉取配置失败:", err)
		return
	}
	defer resp.Body.Close()

	out, err := os.Create(configFile)
	if err != nil {
		fmt.Println("无法写入配置文件:", err)
		return
	}
	defer out.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		fmt.Println("写入配置文件失败:", err)
		return
	}
	fmt.Println("已更新 config.json")
}

func gitClone(rawURL string) {
	cfg, err := loadConfig()
	if err != nil {
		fmt.Println(err)
		return
	}
	if !strings.HasPrefix(rawURL, "http") {
		rawURL = "https://" + rawURL
	}
	target := cfg.ClonePrefix + strings.TrimPrefix(rawURL, "https://")
	runCmd("git", "clone", target)
}

func download(rawURL string) {
	cfg, err := loadConfig()
	if err != nil {
		fmt.Println(err)
		return
	}
	if !strings.HasPrefix(rawURL, "http") {
		rawURL = "https://" + rawURL
	}
	// 优先使用 ghfast.top
	target := cfg.DownloadPrefix[0] + "/" + rawURL
	if !commandExists("wget") && !commandExists("curl") {
		fmt.Println("系统未检测到 wget 或 curl，请安装其中一个后再试。")
		fmt.Print("你可以直接将以下URL复制进浏览器进行下载： ")
		fmt.Println(target)
		return
	}
	if commandExists("wget") {
		runCmd("wget", target)
	} else {
		runCmd("curl", "-L", "-O", target)
	}
}

func openMirror() {
	cfg, err := loadConfig()
	if err != nil {
		fmt.Println(err)
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

func runCmd(name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("执行失败:", err)
	}
}

func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func help() {

	fmt.Println("以上为全部指令及其用法")
	fmt.Println("更详细信息可以查看官网:gitsod.licnoc.top")

}
