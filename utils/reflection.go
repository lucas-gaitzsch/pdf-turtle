package utils

import (
	"reflect"
	"strconv"

	"github.com/rs/zerolog/log"
)

func ReflectDefaultValues[S any](s *S) *S {
	t := reflect.TypeOf(s).Elem()
	tv := reflect.ValueOf(s).Elem()

	count := t.NumField()

	for i := 0; i < count; i++ {
		f := t.Field(i)
		if defaultStr, ok := f.Tag.Lookup("default"); ok {
			var v reflect.Value
			if f.Type.Kind() == reflect.Pointer {
				v, ok = convertTo(defaultStr, f.Type.Elem())
				if !ok {
					continue
				}
			} else {
				v, ok = convertTo(defaultStr, f.Type)
				if !ok {
					continue
				}
			}

			fv := tv.FieldByName(f.Name)
			if fv.CanSet() && fv.IsZero() {
				if fv.Kind() == reflect.Pointer {
					p := reflect.New(f.Type.Elem())
					p.Elem().Set(v)

					fv.Set(p)
				} else {
					fv.Set(v)
				}
			}
		}
	}
	return s
}

func convertTo(str string, toType reflect.Type) (reflect.Value, bool) {
	switch toType.Kind() {
	case reflect.Int:
		if iv, err := strconv.Atoi(str); err == nil {
			return reflect.ValueOf(iv), true
		} else {
			log.Warn().Err(err).Str("toParse", str).Msg("cant parse default int")
			return reflect.Value{}, false
		}
	case reflect.Bool:
		if b, err := strconv.ParseBool(str); err == nil {
			return reflect.ValueOf(b), true
		} else {
			log.Warn().Err(err).Str("toParse", str).Msg("cant parse default bool")
			return reflect.Value{}, false
		}
	default:
		return reflect.ValueOf(str).Convert(toType), true
	}
}
