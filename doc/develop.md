### 平台开发说明文档

1. 每一个平台应继承上述结构体

在开发新的平台时，应在`platforms`文件夹内新建平台文件如`MyPlatform.go`，设置包名为`platforms`, 新平台应继承包`platforms`中的`Model`结构体，以便获得基本的平台功能和属性。

```go
package platforms

type MyPlatform struct {
    *Model
    // 自定义属性和方法
}
```

继承`Model`结构体后，可以使用`Model`结构体中定义的属性和方法，如`Article`、`Cookies`、`Ctx`等。

2. 平台应将需要重写的方法自行重写以实现相关功能

在继承`Model`结构体后，需要根据具体平台的特性和需求，重写以下方法以实现相关功能：

- `SetConfig(foreDefault bool) (err error)`: 设置配置 (需要重写)
- `CheckAuthentication() (authInfo map[string]string, err error)`: 检查授权 (需要重写)
- `Login() (err error)`: 登录 (需要重写)
- `Run() (err error)`: 平台运行方法 (需要重写)

通过重写以上方法，你可以根据自己的需求定制平台的行为，以实现特定的功能。

示例代码：

```go
package platforms


type MyPlatform struct {
    *Model
    // 自定义属性和方法
}

// 重写SetConfig方法
func (m *MyPlatform) SetConfig(foreDefault bool) (err error) {
    // 实现加载配置的逻辑
}

// 重写CheckAuthentication方法
func (m *MyPlatform) CheckAuthentication() (authInfo map[string]string, err error) {
    // 实现检查授权的逻辑
}

// 重写Login方法
func (m *MyPlatform) Login() (err error) {
    // 实现登录的逻辑
}

// 重写Run方法
func (m *MyPlatform) Run() (err error) {
    // 实现平台的运行逻辑
}
```

以上示例代码展示了如何创建一个名为`MyPlatform`的新平台，并重写了`SetConfig`、`CheckAuthentication`、`Login`和`Run`方法，以实现自定义的平台功能。

通过继承`Model`结构体并重写相关方法，你可以开发出适合特定需求的平台，实现各种功能，如加载配置、检查授权、登录和平台运行。