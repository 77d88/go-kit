package xreflect

import (
	"fmt"
	"reflect"
)

type RefVal struct {
	s           interface{}
	InstValue   reflect.Value // 如果是指针，则Inst为指向值
	InstKind    reflect.Kind  // 实例类型
	SourceValue reflect.Value // 实例
	SourceKind  reflect.Kind
}

func Warp(val interface{}) *RefVal {
	if val == nil {
		return &RefVal{}
	}
	inst, kind := GetInst(val)
	of := reflect.ValueOf(val)
	return &RefVal{
		SourceValue: of,
		SourceKind:  of.Kind(),
		s:           val,
		InstValue:   inst,
		InstKind:    kind,
	}
}

func (v *RefVal) Is(kind reflect.Kind) bool {
	if v.s == nil {
		return false
	}
	if v.SourceKind == kind {
		return true
	}
	return v.InstKind == kind
}

func (v *RefVal) StructCallInterface(ifacePtr interface{}, methodName string, args ...interface{}) ([]interface{}, error) {
	return CallInterfaceMethod(v.s, ifacePtr, methodName, args...)
}

func (v *RefVal) InstPath() string {
	t := v.InstValue.Type()
	return fmt.Sprintf("%s:%s", t.PkgPath(), t.Name())
}
