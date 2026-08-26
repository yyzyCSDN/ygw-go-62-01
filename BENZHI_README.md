# ColumnStore

ColumnStore 是一个自包含的列式分析存储引擎。数据按列组织成不可变块，
批量加载写入后由向量化扫描和谓词下推处理分析查询，分区按时间裁剪，
低基数列支持字典压缩，过期块由回收任务清理。系统进程内维护全部状态，
不依赖外部数据库。

## 构建

依赖已 vendor 到 `vendor/`，构建完全离线：

```sh
go build -mod=vendor ./...
```

## 运行

```sh
go run -mod=vendor ./cmd/columnstore -listen 127.0.0.1:18080 -data .columnstore-data
```

启动后访问 http://127.0.0.1:18080/ 打开查询控制台。

## 容器镜像

```sh
sh build_benzhi_docker.sh
docker run -p 18080:18080 columnstore:latest
```

镜像使用离线 vendor 构建（`GOPROXY=off`），容器内同时保留 Go 工具链，
便于在容器内执行构建与检查命令。

## HTTP 接口

- `GET /` 查询控制台页面
- `GET /api/health` 服务健康状态
- `GET /api/query?table=&column=&from=&to=` 分析查询
- `POST /api/load` 批量加载数据
- `POST /api/compress` 触发字典压缩
- `POST /api/gc` 执行过期回收
