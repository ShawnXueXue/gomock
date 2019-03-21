package main

import (
	"flag"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/pprof"
	"main/api"
)

func main() {
	adminPort := flag.String("adminport", "8180", "port for create/update/delete mock api")
	workingPort := flag.String("workingport", "8080", "port exposed to users")
	flag.Parse()
	house, e := api.NewApiHouse("shawn")
	if nil != e {
		panic(e)
	}
	defer house.Close()
	app := iris.Default()
	app.Post("admin/mock/add.json", func(ctx iris.Context) {
		var a api.Api
		ctx.ReadJSON(&a)
		house.AddApi(a)
	})
	app.Delete("admin/mock/del/{apiName}/{type}", func(ctx iris.Context) {
		house.DelApi(ctx.Params().Get("apiName"), ctx.Params().Get("type"))
	})
	app.Get("admin/mock/get", func(ctx iris.Context) {
		ctx.JSON(house.GetAllApi())
	})
	app.Get("/debug", func(ctx iris.Context) {
		ctx.HTML("<h1><a href='/debug/pprof/heap'>heap</a><br>" +
			"<a href='/debug/pprof/goroutine'>goroutine</a><br>" +
			"<a href='/debug/pprof/debug/block'>block</a><br>" +
			"<a href='/debug/pprof/debug/cmdline'>cmdline</a><br>" +
			"<a href='/debug/pprof/debug/profile'>profile</a><br>" +
			"<a href='/debug/pprof/debug/symbol'>symbol</a><br>" +
			"<a href='/debug/pprof/threadcreate'>threadcreate</a><br>")
	})
	app.Any("/debug/pprof/{action:path}", pprof.New())
	go app.Run(iris.Addr(":" + *adminPort))
	house.Start(*workingPort)
}
