# xpub
一款管理远程服务器的辅助工具, 批量管理服务器，以及批量执行命令，并且提供与远程服务器的交互式操作

# Feature
1. 通过ssh同时连接多台服务器
2. 自定义命令，可以将一组命令的集合定义成一条命令
3. 提供与远程服务器的交互式shell
4. 提供单台服务器的批量命令执行与多台服务器批量命令执行

# Configuration
配置文件分为两种，一种是xpub运行的配置文件, 可以在启动时通过-conf参数指定, 如config文件夹下的server.json；另一种是自定义的命令集合

配置参数说明:
1. name：服务器的别名
2. host：服务器的IP地址
3. port：服务器的ssh端口号
4. username：服务器ssh用户名
5. passwd：服务器ssh用户名对应的密码
6. commands：自定义命令
    * command：自定义的命令名
    * command_file：自定义命令对应的命令集合文件，如config文件下的command1、command2文件

# Quick Start
## Linux系统
1. 运行build.sh文件
2. 修改config目录下面的配置文件
3. 运行bin目录xpub程序即可

## Windows系统
1. 安装golang, 配置相关环境
2. 安装依赖包
```
go get github.com/mattn/go-colorable
go get github.com/peterh/liner
go get golang.org/x/crypto/ssh
go get golang.org/x/crypto/ssh/terminal
```
3. 修改config目录下面的配置文件
4. 编译xpub并运行即可

# TODO
1. 多种ssh登录方式
2. 脚本文件运程执行
3. 添加Web接口


