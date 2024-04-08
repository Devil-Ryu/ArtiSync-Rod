## 平台模型 (Package: platforms)

该模块定义了平台模型的结构和方法。

### Model 结构体

```go
type Model struct {
    Article       *utils.Article              // 待上传文章
    Cookies       []*proto.NetworkCookieParam // cookies
    Ctx           context.Context             // 上下文
    Name          string                      // 名称
    PlatformIndex int                         // 平台序号
    RODController *controller.RODController   // 机器人控制器
}
```

该结构体包含以下字段:
- `Article`: 待上传的文章
- `Cookies`: cookies数组
- `Ctx`: 上下文
- `Name`: 平台名称
- `PlatformIndex`: 平台序号
- `RODController`: 机器人控制器

### 方法

#### Init

初始化平台。

```go
func (m *Model) Init(ctx context.Context, article *utils.Article, platformIndex int)
```

参数:
- `ctx`: 上下文对象
- `article`: 待上传的文章
- `platformIndex`: 平台序号

#### SetCookies

设置cookies。

```go
func (m *Model) SetCookies(Cookies []*proto.NetworkCookieParam)
```

参数:
- `Cookies`: cookies数组

#### LoadCookies

读取cookies。

```go
func (m *Model) LoadCookies() (cookies []*proto.NetworkCookieParam, err error)
```

返回值:
- `cookies`: cookies数组
- `err`: 错误信息

#### LoadConfig

加载配置。

```go
func (m *Model) LoadConfig(defaultConfig map[string]interface{}, forceDefault bool) (config map[string]interface{}, err error)
```

参数:
- `defaultConfig`: 默认配置
- `forceDefault`: 是否强制使用默认配置

返回值:
- `config`: 配置信息
- `err`: 错误信息

#### CheckConfig

检查配置。

```go
func (m *Model) CheckConfig(config interface{}) (err error)
```

参数:
- `config`: 配置信息

返回值:
- `err`: 错误信息

#### OpenPage

打开页面。

```go
func (m *Model) OpenPage(pageURL string) (err error)
```

参数:
- `pageURL`: 页面URL

返回值:
- `err`: 错误信息

#### SetConfig

加载配置 (需要重写)。

```go
func (m *Model) SetConfig(foreDefault bool) (err error)
```

返回值:
- `err`: 错误信息

#### CheckAuthentication

检查授权 (需要重写)。

```go
func (m *Model) CheckAuthentication() (authInfo map[string]string, err error)
```

返回值:
- `authInfo`: 授权信息
- `err`: 错误信息

#### Login

登录 (需要重写)。

```go
func (m *Model) Login() (err error)
```

返回值:
- `err`: 错误信息

#### Run

平台运行方法 (需要重写)。

```go
func (m *Model) Run() (err error)
```

返回值:
- `err`: 错误信息