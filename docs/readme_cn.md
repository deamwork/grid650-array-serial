# Grid650 Array 模块串口通信工具

For english version, please visit：[grid650 array module serial tool readme](../readme.md)

## 特性与集成计划

- [x] 自定义文本
  - [ ] 文本内容验证
- [ ] 时间同步
- [ ] 系统性能信息同步

## 配置

该配置简要设定了 http 服务和设备连接的描述信息

示例配置文件如下:

```yaml
http:
  listen: "0.0.0.0:80"
  tls:
    enable: false
    listen: "0.0.0.0:443"
    certificate_chain: "/path/to/your/cert.pem"
    private_key: "/path/to/your/cert.pem"
device:
  name: "/dev/tty.usbmodem14233301"
  baud: 115200
```

- `http` 定义服务器运行的必要信息，目前暂未实现。
- `device` 定义了需要连接的 grid650 设备

### HTTP 服务器

待补充

### Grid650 设备

你可以使用如下步骤，找到配置中需要的 grid650 信息:

#### MacOS

1. 连接 grid650
2. 打开终端 (`iTerm2.app` 等类似工具均可)
3. 输入 `ls -l /dev/tty*` 列出可用串口通信设备
4. 执行后可能会列出非常多的设备，请确保上下文可以用于搜索
5. 查找文件名类似与 `tty.usbmodem14233301` 的文件 (数字部分可能应设备不同而有差异)
6. 找到后，拼装完整路径，类似 `/dev/tty.usbmodem14233301`，然后填入配置文件
7. 填写 Baud 字段时，根据用户手册， grid650 array 的 bit rate 应为 115200。输入该段时，无需引号

#### Windows

待补充

#### Linux & *Unix

待补充

### 使用

```text
$ ./grid650-array-serial
grid 650 array module serial-comm server

Usage:
  grid650-array-serial [flags]
  grid650-array-serial [command]

Available Commands:
  config      Print config parse result
  help        Help about any command
  send        Send one custom text to the grid650 array device
  version     Print version information

Flags:
      --baud int                       Specific device bit rate [env BAUD]
      --config string                  Specific config file path. (default "../config/config.yaml")
      --device string                  Specific device tty [env DEVICE] (default "/dev/tty.usbmodem14233301")
  -d, --environment-development-mode   change environment mode to development [env ENVIRONMENT_DEVELOPMENT_MODE]
  -j, --environment-json-log-style     change log style to JSON [env ENVIRONMENT_JSON_LOG_STYLE]
  -h, --help                           help for grid650-array-serial
      --http-listen string             Specific http listen address. (default "0.0.0.0:80")

Use "grid650-array-serial [command] --help" for more information about a command.
```

#### 快速体验

发送自定义信息:
```bash
./grid650-array-serial send "i am grid650"
```

## 致谢
 - https://github.com/tarm/serial
 - https://github.com/XSAM/go-hybrid

## 交流群

QQ: 964358671

## 引用参考文档
 - [grid_ARRAY_模块数据结构说明_20210331X (PDF)](./grid_ARRAY_模块数据结构说明_20210331X.pdf)