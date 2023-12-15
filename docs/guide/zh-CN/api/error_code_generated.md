# 错误码

！！voyage 系统错误码列表，由 `codegen -type=int -doc` 命令生成，不要对此文件做任何更改。

## 功能说明

如果返回结果中存在 `code` 字段，则表示调用 API 接口失败。例如：

```json
{
  "code": 100101,
  "message": "Database errors"
}
```

上述返回中 `code` 表示错误码，`message` 表示该错误的具体信息。每个错误同时也对应一个 HTTP 状态码，比如上述错误码对应了 HTTP 状态码 500(Internal Server Error)。

## 错误码列表

voyage 系统支持的错误码列表如下：

| Identifier | Code | HTTP Code | Description |
| ---------- | ---- | --------- | ----------- |
| ErrUserNotFound | 110101 | 200 | User not found. |
| ErrUserAlreadyExist | 110102 | 200 | User already exist. |
| ErrFailedAuthentication | 110103 | 200 | username or password error |
| ErrCreateBaseMenu | 110201 | 200 | 创建菜单失败 |
| ErrSuccess | 100000 | 200 | OK. |
| ErrUnknown | 100001 | 500 | Internal server error. |
| ErrBind | 100002 | 400 | Error occurred while binding the request body to the struct. |
| ErrValidation | 100003 | 400 | Validation failed. |
| ErrTokenInvalid | 100004 | 401 | Token invalid. |
| ErrTokenExpired | 100005 | 200 | Token Expired. |
| ErrPageNotFound | 100006 | 404 | Page not found. |
| ErrGetCaptcha | 100007 | 200 | 图形验证码获取失败 |
| ErrInvalidCaptcha | 100008 | 200 | 无效的图形验证码 |
| ErrInsufficientPermissions | 100009 | 200 | Insufficient permissions |
| ErrDatabase | 100101 | 500 | Database error. |
| ErrEncrypt | 100201 | 401 | Error occurred while encrypting the sysuser password. |
| ErrSignatureInvalid | 100202 | 401 | Signature is invalid. |
| ErrExpired | 100203 | 401 | Token expired. |
| ErrInvalidAuthHeader | 100204 | 401 | Invalid authorization header. |
| ErrMissingHeader | 100205 | 401 | The `Authorization` header was empty. |
| ErrPasswordIncorrect | 100206 | 401 | Password was incorrect. |
| ErrPermissionDenied | 100207 | 403 | Permission denied. |
| ErrEncodingFailed | 100301 | 500 | Encoding failed due to an error with the data. |
| ErrDecodingFailed | 100302 | 500 | Decoding failed due to an error with the data. |
| ErrInvalidJSON | 100303 | 500 | Data is not valid JSON. |
| ErrEncodingJSON | 100304 | 500 | JSON data could not be encoded. |
| ErrDecodingJSON | 100305 | 500 | JSON data could not be decoded. |
| ErrInvalidYaml | 100306 | 500 | Data is not valid Yaml. |
| ErrEncodingYaml | 100307 | 500 | Yaml data could not be encoded. |
| ErrDecodingYaml | 100308 | 500 | Yaml data could not be decoded. |

