# gitsod — GitHub 加速小工具

> 一个帮你无痛加速克隆、下载和访问 GitHub 的命令行小工具。

![吸引人用的图片](./picture/attract.jpeg)

---

# 快速上手

## 1. 一键安装gitsod
通过安装脚本下载

## 2. 初始化
第一次使用 只需一步：
``` bash
# 拉取最新镜像配置（仅需一次）
gitsod update
```

成功后会生成 config.json，保存了当前可用的镜像地址。

## 3. 开始使用
- **克隆仓库**
  ```bash
  gitsod clone github.com/tendermint/tendermint.git
  ```
  首次克隆会提示“首次缓存镜像，请稍等”，以后再克隆同一仓库即可获得非常快的速度。

- **下载文件**
  ```bash
  gitsod download github.com/rustdesk/rustdesk/releases/download/1.4.0/rustdesk-1.4.0-x86_64.exe
  ```

- **打开Github镜像站**
    ```bash
  gitsod open
  ```

---

# 命令汇总
| 命令 | 简写 | 说明 |
|------|------|------|
| `gitsod open` | `gitsod` | 打开Github镜像站 |
| `gitsod clone <repo>` | — | 克隆仓库（支持简写 `user/repo`） |
| `gitsod download <url>` | `gitsod -d <url>` | 下载文件 |
| `gitsod update` | — | 更新gitsod本体/更新镜像配置 |
| `gitsod help` | `gitsod -h` | 查看帮助 |

---

# 常见问题

### 1. 找不到 `config.json`？
执行 `gitsod update` 即可自动拉取。

### 2. 系统没有 `git` / `wget` / `curl`？
- `git` 是 **必须** 的，请先安装。  
- 下载文件需要 `wget` 或 `curl`，如两者都没有，程序会给出可直接复制到浏览器的加速链接。

---

# 开源协议
本程序采用[**Apache License 2.0**](./LICENSE)进行代码分发。  