# gitsod — GitHub 加速小工具

> 一个帮你一键加速GitHub的**克隆、下载** 的命令行小工具(同时还能帮你打开Github镜像站的网址)。   

[Gitee仓库地址](https://gitee.com/idiomeo/gitsod)  (所有的Release文件都在Gitee仓库中发布)  
[Github仓库地址](https://github.com/idiomeo/gitsod)  

---

# 快速上手

## 1. 安装gitsod

  
### Linux
打开终端，直接执行以下命令  
```bash
curl -fsSL https://gitee.com/idiomeo/gitsod/raw/master/install/linux.sh | bash  
```  

该命令将自动拉取**一键安装脚本**进行gitsod的安装  

当脚本执行完毕，此时gitsod就已经被下载并添加为你的系统命令了。  


### Windows

在Windows下，安装并不是很复杂，只是步骤有点多，点击[Windows下的安装教程](./Windows下如何安装.md)，跟着教程慢慢来。   


## 2. 初始化
第一次使用时，需要输入一条指令进行初始化:
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