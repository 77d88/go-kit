package api_db

import (
	"github.com/77d88/go-kit/basic/xcore"
	"github.com/77d88/go-kit/plugins/xapi/server/xhs"
	"github.com/77d88/go-kit/plugins/xdb"
	"github.com/77d88/go-kit/plugins/xe"
	"github.com/77d88/go-kit/plugins/xlog"
	"github.com/77d88/go-kit/server/admin_pro/pro"
)

type dbReq struct {
	DbName    string `json:"dbName"`
	ModelPath string `json:"modelPath"`
	Restore   bool   `json:"restore"`
}

func autoMigrate(ctx *xhs.Ctx) {
	var req dbReq
	ctx.ShouldBind(&req)

	db, err := xdb.GetDB(req.DbName)
	ctx.Fatalf(err)
	i := make([]interface{}, 0, len(xdb.RegisterModels))
	for k, v := range xdb.RegisterModels {
		if req.ModelPath == "" {
			i = append(i, v)
		} else {
			if k == req.ModelPath {
				i = append(i, v)
			}
		}
	}
	xlog.Tracef(nil, "auto migrate model %v", i)
	ctx.Fatalf(db.AutoMigrate(i...))
	if err != nil {
		xlog.Errorf(nil, "自动迁移失败: %+v", err)
		ctx.Fatalf(err)
	}
}

func initDataRows(ctx *xhs.Ctx) {
	var req dbReq
	ctx.ShouldBind(&req)

	for k, dist := range xdb.RegisterModels {
		data := dist.InitData()
		for _, v := range data {
			if v.GetID() <= 0 {
				continue
			}
			// 创建一个v实例
			d := xcore.NewBy(v)
			result := xdb.Ctx(ctx, req.DbName).First(&d, v.GetID())
			if result.IsNotFound() { // 记录不存在，创建新记录
				err := xdb.Ctx(ctx, req.DbName).Create(v)
				ctx.Fatalf(err)
				xlog.Infof(nil, "创建记录[%s]成功: %v", k, v.GetID())
			} else {
				// 删除在创建
				if req.Restore { // 重建建立所有数据
					ctx.Fatalf(xdb.Ctx(ctx, req.DbName).DeleteUnscoped(v))
					ctx.Fatalf(xdb.Ctx(ctx, req.DbName).Create(v))
				}
			}
		}
	}
}

func Register(api *xe.Engine, path string) {
	api.RegisterPost(path+"/auto_migrate", pro.SuperAdmin, autoMigrate)
	api.RegisterPost(path+"/init_data_rows", pro.SuperAdmin, initDataRows)
}
