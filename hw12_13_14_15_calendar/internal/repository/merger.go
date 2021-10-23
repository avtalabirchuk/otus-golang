package repository

import (
	"reflect"
	"time"

	"github.com/imdario/mergo"
)

type timeTransformer struct {
}

func (t timeTransformer) Transformer(typ reflect.Type) func(dst, src reflect.Value) error {
	if typ == reflect.TypeOf(time.Time{}) {
		return func(dst, src reflect.Value) error {
			if dst.CanSet() {
				isZero := dst.MethodByName("IsZero")
				result := isZero.Call([]reflect.Value{})
				if result[0].Bool() {
					dst.Set(src)
				}
			}
			return nil
		}
	}
	return nil
}

func MergeEvents(dst *Event, src Event) error {
	return mergo.Merge(dst, src, mergo.WithOverride, mergo.WithTransformers(timeTransformer{}))
}
