# API 文档

| 方法 | 路径              | 参数                    | 用途                                         | 响应                                  | 是否需要身份认证 |
| ---- | ----------------- | ----------------------- | -------------------------------------------- | ------------------------------------- | ---------------- |
| GET  | "/"               | 无                      | 主页面                                       | HTML                                  | 0                |
| GET  | "/login"          | 无                      | 登录页面                                     | HTML                                  | 0                |
| GET  | "/signup"         | 无                      | 注册页面                                     | HTML                                  | 0                |
| GET  | "/wj"             | 无                      | 忘记密码页面                                 | HTML                                  | 0                |
| POST | "/user/login"     | name  && pwd            | 登录表单,设定cookie:user:~~                  | 200 : string: "200" ,201:string:"400" | 0                |
| POST | "/signup/email"   | email                   | 注册页面                                     | 200 : string: "YES"                   | 0                |
| POST | "/signup/up"      | code,pwd和emailname     | 注册页面                                     | 200 : string: "YES" 或string:"NO"     | 0                |
| GET  | "/download"       | url                     | 下载                                         | 图片文件                              | 0                |
| GET  | "/img/rand"       | 无                      | 获取32张随机公有图片的缩略图                 | JSON:URL......                        | 0                |
| POST | "/img/big"        | minurl                  | 获取图片的大图URL                            | string: URL 大图地址                  | 0                |
| GET  | "/img/user/like"  | 无                      | 个人点赞过的图片的缩略图及其信息             | jSON:URL......                        | 1                |
| GET  | "/img/user/all"   | 无                      | 个人所有图片缩略图及其信息                   | JSON:URL......                        | 1                |
| GET  | "/img/user/share" | 无                      | 个人分享的图片                               | JSON:URL......                        | 1                |
| GET  | "/img/user/time"  | 无                      | 获取用户某年某月上传的个人图片的第一张缩略图 | JSON:year && month && url             | 1                |
| POST | "/img/user/month" | year && month           | 获取用户某年某月上传的个人图片的所有缩略图   | JSON:url...信息                       | 1                |
| POST | "/upload"         | "files[]"做分隔符的表单 | 上传                                         | 200 : string: "YES" ,201:string:"NO"  | 1                |
| POST | "/change"         | minurl && old  && new   | 更改私有公有权限                             | 200 : string: "YES"                   | 1                |
| POST | "/good"           | minurl                  | 点赞                                         | 200 : string: "YES"                   | 1                |

