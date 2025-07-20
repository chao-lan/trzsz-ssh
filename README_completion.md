# tssh Zsh Tab Completion 安装指南

## 方法1：自动安装（推荐）

运行安装脚本：
```bash
./install_completion.sh
```

然后重新加载zsh配置：
```bash
source ~/.zshrc
```

## 方法2：手动安装

### 步骤1：创建补全目录
```bash
mkdir -p ~/.zsh/completions
```

### 步骤2：复制补全文件
```bash
cp _tssh_completion ~/.zsh/completions/_tssh
```

### 步骤3：修改 ~/.zshrc
在您的 `~/.zshrc` 文件中添加以下内容：

```bash
# tssh completion
fpath=(~/.zsh/completions $fpath)
autoload -U compinit && compinit

# Make completion work for ssh alias (which points to tssh)
compdef _tssh ssh
```

### 步骤4：重新加载配置
```bash
source ~/.zshrc
```

## 使用方法

安装完成后，您可以使用以下tab补全功能：

1. **参数补全**：
   ```bash
   ssh --de<TAB>     # 补全为 --debug
   ssh --up<TAB>     # 补全为 --upload-file
   ssh --dr<TAB>     # 补全为 --dragfile
   ssh -<TAB>        # 显示所有可用的短参数
   ```

2. **主机名补全**：
   ```bash
   ssh my<TAB>       # 从 ~/.ssh/config 中补全主机名
   ```

3. **文件路径补全**：
   ```bash
   ssh -i <TAB>      # 补全身份文件路径
   ssh -F <TAB>      # 补全配置文件路径
   ssh --upload-file <TAB>  # 补全要上传的文件路径
   ```

## 测试补全功能

安装完成后，您可以测试以下命令：

```bash
# 测试参数补全
ssh --<TAB><TAB>    # 显示所有长参数选项
ssh -<TAB><TAB>     # 显示所有短参数选项

# 测试特定参数补全
ssh --deb<TAB>      # 应该补全为 --debug
ssh --upl<TAB>      # 应该补全为 --upload-file

# 测试主机名补全（如果您的 ~/.ssh/config 中有配置的主机）
ssh <TAB><TAB>      # 显示所有配置的主机名
```

## 故障排除

如果补全功能不工作，请检查：

1. 确保 `~/.zsh/completions/_tssh` 文件存在
2. 确保 `~/.zshrc` 中包含了必要的配置行
3. 重新启动终端或运行 `source ~/.zshrc`
4. 运行 `which compinit` 确保 compinit 可用

如果仍然有问题，可以运行以下命令进行调试：
```bash
# 检查补全函数是否加载
which _tssh

# 重新初始化补全系统
autoload -U compinit && compinit -D
```