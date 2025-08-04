# SlothTracker

> **树懒追踪器**
>
> 一个设备状态查看器(~~摸鱼监控神器~~)

## API

[API文档](https://s8y25sfnie.apifox.cn)

使用的技术栈:

* 后端: GoLang + Gin + GORM
* 数据库: SQLite

## Desktop

使用的技术栈:

* GoLang + Wails
* Vite + Vue3

## 如何调试

拉取仓库后, 需要先安装依赖.

### 启动后端

```bash
cd API
go run main.go
```

### 启动桌面端

```bash
cd Desktop
wails dev
```

## 如何构建

### 网页端

```bash
cd Desktop/frontend
npm run build
```

### 桌面端

```bash
cd Desktop
wails build
```
