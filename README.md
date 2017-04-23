EConfig
---------

令人愉悦的自带UI的配置管理工具, 灵感来自于[confd](https://github.com/kelseyhightower/confd)

## 目标(Features)

- 容易部署,最简单可以单文件运行. (也可以使用配置文件,管理多个项目)
- 支持多种配置数据源
   - redis
   - mysql, pg
   - vault (TODO)
   - consul, etcd, zookeeper (TODO)

- 配置项变化后自动更新配置文件  
- 多项目配置, 支持项目配置key前缀  
- 默认自己识别`.env.tmpl`为模板, 生成目标文件`.env`  
- 配置项完整性检查
- 友好的web界面体验
- 配置项加密安全存储 (TODO)


## 部署(Setup)


## 运行(Run)

- 以服务的方式运行(Run as serve)

`econfig serve -v --config /etc/econfig/.econfig.toml --backend postgres --http-port 1520`

在浏览器中访问(Open it in browser) `http://localhost:1520`

- 立即执行(Run only once)

`econfig sync -v --config /etc/econfig/.econfig.toml --backend postgres --app myapp`

## Thanks

[iris](https://github.com/kataras/iris)
[vue](https://vuejs.org)
[iviewui](https://www.iviewui.com)
