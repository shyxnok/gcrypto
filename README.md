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
./bin/gcrypto encrypt secret.txt secret.bin

# 解密文件
./bin/gcrypto decrypt secret.bin secret.txt
```

## 使用

```
gcrypto encrypt <srcPath> <dstPath>   # 加密文件
gcrypto decrypt <srcPath> <dstPath>   # 解密文件
```

| 说明 | 行为 |
|---|---|
| 成功 | 输出 `File encrypted/decrypted successfully.`，退出码 0 |
| 参数不足 | 打印用法并返回退出码 1 |
| 文件不存在 / 解密失败 | 打印错误信息，退出码 1 |

示例：

```bash
# 加密一篇文档
./bin/gcrypto encrypt diary.txt diary.enc

# 还原（必须在加密后的文件上执行）
./bin/gcrypto decrypt diary.enc diary-backup.txt
```

> ⚠️ 使用内置固定密钥，加密主要用于防误读（普通人打开为乱码），**不提供真正的安全强度**。若需强加密，请自行基于 `secbox` 包传入独立密钥。

## 构建与发布

| 命令 | 说明 |
|---|---|
| `make` | 构建全部 6 个平台到 `bin/` |
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
├── main.go            入口，调用 secbox.Execute()
├── secbox/            全部业务代码（package secbox）
│   ├── cli.go         rootCmd + Execute + encrypt/decrypt 命令
│   ├── crypto.go      XXTEA 加密核心
│   ├── secbox.go      EncryptFile / DecryptFile
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
