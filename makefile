NAME    := gcrypto
BINDIR  := bin
DISTDIR := dist

# 纯 Go 项目：CGO_ENABLED=0 保证交叉编译可靠、产物为静态二进制
# -trimpath 去掉本机绝对路径，构建可复现
GOFLAGS := -trimpath -ldflags '-s -w'
GOBUILD := CGO_ENABLED=0 go build $(GOFLAGS)

# 目标平台：GOOS-GOARCH
UNIX_PLATFORMS    := darwin-amd64 darwin-arm64 linux-amd64 linux-arm64
WINDOWS_PLATFORMS := windows-386 windows-amd64
PLATFORMS         := $(UNIX_PLATFORMS) $(WINDOWS_PLATFORMS)

# 实际产物（windows 追加 .exe，统一放 $(BINDIR)）
UNIX_BINS    := $(addprefix $(BINDIR)/$(NAME)-,$(UNIX_PLATFORMS))
WINDOWS_BINS := $(addprefix $(BINDIR)/$(NAME)-,$(WINDOWS_PLATFORMS))
BINS         := $(UNIX_BINS) $(addsuffix .exe,$(WINDOWS_BINS))

# 发布压缩包（放 $(DISTDIR)，与二进制目录分开，避免模式规则冲突）
GZ  := $(addprefix $(DISTDIR)/$(NAME)-,$(addsuffix .gz,$(UNIX_PLATFORMS)))
ZIP := $(addprefix $(DISTDIR)/$(NAME)-,$(addsuffix .zip,$(WINDOWS_PLATFORMS)))

.PHONY: all all-arch releases clean help $(PLATFORMS)

all: $(BINS)
all-arch: all

# 平台别名：make linux-amd64 只构建单个平台
$(UNIX_PLATFORMS):    %: $(BINDIR)/$(NAME)-%
$(WINDOWS_PLATFORMS): %: $(BINDIR)/$(NAME)-%.exe

# 统一构建规则：由目标名解析 GOOS/GOARCH，$@ 已含 .exe 后缀
$(BINDIR)/$(NAME)-%:
	@mkdir -p $(BINDIR)
	@os=$$(echo '$(basename $*)' | cut -d- -f1); \
	arch=$$(echo '$(basename $*)' | cut -d- -f2); \
	echo "==> building $(NAME) for $$os/$$arch"; \
	GOOS=$$os GOARCH=$$arch $(GOBUILD) -o "$@"

# 打包：unix -> gzip，windows -> zip（仅匹配 $(DISTDIR) 下的归档）
$(DISTDIR)/$(NAME)-%.gz: $(BINDIR)/$(NAME)-%
	@mkdir -p $(DISTDIR)
	chmod +x "$<"
	gzip -c "$<" > "$@"

$(DISTDIR)/$(NAME)-%.zip: $(BINDIR)/$(NAME)-%.exe
	@mkdir -p $(DISTDIR)
	zip -j "$@" "$<"

releases: $(GZ) $(ZIP)

clean:
	rm -f $(BINDIR)/$(NAME)* $(DISTDIR)/$(NAME)*

help:
	@echo "make [all]         构建全部平台二进制到 $(BINDIR)/"
	@echo "make <platform>    只构建一个平台，如 make linux-amd64"
	@echo "make releases      构建二进制并打包 gz/zip 到 $(DISTDIR)/"
	@echo "make clean         删除 $(NAME) 的构建产物"
