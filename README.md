# milkTrace 牛奶溯源系统后端部分代码
使用Go语言进行开发，合约使用Solidity进行编写

终端执行 

`go build`

`./server`

即可启动服务器，访问入口默认为`http://127.0.0.1:8080/milkTrace/...`
## accounts
用于存储登录区块链的账户

## dal
与数据库交互

## handler
与前端交互

## sdk
区块链的sdk文件

## service
提供公共服务方法

## Tls
ym写的solidity