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

# 🧭 项目整体流程规划

## 第一阶段: 基础服务端搭建(Go)

✅ 目标: 实现注册/状态上传/状态查看 API, 支持唯一 ID 识别终端

## 功能点

* ✅ 注册终端(首次请求, 返回唯一 ID)
* ✅ 上传终端状态(附带 ID)
* ✅ 查询终端状态(公开接口, 可嵌入前端)
* ✅ SQLite 数据存储(简单, 便携)
* ✅ 基础结构清晰, 便于后续拓展(用户绑定/Token/标签等)

## 第二阶段: Windows 客户端(Go)

✅ 目标: 定时采集状态 + 配置存储 + 自动注册 + 上传状态

### 功能点

* ✅ 第一次启动时自动注册并保存 ID
* ✅ 每 60 秒采集一次:
  * 当前窗口标题
  * WiFi 状态(连网与否)
  * 电池状态(充电/电量)
  * 等
* ✅ 使用本地配置文件(JSON)
* ✅ CLI 输出调试信息

## 第三阶段: 前端展示(网页)

✅ 目标: 支持网页查看终端状态, 表格形式, 排序高亮在线状态

### 功能点

* ✅ Vue/React 项目展示终端状态表
* ✅ 实时刷新(轮询或 WebSocket)
* ✅ 在线状态高亮/状态图标美化

## 第四阶段: 安卓客户端(Android Studio)

# 📐 后端 API v1 设计

**通用说明**

* 所有请求使用 JSON 格式
* 所有状态请求需附带唯一 deviceId 字段(除注册接口)

## 1. `POST /api/register`

**说明**: 首次运行客户端自动调用, 返回唯一设备 ID

**请求示例**:

```json
{
  "deviceName": "洱海的电脑",
  "platform": "windows",
  "description": "主力开发机"
}
```

**响应示例**:

```json
{
  "code": 0,
  "deviceId": "a6e4b1c8-xxxx-xxxx-xxxx-71d112de13f2"
}
```

## 2. `POST /api/status/update`

**说明**: 客户端上报状态

**请求示例**:

```json
{
  "deviceId": "a6e4b1c8-xxxx-xxxx-xxxx-71d112de13f2",
  "timestamp": 1721744460,
  "battery": {
    "charging": true,
    "level": 87
  },
  "wifi": {
    "connected": true,
    "ssid": "MyHomeWiFi"
  },
  "foregroundApp": {
    "name": "Chrome",
    "title": "ChatGPT"
  },
  "screenOn": true
}
```

**响应示例**:

```json
{
  "code": 0,
  "message": "状态更新成功"
}
```

## 3. `GET /api/status/list`

**说明**: 前端拉取所有设备当前状态

**响应示例**:

```json
{
  "devices": [
    {
      "deviceId": "a6e4b1c8-xxxx-xxxx-xxxx-71d112de13f2",
      "name": "洱海的电脑",
      "platform": "windows",
      "lastUpdate": 1721744460,
      "battery": {
        "charging": true,
        "level": 87
      },
      "wifi": {
        "connected": true,
        "ssid": "MyHomeWiFi"
      },
      "foregroundApp": {
        "name": "Chrome",
        "title": "ChatGPT"
      },
      "screenOn": true
    }
  ]
}
```

# 🧩 数据表结构简化(SQLite)

```sql
-- 设备表
CREATE TABLE devices (
  id TEXT PRIMARY KEY,
  name TEXT,
  platform TEXT,
  description TEXT,
  registered_at TIMESTAMP
);

-- 状态表
CREATE TABLE device_status (
  device_id TEXT,
  timestamp INTEGER,
  battery_charging BOOLEAN,
  battery_level INTEGER,
  wifi_connected BOOLEAN,
  wifi_ssid TEXT,
  app_name TEXT,
  app_title TEXT,
  screen_on BOOLEAN,
  FOREIGN KEY(device_id) REFERENCES devices(id)
);
```
