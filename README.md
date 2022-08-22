# gopixiu-eventer
`gopixiu-eventer` 旨在对 kubernetes 原生功能的补充和强化
- 提供一个控制器收集 `Event` 信息到 `ES` 中实现 `Event` 数据持久化
  - 可通过 `go run main.go -h` 查看参数
```
go run main.go -h                                                              
Usage of /var/folders/h7/dnghff2d3f11mkbsrnqp75880000gn/T/go-build3024363050/b001/exe/main:
  -elasticsearch.address string
        (可选) elasticsearch address 地址
  -elasticsearch.password string
        (可选) elasticsearch 用户名
  -elasticsearch.username string
        (可选) elasticsearch 用户名
  -kubernetes.kubeConfig string
        (可选) kubeconfig 文件的绝对路径 (default "/Users/jimingyu/.kube/config")

```

- `ES` 信息查询
```
GET http://localhost:9200/event/_search?pretty
```