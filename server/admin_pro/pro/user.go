package pro

import (
	"context"
	"github.com/77d88/go-kit/basic/xarray"
	"github.com/77d88/go-kit/basic/xencrypt/xmd5"
	"github.com/77d88/go-kit/plugins/xdb"
	"github.com/77d88/go-kit/plugins/xlog"
)

const TableNameUser = "s_sys_user"

var (
	UserPwdSalt = "$$SS$#Y#X#@R@@!!!" //  密码盐
)

func init() {
	xdb.AddModels(&User{})
}

// User 用户表
type User struct {
	xdb.BaseModel
	UpdateUser          int64          `gorm:"comment:更新人"`
	Password            string         `gorm:"comment:后台登录密码" json:"password"`   // 后台登录密码
	Disabled            bool           `gorm:"comment:是否禁用" json:"disabled"`     // 是否禁用
	Username            string         `gorm:"comment:后台登录名称" json:"username"`   // 后台登录名称
	Nickname            string         `gorm:"comment:后台显示名称" json:"nickname"`   // 后台显示名称
	Avatar              *xdb.Int8Array `gorm:"comment:头像" json:"avatar"`         // 头像
	Roles               *xdb.Int8Array `gorm:"comment:系统角色" json:"roles"`        // 系统角色
	Permission          *xdb.Int8Array `gorm:"comment:系统独立权限" json:"permission"` // 系统独立权限
	Email               string         `gorm:"comment:邮箱" json:"email"`
	IsReLogin           bool           `gorm:"comment:是否需要重新登录" json:"isReLogin"`
	ReLoginDesc         string         `gorm:"comment:重新登录描述" json:"reLoginDesc"`
	PermissionCodes     *xdb.TextArray `gorm:"comment:权限码" json:"permissionCodes"`     // 冗余 集合Permission里面的所有
	RolePermissionCodes *xdb.TextArray `gorm:"comment:角色码" json:"RolePermissionCodes"` // 冗余 集合Roles里面的所有Permission Code
	IsSuperAdmin        bool           `gorm:"comment:是否是超级管理员" json:"isSuperAdmin"`
	_codes              []string       `gorm:"-"` // 本地计算的code
	_isCalcCodes        bool           `gorm:"-"` // 是否计算
}

// TableName Res's table name
func (*User) TableName() string {
	return TableNameUser
}

func (d *User) ToResponse() *UserDst {
	return &UserDst{
		Id:          d.ID,
		Disabled:    d.Disabled,
		Nickname:    d.Nickname,
		Username:    d.Username,
		Avatar:      d.Avatar,
		Roles:       d.Roles,
		Permission:  d.Permission,
		IsReLogin:   d.IsReLogin,
		ReLoginDesc: d.ReLoginDesc,
	}
}

// AllPermissionCode 获取所有权限码 去重
func (d *User) AllPermissionCode() []string {
	if d._isCalcCodes {
		return d._codes
	}
	codes := make([]string, 0)
	codes = append(codes, d.PermissionCodes.ToSlice()...)
	codes = append(codes, d.RolePermissionCodes.ToSlice()...)
	codes = xarray.Unique(codes)
	if d.IsSuperAdmin {
		codes = append(codes, RoleSuperAdmin)
	}
	d._codes = codes
	d._isCalcCodes = true
	return codes
}

type UserDst struct {
	Id          int64          `json:"id,string"`
	Password    string         `json:"password,omitempty"`
	Disabled    bool           `json:"disabled"`
	Username    string         `json:"username"`
	Nickname    string         `json:"nickname"`
	Avatar      *xdb.Int8Array `json:"avatar"`
	Roles       *xdb.Int8Array `json:"roles"`
	Permission  *xdb.Int8Array `json:"permission"`
	Email       string         `json:"email"`
	IsReLogin   bool           `json:"isReLogin"`
	ReLoginDesc string         `json:"reLoginDesc"`
}

func ToUserResponses(dicts []*User) []*UserDst {
	res := make([]*UserDst, 0)
	if len(dicts) == 0 {
		return res
	}

	for _, dict := range dicts {
		res = append(res, dict.ToResponse())
	}
	return res
}

func (*User) InitData() []xdb.GromModel {
	return []xdb.GromModel{
		// 初始化一个普通管理员
		&User{
			BaseModel: xdb.NewBaseModel(1),
			Username:  "admin",
			Password:  xmd5.EncryptSalt("123456", UserPwdSalt),
			Nickname:  "管理员",
			Disabled:  false,
			Email:     "admin@admin.com",
		},
	}
}

func (d *User) HasPermission(code ...string) bool {
	return xarray.ContainAny(d.AllPermissionCode(), code)
}

// GetUserAllPermissionCode 获取用户所有权限
func GetUserAllPermissionCode(c context.Context, userId int64) ([]string, error) {
	var user User
	result := xdb.Ctx(c).Model(&User{}).WithId(userId).First(&user)
	if result.Error != nil {
		xlog.Errorf(c, "获取用户%d 异常", userId)
		return nil, result.Error
	}
	return user.AllPermissionCode(), nil
}
