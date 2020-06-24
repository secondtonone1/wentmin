## 简介
wentmin是基于wentby改造的轻量级tcp服务器，支持实时socket通信，websocket，http等协议。
## 项目结构
### common
common包用来定义常用的变量和常量
errdef.go中定义了错误码和常用错误
msgdef.go中定义了协议ID，0~1000为系统级别协议，1001以上为用户协议
### components
components为组件包，插件可以写在这里，以后会扩充数据库插件等。目前只做了日志插件读取配置
### config
config目录中存放了服务器的配置信息
### logic
用于实现服务器逻辑处理，包括各种逻辑类和消息对应的回调函数
### netmodel
网络模块，内部实现了tcp处理逻辑，包括
session: 会话层，实现tcp连接的管理。
outmsgque输出队列，目前配置两个，无锁并发输出。
msgqueue 消息队列， 目前配置一个，逻辑处理单线程更安全容易一些，也可以根据配置修改多个，服务器自动实现了加锁。
msghandler 消息处理函数模板，主要用来管理注册消息和回调函数
tcpserver tcp逻辑处理，管理连接，心跳检测，断线处理。tcpserver用于网络协程，当所有协程退出后其才退出
client 客户端封装
webserver 支持http和websocket协议
### example
example中定义了几个测试客户端，包括简单测试，loop测试，长连接压力测试等。
### proto
协议定义在这里，主要是protobuf定义的proto和转化生成的pb文件
### jsonproto
json协议定义在该目录下，websocket用jsonproto通信
### protocol
msgdef定义了消息结构
/*
-----------------------------------------------
               msgpacket
-----------------------------------------------
      msghead     |  msgbody
-----------------------------------------------
id      |   len   |   data
-----------------------------------------------
*/
msghead占用四字节，其中id和len各两字节，id表示消息id，从1001~9999.
len表示后续报文的长度，也就是data中数据的长度
msgbody为消息体，其内部其实是一个字节数组data,具体长度为msghead中len定义的大小

protocol 定义了报文的读写，目前按照大端模式读写。

### main函数
main函数是server的入口函数， 直接go run -race main.go 就可以启动server了。

## 服务器发送消息
如我们要发送"Haaaa"， 服务器先将"Haaaa"通过proto序列化，假设生成的字节流为"hello"
id 为我们定义的消息id，假设为1001, 将1001以大端模式写入两字节
len 为"hello"大小，为5， 将5以大端模式写入两字节
data 为 "hello", 将"hello"以大端模式写入len字节。
这样我们将这9字节发送即可(id2字节+len2字节+data五字节)
## 服务器接收消息
服务器接收先等待包头都接收全再接收包体
先接受4字节，分别通过大端模式，读取两字节id和两字节len，接下来读取len字节的数据data，
接下来将data通过proto反序列化获取实际的数据内容。