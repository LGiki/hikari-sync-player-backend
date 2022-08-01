# ✨ Hikari Sync Player Backend

<img src="icon.svg" width="120" height="120" align="right" />

一个用来和朋友远程一起听播客~~、看视频~~的在线工具的后端。

该项目为个人行为，与小宇宙无关。

更多的说明、截图请查看[hikari-sync-player-frontend](https://github.com/LGiki/hikari-sync-player-frontend)的[README](https://github.com/LGiki/hikari-sync-player-frontend/blob/main/README.md)。

# Deploy

1. 安装 [Redis](https://redis.io/)

2. ```bash
   git clone https://github.com/LGiki/hikari-sync-player-backend.git
   cd hikari-sync-player-backend/
   go build
   ```

3. 编辑 [conf/app.toml](conf/app.toml) 文件，主要的配置项有

   - `[App]`
     - `Host`：后端程序监听的IP，例如：0.0.0.0
     - `Port`：后端程序监听的端口，例如：12321
   - `[Redis]`
     - `Host`：Redis服务的地址，例如：127.0.0.1:6379
     - `Password`：Redis密码，如果未设置密码请留空

4. ```bash
   ./hikari_sync_player
   ```

# License

MIT License