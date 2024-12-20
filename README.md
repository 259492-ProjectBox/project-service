# set environment for production

$env:ENV="prod"

% this is how to clear env
Remove-Item Env:ENV

เอาชื่อ project มาใส่ในเป็น file path ที่ อัพ

# Use Swagger

1.Install swag for generating Swagger documentation.
`go install github.com/swaggo/swag/cmd/swag@latest`

2.Install gin-swagger and swag packages in your project

```
go get -u github.com/swaggo/gin-swagger
go get -u github.com/swaggo/gin-swagger/swaggerFiles
```

3.Install swag Go annotations.

```
go get -u github.com/swaggo/files
```

4.Gen swagger using powershell need to do this every time api change

```
swag init
```

# CORS

```
go get github.com/gin-contrib/cors
```

wire.go

```
go install github.com/google/wire/cmd/wire@latest
```

# Minio

`go get github.com/minio/minio-go/v7`

if you can not run project you need to change the env in minioendpoint

# Gorm

gorm.io/gorm

# Hot reload air

https://github.com/air-verse/air
go install github.com/air-verse/air@latest

# PDF To Text

1. Install Scoop Cli
2. `scoop install poppler`

air init

air
