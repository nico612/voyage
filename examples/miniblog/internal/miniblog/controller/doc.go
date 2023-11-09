package controller

// Controller 层主要完成：接收 HTTP 请求，并进行参数解析、参数校验、逻辑分发处理、请求返回操作。
// 参数验证工具govalidator: https://github.com/asaskevich/govalidator
// 通过在 Controller 层实现有限的功能（参数解析、校验、逻辑分发、请求聚合和返回），并将负责的业务逻辑放在 Biz 层实现，可以使 Controller 层代码逻辑结构清晰，利于后期的代码维护
