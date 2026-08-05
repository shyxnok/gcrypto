# gcrypto

基于 XXTEA 算法的文件加解密命令行小工具。单个静态二进制，开箱即用。

## 特性

- 🔐 **XXTEA 对称加密**：块长可变的分组加密算法，内置固定密钥
- 📦 **任意长度文件**：自动 PKCS#7 填充，支持任意二进制内容
- 🧱 **纯 Go 实现**：`CGO_ENABLED=0` 交叉编译，产物为无依赖的静态二进制
- 🔨 **一键跨平台构建**：makefile 支持 6 个平台，可增量构建、打包发布

## 快速开始

```bash
# 构建本平台二进制到 bin/
make

# 加密文件
./bin/gcrypto encrypt secret.txt enc

# 解密文件
./bin/gcrypto decrypt secret.enc

# 便捷命令：elh 加密为 .lh，dlh 解密还原
./bin/elh secret.txt
./bin/dlh secret.lh
```

## 使用

```
gcrypto encrypt <srcPath> <suffix>    # 加密为任意后缀（如 enc/lh/xxx）
gcrypto decrypt <srcPath>             # 解密并自动还原原文件名
elh <srcPath>                         # 加密，原后缀写入文件头，输出同名 .lh
dlh <srcPath>                         # 解密，读文件头后缀并还原原文件名
```

`.lh` 文件的格式：

- 加密时把原文件后缀（如 `.hxl`）以 `[1 字节长度][后缀]` 写入 `.lh` 文件首部，正文用与
  `gcrypto encrypt` 相同的算法加密，输出文件名为去掉原后缀后加上指定后缀
  （`gcrypto encrypt my.hxl enc` → `my.enc`；`elh my.hxl` → `my.lh`）
- 解密时从首部读出原后缀，还原文件名（`my.lh` → `my.hxl`）
- 加密或解密的目标文件已存在时，直接删除旧的再写入
- `gcrypto encrypt/decrypt` 与 `elh/dlh` 使用同一格式，可互相解密

| 说明 | 行为 |
|---|---|
| 成功 | 输出 `File encrypted/decrypted successfully.`，退出码 0 |
| 参数不足 | 打印用法并返回退出码 1 |
| 文件不存在 / 解密失败 | 打印错误信息，退出码 1 |

示例：

```bash
# 加密一篇文档为任意后缀
./bin/gcrypto encrypt diary.txt enc

# 还原（自动恢复为 diary.txt）
./bin/gcrypto decrypt diary.enc

# 用便捷命令：elh 加密，dlh 还原（自动保留/恢复后缀）
./bin/elh diary.txt            # 生成 diary.lh
./bin/dlh diary.lh             # 还原为 diary.txt
```

> ⚠️ 使用内置固定密钥，加密主要用于防误读（普通人打开为乱码），**不提供真正的安全强度**。若需强加密，请自行基于 `secbox` 包传入独立密钥。

## 构建与发布

| 命令 | 说明 |
|---|---|
| `make` | 构建全部 6 个平台到 `bin/` |
| `make tools` | 构建便捷命令 `elh`/`dlh` 到 `bin/` |
| `make linux-amd64` | 只构建单个平台（`darwin-amd64` / `darwin-arm64` / `linux-amd64` / `linux-arm64` / `windows-386` / `windows-amd64`） |
| `make releases` | 构建并打包：unix 平台生成 `.gz`，windows 生成 `.zip`，归档到 `dist/` |
| `make clean` | 清理构建产物 |
| `make help` | 查看全部目标 |

构建参数：`CGO_ENABLED=0`（静态链接）+ `-trimpath`（可复现构建）+ `-ldflags '-s -w'`（去除符号表、减小体积）。

## 加密算法

- **算法**：XXTEA（eXtended Tiny Encryption Algorithm）
- **密钥**：内置固定 16 字节密钥（包初始化时派生一次）
- **填充**：PKCS#7，加密前补齐，解密时校验
- **接口**：`secbox` 包提供通用 API——`Encrypt`/`Decrypt`（可指定密钥、填充、轮数）、`EncryptText`/`DecryptText`（内置密钥）、`EncryptBase64`/`EncryptHex`（编码输出）、`EncryptFile`/`DecryptFile`（文件封装）、`URandom`（随机字节）

## 目录结构

```
gcrypto/
├── main.go            入口，调用 cli.Execute()
├── cli/               package cli：cobra 命令定义（encrypt/decrypt）
│   └── cli.go
├── elh/               package elh：.lh 加密包（后缀头 + 输出 .lh）
│   ├── elh.go
│   └── elh_test.go
├── dlh/               package dlh：.lh 解密包（读后缀头 + 还原文件名）
│   ├── dlh.go
│   └── dlh_test.go
├── cmd/
│   ├── elh/           elh 命令行入口
│   └── dlh/           dlh 命令行入口
├── secbox/            XXTEA 加解密核心（package secbox）
│   ├── crypto.go      XXTEA 加密核心
│   ├── secbox.go      EncryptFile / DecryptFile 等文件 API
│   ├── crypto_test.go XXTEA 单元测试
│   └── secbox_test.go 文件往返测试
├── makefile           跨平台构建脚本
└── go.mod / go.sum
```

## 测试

```bash
go test ./...     # 全部测试
go test ./secbox/ -cover   # 覆盖率
```

## Git

项目使用 Git 管理，代码在 `main` 分支：

```bash
git status     # 查看改动
git diff       # 查看改动内容
git add .      # 暂存改动
git commit -m "提交说明"   # 提交
```

提交说明建议使用 `类型: 简述` 的格式：

| 类型 | 用途 | 示例 |
|---|---|---|
| `feat` | 新功能 | `feat: 新增 elh/dlh 加解密命令` |
| `fix` | 修复问题 | `fix: 支持任意后缀的解密还原` |
| `docs` | 文档改动 | `docs: 更新 README` |
| `refactor` | 重构 | `refactor: CLI 独立为 cli 包` |

> 注意：`my.hxl` 是个人文件，不会纳入版本控制，提交时请勿将其加入暂存区。
