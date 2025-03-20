为lot demo
基于https://www.emqx.com/zh/blog/the-easiest-guide-to-getting-started-with-mqtt
可以对于低性能，高延迟的设备环境进行更为便捷的处理

其中该demo功能为，通过1.html实现了发布topic，对于该项目编译成可执行文件，放置于ubuntu中进行运行测试，
1.html通过请求，可对于ubuntu实现打开浏览器中的www.baidu.com的操作
以此类推

管理方面使用官方文档中的MQTTX客户端，可以借用其官方对于websocket的支持，对于实际的逻辑处理进行相对应的封装
