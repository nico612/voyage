# 错误码

！！voyage 系统错误码列表，由 `codegen -type=int -doc` 命令生成，不要对此文件做任何更改。

## 功能说明

如果返回结果中存在 `code` 字段，则表示调用 API 接口失败。例如：

```json
{
  "code": 100101,
  "message": "Database error"
}
```

上述返回中 `code` 表示错误码，`message` 表示该错误的具体信息。每个错误同时也对应一个 HTTP 状态码，比如上述错误码对应了 HTTP 状态码 500(Internal Server Error)。

## 错误码列表

voyage 系统支持的错误码列表如下：

| Identifier | Code | HTTP Code | Description |
| ---------- | ---- | --------- | ----------- |
| Success | 100001 | 200 | OK. |
| ErrUnknown | 100002 | 200 | 内部服务错误 |
| ErrBind | 100003 | 200 | 请求参数绑定错误 |
| ErrValidation | 100004 | 200 | 参数验证错误 |
| ErrTokenInvalid | 100005 | 200 | 无效的token |
| ErrDatabase | 100101 | 200 | Database error. |
| ErrEncrypt | 100201 | 200 | Error occurred while encrypting the user password. |
| ErrSignatureInvalid | 100202 | 200 | Signature is invalid. |
| ErrExpired | 100203 | 200 | Token expired. |
| ErrInvalidAuthHeader | 100204 | 200 | Invalid authorization header. |
| ErrMissingHeader | 100205 | 200 | The `Authorization` header was empty. |
| ErrPasswordIncorrect | 100206 | 200 | Password was incorrect. |
| ErrPermissionDenied | 100207 | 200 | Permission denied. |
| ErrEncodingFailed | 100301 | 200 | Encoding failed due to an error with the data. |
| ErrDecodingFailed | 100302 | 200 | Decoding failed due to an error with the data. |
| ErrInvalidJSON | 100303 | 200 | Data is not valid JSON. |
| ErrEncodingJSON | 100304 | 200 | JSON data could not be encoded. |
| ErrDecodingJSON | 100305 | 200 | JSON data could not be decoded. |
| ErrInvalidYaml | 100306 | 200 | Data is not valid Yaml. |
| ErrEncodingYaml | 100307 | 200 | Yaml data could not be encoded. |
| ErrDecodingYaml | 100308 | 200 | Yaml data could not be decoded. |

