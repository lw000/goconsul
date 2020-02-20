# grpc版本 v1.26.0
go mod edit -require=google.golang.org/grpc@v1.26.0
go get -u -x google.golang.org/grpc@v1.26.0


# 修改mod配置
go mod edit -require=github.com/marten-seemann/qtls@v0.4.1
# 下载对应版本
go get -u -x github.com/marten-seemann/qtls@v0.4.1
