# photouch

​		一款图片管理分享系统, 范围可大到面向所有人, 小到面向一个家族一个班级, 分享生活中的图片, 分享生活中的感动, 记录生活中的故事.



## linux上 安装

1. 下载photouch文件后,给予可执行权限
2. mysql数据库中创建数据库, 导入create.sql

3. 注意要创建config.json文件,内部存关键信息,

```json
{
    "email_user": "@qq.com",
    "email_key": "...",
    "mysql_user": "...",
    "mysql_key": "...",
    "mysql_db": "...",
    "token": "..."
}
```

"email_user"和"email_key"为注册发邮件所需的邮箱账号和stmp密码

"mysql_user", "mysql_key"和"mysql_db"为Linux上mysql服务器的账户密码和数据库名

token为暂时图片存储的图床的token, 暂时使用[去不图床](https://7bu.top/)

设置好后执行即可,默认端口:8000
