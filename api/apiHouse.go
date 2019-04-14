package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"reflect"
	"time"
)

const (
	dbName = "ApiHouse.db"
)

type ApiHouse struct {
	db      *bolt.DB
	bucket  string
	aServer *iris.Application
}

func NewApiHouse(bucket string) (*ApiHouse, error) {
	ah := ApiHouse{bucket: bucket}
	db, e := bolt.Open(dbName, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if e != nil {
		return nil, e
	}
	db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte(ah.bucket))
		return nil
	})
	ah.db = db
	ah.aServer = iris.Default()
	return &ah, nil
}

func (ah *ApiHouse) Close() error {
	e := ah.db.Close()
	if nil != e {
		return e
	}
	return nil
}

func (ah *ApiHouse) GetAllApi() []Api {
	var apiList []Api
	ah.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(ah.bucket))
		apiList = make([]Api, 0, b.Stats().KeyN)
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			line := Api{}
			json.Unmarshal(v, &line)
			apiList = append(apiList, line)
		}
		return nil
	})
	return apiList
}

func (ah *ApiHouse) AddApi(api Api) error {
	if api.ApiName == "" {
		return errors.New("can not add api with out api name")
	}
	apiId := api.getId()
	_ = ah.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(ah.bucket))
		if nil != err {
			return err
		}
		j, _ := json.Marshal(api)
		b.Put([]byte(apiId), j)
		return nil
	})
	route := ah.aServer.GetRoute(apiId)
	if nil != route {
		if iris.MethodNone == route.Method {
			route.RestoreStatus()
		}
		route.Handlers[len(route.Handlers)-1] = handlerFunc(api)
	} else {
		ah.aServer.Handle(api.RequestType, api.ApiName, record, handlerFunc(api))
	}
	ah.aServer.RefreshRouter()
	return nil
}

func (ah *ApiHouse) DelApi(apiName, requestType string) error {
	if apiName == "" || requestType == "" {
		return errors.New("can not add api with out api name")
	}
	api, _ := ah.getApi(apiName, requestType)
	if nil != api {
		apiId := api.getId()
		ah.db.Update(func(tx *bolt.Tx) error {
			return tx.Bucket([]byte(ah.bucket)).Delete([]byte(apiId))
		})
		route := ah.aServer.GetRoute(apiId)
		route.SetStatusOffline()
		ah.aServer.RefreshRouter()
	}
	return nil
}

func (ah *ApiHouse) getApi(apiName, requestType string) (*Api, error) {
	var b []byte
	ah.db.View(func(tx *bolt.Tx) error {
		b = tx.Bucket([]byte(ah.bucket)).Get([]byte(getId(apiName, requestType)))
		return nil
	})
	if nil == b {
		return nil, fmt.Errorf("can not get api: %s", apiName)
	}
	a := &Api{}
	json.Unmarshal(b, a)
	return a, nil
}

func (ah *ApiHouse) Start(port string) {
	ah.loadAll()
	ah.aServer.Run(iris.Addr(":" + port))
}

func (ah *ApiHouse) loadAll() {
	for _, api := range ah.GetAllApi() {
		ah.aServer.Handle(api.RequestType, api.ApiName, record, handlerFunc(api))
	}
}

func record(ctx context.Context) {
	ctx.Application().Logger().Infof("(%s) Handler is executing from: '%s'", ctx.Path(), reflect.TypeOf(ctx).Elem().Name())
	ctx.Next()
}
func handlerFunc(api Api) context.Handler {
	return func(ctx iris.Context) {
		// should before .Write
		ctx.StatusCode(api.Status)
		// .Write
		ctx.JSON(api.Response)
	}
}
