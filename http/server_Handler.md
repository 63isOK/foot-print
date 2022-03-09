# Handler

Handler用于处理http的请求,有以下几点要注意:

1. ServeHTTP需要将响应的头和数据写入到ResponseWriter,并发返回,返回意味这请求完成.
2. 在ServeHTTP调用完成的同时或完成之后,使用ResponseWriter或从请求体中读数据都是无效的.
3. 不同的http客户端,不同的http版本以及中间人,在向ResponserWriter写之后,再读请求.
   体基本不可能了.正确的做法是先读请求体,再响应.
4. handler不应该修改请求.
5. 如果ServeHTTP panic了,panic的影响和请求是隔离的.
6. 如果recover panic了,就要打印堆栈日志并关闭网络连接.
7. 如果要取消一个handler,客户端会得到中断响应,服务端不会记录错误,而是以
   ErrAbortHandler进行panic.

Handler是一个接口

```go
type Handler interface {
  ServeHTTP(ResponseWriter, *Request)
}

// 请求是一个结构体,在ServeHTTP中以指针形式出现
// 服务端收的或客户端发的都叫请求
// 客户端和服务端的请求有细微差别
type Request struct {
  // http方法 get/post/put等,空表示get
  // go的http客户端不支持connect方法,connect用于代理(翻墙)
  Method string

  // 要请求的uri(服务端);要访问的url(客户端)
  // 在服务端,
  //    url存的是从RequestURI解析之后的uri,
  //    大部分场合,URL.Path/URL.RawQuery之外的字段都为空
  // 在客户端,
  //    URL.Host表示要连接的server,会出现在头信息
  //    Request.Host是可选的,如果指定了,就会覆盖头信息中URL.Host
  URL *url.URL

  // 服务端即将接收的请求版本
  // 客户端会忽略这几个字段,取值范围是HTTP/1/1或HTTP/2
  Proto      string // "HTTP/1.0"
  ProtoMajor int    // 1
  ProtoMinor int    // 0

  // 头信息, 存储格式为map[string][]string
  // 在服务端,会将头信息中的Host提取至Request.Host,并从map中移除
  // 头信息名大小写不敏感
  //    解析的实现是统一的:首字母和-后的字母大写,其他小写
  // 在客户端,Content-Length/Connection在需要时自动添加,对应的值可能被忽略
  Header Header

  // 请求体
  // 在客户端,
  //    空请求体表示没有body,get请求就没有body
  //    Transport负责调用Body.Close
  // 在服务端,
  //    body永远非空,如果遇到请求体是空,则立马返回EOF
  //    server会调用Body.Close,所以在handler中不需要调用Request.Body.Close
  Body io.ReadCloser

  // 客户端专用,可选,获取body的拷贝
  // 适用场景:重定向请求多次读body
  GetBody func() (io.ReadCloser, error)

  // 内容长度,-1表示未知
  // 大于0时,表示可以从body中读取的字节数
  // 在客户端,在body非空的情况下,内容长度为0同样表示长度未知
  ContentLength int64

  // 传输编码列表,从外到内的编码都丢到列表中了
  // 空列表表示identity编码,意思就是不编码
  // 编码列表常常被忽略
  // 另一种编码是chunked,发送和接收请求时会自动添加到列表或自动删除
  TransferEncoding []string

  // 是否是短连接
  // 对于服务端,标识了响应之后是否要关闭连接
  //    server对此字段会自动处理,handler无需处理
  // 对于客户端,标识了读响应之后是否关闭连接
  //    设置这个值是为了复用连接,前提是请求都指向同一个host
  //    类似Transport的keep alive
  Close bool

  // 对于服务端,
  //    Host表示了要寻找的URL,要么是头信息中的host,要么是url本身
  //    host可以是host:port模式,对于idn(国际域名)可能有两种风格
  //    punycode/unicode,不管哪种格式,golang.org/x/net/idna都可以按需转换
  //    为了防止dns重绑定攻击,handler需要校验头信息中的host是否有证明自己是权威的信息
  //    ServeMux支持注册具体的host名,用于保护注册的handler
  // 对于客户端,
  //    host是可选的,如果指定了,就在发送时覆盖头信息中的host
  //    如果为空,Request.Write方法就使用URL.Host
  Host string

  // 表单,包含了解析后的表单数据,要在ParseForm调用之后使用
  // 客户端会忽略此字段,而直接使用Body
  // 表单中包含拿了URL字段的查询参数,patch/post/put的表单数据
  Form url.Values

  // 解析后的表单数据,是patch/post/put的body参数解析出的表单数据
  // 客户端会忽略此字段,只能在ParseForm之后使用
  PostForm url.Values

  // 多表单数据,包括文件上传,只能在ParseMultipartForm调用之后使用
  // 客户端会忽略此字段,而直接使用Body
  MultipartForm *multipart.Form

  // 只有少数的客户端/服务端/代理才支持"预告"
  Trailer Header

  // 远端地址,客户端会忽略此字段
  // ReadRequest不会填充此字段,也没有固定的格式
  // 当前go实现,将远端地址设置为ip:port,且在handler处理之前设置
  // 这个字段的意义在于让server或其他软件能记录发送请求的地址
  // 实际上,这个远端地址,就是客户端的地址
  // 此字段常用于日志
  RemoteAddr string

  // 请求uri,是客户端发给服务端的原始uri
  // 一般会被URL代替
  // 客户端设置此字段会报错
  RequestURI string

  // tls允许server或其他软件记录tls连接的信息
  // ReadRequest不会填充此字段\
  // 调用handler之前,会填充连接信息,前提是tls启用了;否则置为nil
  // 和远程地址一样,客户端会忽略此字段
  TLS *tls.ConnectionState

  // 取消通道,可选,弃用阶段
  Cancel <-chan struct{}

  // 请求中的响应,适用于客户端重定向的场景
  // 这里的响应是引发当前请求的响应
  Response *Response

  // 只能在Request.WithContext修改ctx
  // 为了防止误用改变调用者的上下文,所以非暴露
  ctx context.Context
}

// 响应编辑器,在handler内部使用,用于构造一个响应
// Handler.ServeHTTP方法返回后,不得使用响应编辑器
type ResponseWriter interface {
  // 获取由WriteHeader发送的头信息
  // 在调用完WriteHandler/Write之后修改头信息是无效的
  // 除非用trailer去修改
  // 支持自动响应头,只需要将头里的值设为nil即可,eg:Date
  Header() Header

  // 将响应数据写入连接
  // 如果没有调用WriteHeader,Write会在写数据之前调用WriteHeader(200)
  // 头信息中如果不包含Content-Type,Write会自动添加[]Content-Type
  // 另外,如果响应数据才几KB,且没有调用Flush,就自动太难加Content-Length
  // 调用Write/WriteHeader之后不能再读Request.Body
  //    http1.x,主动调用Flusher.Flush或写数据太多触发Flush之后,请求体就不可用了
  Write([]byte) (int, error)

  // 写响应头,并设置响应码
  // 如果不显式调用WriteHeader,Write会自动调用,并设置响应码为200
  // 一般WriteHeader用于发送错误码
  // 有效的响应码是1xx-5xx,目前go实现不支持用户定义的1xx,除了100(继续)
  WriteHeader(statusCode int)
}

// 
type Flusher interface {
  Flush()
}
```
