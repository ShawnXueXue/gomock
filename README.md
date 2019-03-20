# gomock
mock server for rest api written by go

## build
1. Update go to at least version 1.12
2. Download code
3. Enable go modules
4. `go mod tidy`
5. 访问不了`golang.org/x`的朋友, 需要手动下载代码, 使用`go mod edit -replace`替换为本地代码.  
例如:
`go mod edit -replace=golang.org/x/net@v0.0.0-20180906233101-161cd47e91fd={local code path}`

## deploy
1. `go build`
2. `main.exe -adminport 8180 -workingport 8080`  
`adminport`:port for create/update/delete mock api  
`workingport`:port exposed to users

## usage
* create/update api. use `RequestType`/`ApiName` as its unique key.  
```
http://localhost:8180/admin/mock/add.json

{
	"ApiName":"1",
	"Response":{"key":"value"},
	"RequestType":"GET",
	"Status":200
}
```

* delete api
```
http://localhost:8180/admin/mock/del/{ApiName}/{RequestType}
```
* get api
```
http://localhost:8180/admin/mock/get
```

## todo
- status code not working
