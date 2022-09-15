# nico

对外网关，负责将请求经过解包，寻址找到对应的微服务进行处理

### 支持的寻址方式

> 以下列表存在优先级，如果找到即使用该方式，并且跳过后续的寻址方式

- header寻址：`header: X-Fc-Service: <server name>`
- 路径寻址：`/<server name>/<path>`
