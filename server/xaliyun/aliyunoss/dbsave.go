package aliyunoss

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/77d88/go-kit/basic/xstr"
	"github.com/77d88/go-kit/plugins/xdatabase/xdb"
	"github.com/77d88/go-kit/plugins/xlog"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"gorm.io/gorm"
)

type OFile struct {
	Id       int64  `form:"id" json:"id,string"`
	Key      string `form:"key" json:"key,omitempty"`
	ETag     string `form:"etag" json:"etag,omitempty"`
	Type     int32  `form:"type" json:"type"`
	MaxWidth int64  `form:"maxWidth" json:"maxWidth"` // 图片最大宽度
}

type ObjFun func(c context.Context, key, toPath string) error

// CopyObjFun 拷贝对象
var CopyObjFun ObjFun = func(c context.Context, key, toPath string) error {
	_, err := Client.CopyObject(c, &oss.CopyObjectRequest{
		Bucket:       &Config.OssBucket,
		Key:          &toPath,
		SourceBucket: &Config.OssBucket,
		SourceKey:    &key,
	})
	return err
}

func DbSave(c context.Context, db *gorm.DB, r OFile, objFun ...ObjFun) (*OFile, error) {
	var glObjFun ObjFun

	if len(objFun) == 0 {
		// 默认使用拷贝对象
		glObjFun = CopyObjFun
	} else {
		glObjFun = objFun[0]
	}

	// 有id的情况 直接返回
	if r.Id > 0 {
		return &OFile{Id: r.Id}, nil
	} else {
		if xstr.IsBlank(r.Key) {
			xlog.Errorf(c, "参数错误 未填写key")
			return nil, errors.New("参数错误")
		}
		if xstr.IsBlank(r.ETag) {
			xlog.Errorf(c, "参数错误 未填写ETag")
			return nil, errors.New("参数错误")
		}
		o := Config

		var res Res
		result := db.WithContext(c).Where("ali_etag = ?", r.ETag).Limit(1).Find(&res)
		if result.Error != nil {
			return nil, result.Error
		}

		if res.ID > 0 {
			if !res.IsOptimize { // 是否优化过
				err := OptimizeRes(c, db, res, r)
				if err != nil {
					return nil, err
				}
			}
			return &OFile{Id: res.ID}, nil
		}

		err := db.WithContext(c).Transaction(func(tx *gorm.DB) error {
			base := xdb.NewBaseModel()
			res = Res{
				BaseModel:  base,
				RefTime:    time.Now(),
				AliEtag:    r.ETag,
				Path:       fmt.Sprintf("%s/%d", o.SavePrefix, base.ID),
				MimeType:   r.Type,
				IsOptimize: r.Type == ResTypeImage, // 图片资源默认优化
			}
			if result := tx.Create(&res); result.Error != nil {
				return result.Error
			}

			if err := glObjFun(c, r.Key, res.Path); err != nil {
				return err
			}

			if err := imgHandler(c, &res, r); err != nil {
				return err
			}

			return nil
		})
		if err != nil {
			return nil, err
		}
		return &OFile{Id: res.ID}, nil
	}
}

func OptimizeRes(c context.Context, db *gorm.DB, res Res, r OFile) error {
	return db.WithContext(c).Transaction(func(d *gorm.DB) error {
		result := d.Model(&Res{}).Where("id = ?", res.ID).Update("is_optimize", true)
		if result.Error != nil {
			return result.Error
		}
		err := imgHandler(c, &res, r)
		if err != nil {
			return err
		}
		return nil
	})
}

// imgHandler 处理图片
func imgHandler(ctx context.Context, res *Res, r OFile) error {
	// 如果是文件类型生成相关缩略图
	if r.Type == ResTypeImage {
		_, err := resetImageSize(ctx, res.Path, 2560, 100) // 仅保留 2560 的原图
		if err != nil {
			return err
		}
		_, err = genOtherImg(ctx, res.Path, 100, 30) // 最低质量底图
		if err != nil {
			return err
		}
		if r.MaxWidth >= 720 || r.MaxWidth == 0 {
			_, err = genOtherImg(ctx, res.Path, 720, 90) // 通用列表可以用
			if err != nil {
				return err
			}
		} else {
			if r.MaxWidth > 100 {
				_, err = genOtherImg(ctx, res.Path, r.MaxWidth, 90)
				if err != nil {
					return err
				}
			}
		}
		if r.MaxWidth >= 1280 || r.MaxWidth == 0 {
			_, err = genOtherImg(ctx, res.Path, 1280, 90) // 通用列表大图可以用
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// GenOtherImg 生成其他尺寸图片
func genOtherImg(ctx context.Context, key string, width, quality int64) (*oss.ProcessObjectResult, error) {
	targetImageName := fmt.Sprintf("%s_%d", key, width)
	// 将图片缩放为固定宽高100 px后转存到指定存储空间
	style := fmt.Sprintf("image/auto-orient,1/resize,m_lfit,w_%d/quality,q_%d/format,webp", width, quality)
	process := fmt.Sprintf("%s|sys/saveas,o_%v,b_%v", style,
		base64.URLEncoding.EncodeToString([]byte(targetImageName)),
		base64.URLEncoding.EncodeToString([]byte(Config.OssBucket)))

	// 构建一个ProcessObject请求，用于发起对特定对象的同步处理
	request := &oss.ProcessObjectRequest{
		Bucket:  oss.Ptr(Config.OssBucket), // 指定要操作的存储空间名称
		Key:     oss.Ptr(key),              // 指定要处理的图片名称
		Process: oss.Ptr(process),          // 指定处理指令
	}

	// 构建一个ProcessObject请求，用于发起对特定对象的同步处理
	return Client.ProcessObject(ctx, request)
}

func resetImageSize(ctx context.Context, key string, width, quality int64) (*oss.ProcessObjectResult, error) {
	// 将图片缩放为固定宽高100 px后转存到指定存储空间
	style := fmt.Sprintf("image/auto-orient,1/interlace,1/resize,m_lfit,w_%d/quality,q_%d/format,jpeg", width, quality)
	process := fmt.Sprintf("%s|sys/saveas,o_%v,b_%v", style,
		base64.URLEncoding.EncodeToString([]byte(key)),
		base64.URLEncoding.EncodeToString([]byte(Config.OssBucket)))

	// 构建一个ProcessObject请求，用于发起对特定对象的同步处理
	request := &oss.ProcessObjectRequest{
		Bucket:  oss.Ptr(Config.OssBucket), // 指定要操作的存储空间名称
		Key:     oss.Ptr(key),              // 指定要处理的图片名称
		Process: oss.Ptr(process),          // 指定处理指令
	}

	// 构建一个ProcessObject请求，用于发起对特定对象的同步处理
	return Client.ProcessObject(ctx, request)
}
