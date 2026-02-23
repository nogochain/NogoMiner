# NogoChain 矿工文档 / NogoChain Miner Documentation

## 项目概述 / Project Overview

NogoChain 矿工是一个独立的挖矿软件，用于在 NogoChain 网络上进行挖矿。它使用 NogoPow 共识算法，通过 CPU 进行挖矿操作。

NogoChain Miner is a standalone mining software for mining on the NogoChain network. It uses the NogoPow consensus algorithm and performs mining operations through CPU.

### 主要功能 / Key Features
- 连接到 NogoChain 节点并获取挖矿任务
- 使用多线程进行 CPU 挖矿
- 提交挖矿结果到节点
- 监控哈希率和挖矿状态
- 支持详细的日志记录

- Connect to NogoChain nodes and obtain mining tasks
- Use multi-threading for CPU mining
- Submit mining results to nodes
- Monitor hashrate and mining status
- Support detailed logging

## 启动方法 / Startup Method

### 前提条件 / Prerequisites
- Go 1.18 或更高版本 / Go 1.18 or higher
- 运行中的 NogoChain 节点（启用 RPC 服务） / Running NogoChain node (with RPC service enabled)

### 编译步骤 / Compilation Steps

1. 进入 nogominer-cpu 目录：
   ```bash
   cd nogominer-cpu
   ```

1. Navigate to the nogominer-cpu directory:
   ```bash
   cd nogominer-cpu
   ```

2. 编译矿工软件：
   ```bash
   go build -o miner.exe .
   ```

2. Compile the miner software:
   ```bash
   go build -o miner.exe .
   ```

   或者使用项目根目录的构建脚本：
   ```bash
   go run ../build/ci.go install ./nogominer-cpu
   ```

   Or use the build script from the project root directory:
   ```bash
   go run ../build/ci.go install ./nogominer-cpu
   ```

### 运行命令 / Run Command

基本运行命令：
```bash
./miner.exe
```

Basic run command:
```bash
./miner.exe
```

指定 RPC 服务器和其他参数：
```bash
./miner.exe --rpcaddr 127.0.0.1 --rpcport 8545 --etherbase 0xYourAddress --threads 4 --verbose
```

Specify RPC server and other parameters:
```bash
./miner.exe --rpcaddr 127.0.0.1 --rpcport 8545 --etherbase 0xYourAddress --threads 4 --verbose
```

## 配置方法 / Configuration Method

### 命令行参数 / Command Line Parameters

| 参数 / Parameter | 描述 / Description | 默认值 / Default Value | 环境变量 / Environment Variable |
|------------------|-------------------|------------------------|--------------------------------|
| `--rpcaddr` | RPC 服务器地址 / RPC server address | 127.0.0.1 | NOGO_MINER_RPCADDR |
| `--rpcport` | RPC 服务器端口 / RPC server port | 8545 | NOGO_MINER_RPCPORT |
| `--etherbase` | 矿工地址 / Miner address | 0x0000000000000000000000000000000000000000 | NOGO_MINER_ETHERBASE |
| `--threads` | 挖矿线程数 / Mining threads | 4 | NOGO_MINER_THREADS |
| `--verbose` | 启用详细日志 / Enable verbose logging | false | NOGO_MINER_VERBOSE |
| `--logfile` | 日志文件路径 / Log file path | 空（标准输出） / Empty (standard output) | NOGO_MINER_LOGFILE |

### 环境变量配置 / Environment Variable Configuration

除了命令行参数，还可以通过环境变量来配置矿工：

In addition to command line parameters, you can also configure the miner through environment variables:

```bash
# 设置 RPC 服务器地址 / Set RPC server address
export NOGO_MINER_RPCADDR=127.0.0.1

# 设置 RPC 服务器端口 / Set RPC server port
export NOGO_MINER_RPCPORT=8545

# 设置矿工地址 / Set miner address
export NOGO_MINER_ETHERBASE=0xYourAddress

# 设置挖矿线程数 / Set mining threads
export NOGO_MINER_THREADS=4

# 启用详细日志 / Enable verbose logging
export NOGO_MINER_VERBOSE=true

# 设置日志文件路径 / Set log file path
export NOGO_MINER_LOGFILE=miner.log

# 运行矿工 / Run miner
./miner.exe
```

## 常见问题和故障排除 / Common Issues and Troubleshooting

### 1. 无法连接到 RPC 服务器 / 1. Unable to Connect to RPC Server

**问题 / Issue**：矿工无法连接到 NogoChain 节点的 RPC 服务。

**问题 / Issue**: Miner cannot connect to the RPC service of the NogoChain node.

**解决方案 / Solution**：
- 确保 NogoChain 节点正在运行 / Ensure the NogoChain node is running
- 确保节点启用了 RPC 服务（使用 `--http` 参数） / Ensure the node has RPC service enabled (using `--http` parameter)
- 检查 RPC 服务器地址和端口是否正确 / Check if the RPC server address and port are correct
- 确保防火墙没有阻止连接 / Ensure the firewall is not blocking the connection

### 2. 挖矿速度慢 / 2. Slow Mining Speed

**问题 / Issue**：哈希率较低，挖矿速度慢。

**问题 / Issue**: Low hashrate, slow mining speed.

**解决方案 / Solution**：
- 增加挖矿线程数（`--threads` 参数） / Increase the number of mining threads (`--threads` parameter)
- 确保系统有足够的 CPU 资源 / Ensure the system has sufficient CPU resources
- 关闭其他占用 CPU 的应用程序 / Close other CPU-intensive applications

### 3. 无法找到区块 / 3. Unable to Find Blocks

**问题 / Issue**：矿工运行但没有找到区块。

**问题 / Issue**: Miner is running but not finding blocks.

**解决方案 / Solution**：
- 确保连接到的节点是同步的 / Ensure the connected node is synchronized
- 检查网络哈希率，可能需要更多的算力 / Check the network hashrate, more computing power may be needed
- 验证矿工地址是否正确 / Verify if the miner address is correct

### 4. 日志中出现错误 / 4. Errors in Logs

**问题 / Issue**：日志中出现 RPC 错误或其他错误。

**问题 / Issue**: RPC errors or other errors appear in the logs.

**解决方案 / Solution**：
- 检查 RPC 服务器是否正常运行 / Check if the RPC server is running properly
- 验证命令行参数是否正确 / Verify if the command line parameters are correct
- 查看详细日志以获取更多信息（使用 `--verbose` 参数） / Check detailed logs for more information (using `--verbose` parameter)

## 最佳实践 / Best Practices

1. **使用合适的线程数**：根据系统 CPU 核心数设置线程数，一般建议设置为 CPU 核心数的 75-80%。

1. **Use appropriate number of threads**: Set the number of threads based on the system's CPU cores, generally recommended to be 75-80% of the CPU cores.

2. **监控挖矿状态**：使用 `--verbose` 参数启用详细日志，定期检查哈希率和挖矿状态。

2. **Monitor mining status**: Use the `--verbose` parameter to enable detailed logging and regularly check hashrate and mining status.

3. **设置合适的日志文件**：使用 `--logfile` 参数将日志输出到文件，便于后续分析。

3. **Set appropriate log file**: Use the `--logfile` parameter to output logs to a file for easier analysis later.

4. **定期更新矿工软件**：保持矿工软件更新到最新版本，以获取性能改进和 bug 修复。

4. **Regularly update miner software**: Keep the miner software updated to the latest version to get performance improvements and bug fixes.

5. **确保节点同步**：在开始挖矿前，确保连接的 NogoChain 节点已完全同步。

5. **Ensure node synchronization**: Before starting mining, ensure the connected NogoChain node is fully synchronized.

## 示例配置 / Example Configurations

### 基本配置 / Basic Configuration
```bash
./miner.exe --rpcaddr 127.0.0.1 --rpcport 8545 --etherbase 0xYourAddress --threads 4
```

### 详细日志配置 / Detailed Log Configuration
```bash
./miner.exe --rpcaddr 127.0.0.1 --rpcport 8545 --etherbase 0xYourAddress --threads 4 --verbose --logfile miner.log
```

### 环境变量配置 / Environment Variable Configuration
```bash
export NOGO_MINER_RPCADDR=127.0.0.1
export NOGO_MINER_RPCPORT=8545
export NOGO_MINER_ETHERBASE=0xYourAddress
export NOGO_MINER_THREADS=4
export NOGO_MINER_VERBOSE=true
export NOGO_MINER_LOGFILE=miner.log
./miner.exe
```

## 联系方式 / Contact Information

如果您在使用过程中遇到问题，可以通过以下方式获取帮助：

If you encounter any issues during use, you can get help through the following channels:

- 项目 GitHub 仓库：[NogoChain](https://github.com/nogochain/core-nogo)
- 社区论坛：[NogoChain 社区](https://community.nogochain.org)
- 技术支持邮箱：support@nogochain.org

- Project GitHub Repository: [NogoChain](https://github.com/nogochain/core-nogo)
- Community Forum: [NogoChain Community](https://community.nogochain.org)
- Technical Support Email: support@nogochain.org