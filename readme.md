## 云函数验证码

### 依赖

- 云函数平台使用腾讯云 SCF
- 验证码生成使用 https://github.com/liujiawm/gocaptcha

### 使用说明

- `SCF_BASE_URL` 为腾讯云函数相应 API 网关的访问路径，如 `https://service-dlxxjcx0-1xxx02252.gz.apigw.tencentcs.com/release/simple-captcha`

* 获取验证码

  ```
  POST SCF_BASE_URL HTTP/1.1
  Content-Type: application/json; charset=utf-8

  ```

  请求参数(body)：

  | 参数名 | 类型   | 必填 | 值    | 说明                      |
  | ------ | ------ | ---- | ----- | ------------------------- |
  | action | string | 是   | `new` | 获取验证码时传固定值`new` |

  响应参数(body)：

  | 参数名     | 类型   | 值            | 说明                                           |
  | ---------- | ------ | ------------- | ---------------------------------------------- |
  | code       | int    | `2000`,`2901` | 状态码,`2000` 成功，`2901` 失败                |
  | msg        | string |               | 提示信息                                       |
  | img        | string |               | 验证码图片的 base64 编码                       |
  | ciphertext | string |               | 验证码 hash 值，用于校验用户输入验证码是否正确 |

- 校验验证码

  ```
  POST SCF_BASE_URL HTTP/1.1
  Content-Type: application/json; charset=utf-8

  ```

  请求参数(body)：

  | 参数名         | 类型   | 必填 | 值      | 说明                               |
  | -------------- | ------ | ---- | ------- | ---------------------------------- |
  | action         | string | 是   | `check` | 校验验证码时传入固定值`check`      |
  | usercode       | string | 是   |         | 用户输入的验证码                   |
  | userciphertext | string | 是   |         | 获取验证码时获取的 ciphertext 参数 |

  响应参数(body)：

  | 参数名      | 类型   | 参数           | 说明                                                          |
  | ----------- | ------ | -------------- | ------------------------------------------------------------- |
  | code        | int    | `2000`,`2901`  | 状态码,`2000`成功，`2901` 失败                                |
  | msg         | string |                | 提示信息                                                      |
  | checkstatus | string | `1`,`-1` ,`-2` | 校验验证码状态：成功(`1`)；验证码错误(`-1`)；验证码过期(`-2`) |

### 部署

- [编译打包](https://cloud.tencent.com/document/product/583/18032#.E7.BC.96.E8.AF.91.E6.89.93.E5.8C.85)
