package v1

//解释下为什么会将 CreateUserRequest 结构体的定义文件 user.go 存放在 pkg/api/miniblog/v1 目录下：
//
//CreateUserRequest 对用户暴露，作为 POST /v1/users 接口的请求 Body，所以我将 user.go 放在了 pkg/ 目录下；
//
//另外，为了能让后来的代码维护者或者包的使用者感知到 CreateUserRequest 是专门用来做请求参数的结构体，我将 user.go 放在了 pkg/api 目录下；
//
//另外，考虑到未来 miniblog 可能会加入多个服务，每个服务都有自己的对外 API 定义，这里我新建了一个目录 miniblog，将 user.go 存放在了 pkg/api/miniblog 目录下；
//
//最后，考虑到未来 API 接口的版本升级和维护，这里我们也创建了一层 v1 目录用来保存不同版本包的请求参数结构体。所以最终，user.go 存放在了 pkg/api/miniblog/v1/ 目录下。
//
//可以看到，在开发过程中，我们需要时刻保持一个功能的可扩展性、可维护性。
