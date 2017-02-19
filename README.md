hi-config
---------

令人愉悦的自带UI的配置管理工具, 灵感来自于[confd](https://github.com/kelseyhightower/confd)

## 目标

- 容易部署,最简单可以单文件运行. (也可以使用配置文件,管理多个项目)
- 支持多种配置数据源
   - redis
   - vault
   - mysql, pg
   - consul, etcd, zookeeper

- 配置项加密安全存储
- 配置项变化后自动更新配置文件  
- 多项目配置, 支持项目配置key前缀  
- 默认自己识别`.env.tmpl`为模板, 生成目标文件`.env`  
- 配置项完整性检查
- 友好的web界面体验


## setup

- test

`go test $(go list ./... | grep -v /vendor/)`
