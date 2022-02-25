# http包文档阅读

http包提供了client和server的实现.  
http/https的请求可由Get/Head/Post/PostForm函数创建.  
响应的body用完需要关闭:defer resp.Body.Close().  
Client可用于控制http client头,重定向策略.  
Transport可用于控制代理,tls配置,保活,压缩.  
Client/Transport可创建一次,多次复用,协程安全,且高效.  
ListenAndServe启动一个服务,如果不指定handler,默认使用DefaultServeMux.  
Handler/HandleFunc会向DefaultServeMux添加handler.  
如果要控制更多server的细节,用自定义Server.  
Transport/Server自动启用了http2,仅支持简单配置.复杂配置可考虑http2包.

