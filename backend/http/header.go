// Copyright (c) 2025 The Consee Authors. All rights reserved.
// SPDX-License-Identifier: MulanPSL-2.0

package http2

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var timeType = reflect.TypeOf(time.Time{})
var durationType = reflect.TypeOf(time.Duration(0))

var baseTypeParsers = map[reflect.Kind]func(string, string) (any, error){
	reflect.Int:    parseInt,
	reflect.Int8:   parseInt8,
	reflect.Int16:  parseInt16,
	reflect.Int32:  parseInt32,
	reflect.Int64:  parseInt64,
	reflect.Uint:   parseUint,
	reflect.Uint8:  parseUint8,
	reflect.Uint16: parseUint16,
	reflect.Uint32: parseUint32,
	reflect.Uint64: parseUint64,
	reflect.Bool:   parseBool,
	reflect.String: parseString,
}

var (
	errInvalidTarget = errors.New("target is not a pointer to a struct")
)

var (
	errFInvalidDurationFormat = "invalid duration format: %s"
)

func ParseHeader(h http.Header, target any) error {
	// 判断target是否为指针
	t1 := reflect.TypeOf(target)
	if t1.Kind() != reflect.Ptr {
		return errInvalidTarget
	}
	// 判断target的Elem是否为结构体
	t2 := t1.Elem()
	if t2.Kind() != reflect.Struct {
		return errInvalidTarget
	}
	v := reflect.ValueOf(target).Elem()
	// 找到 header 的 tag，解析
	for i, nf := 0, t2.NumField(); i < nf; i++ {
		f := t2.Field(i)
		if f.Anonymous {
			continue
		}
		tag := f.Tag.Get("header")
		if tag == "-" {
			continue
		}
		ss := strings.Split(tag, ",")
		var key string = ss[0]
		var format string
		if len(ss) > 1 {
			format = ss[1]
		}
		value := h.Values(key)
		if len(value) == 0 {
			continue
		}

		if f.Type.Kind() == reflect.Slice && f.Type.Elem().Kind() == reflect.String {
			// 字段是个 string[] 类型，直接 copy 后赋值
			vCopy := make([]string, len(value))
			copy(vCopy, value)
			v.Field(i).Set(reflect.ValueOf(vCopy))
			continue
		}

		v0 := value[0]
		if f.Type == timeType {
			if format == "" {
				format = time.DateTime
			}
			t, err := time.ParseInLocation(format, v0, time.Local)
			if err != nil {
				return err
			}
			v.Field(i).Set(reflect.ValueOf(t))
			continue
		}

		if f.Type == durationType {
			var d time.Duration
			if i64, err := strconv.ParseInt(v0, 10, 64); err == nil {
				// 目前很多人将整型直接用来表示时段，时段单位则由前后端直接约定，这里需要兼容
				// 用 format 指定时段单位，默认为秒
				times := time.Second
				switch format {
				case "ns":
					times = time.Nanosecond
				case "us", "µs":
					times = time.Microsecond
				case "ms":
					times = time.Millisecond
				case "", "s": // 默认秒
				case "m", "min":
					times = time.Minute
				case "h":
					times = time.Hour
				default:
					return fmt.Errorf(errFInvalidDurationFormat, format)
				}
				d = time.Duration(i64 * int64(times))
			} else {
				// 不是整型，直接尝试解析
				d, err = time.ParseDuration(v0)
				if err != nil {
					return err
				}
			}
			v.Field(i).Set(reflect.ValueOf(d))
			continue
		}

		parseFn, ok := baseTypeParsers[f.Type.Kind()]
		if !ok {
			continue
		}
		v1, err := parseFn(v0, format)
		if err != nil {
			return err
		}
		v.Field(i).Set(reflect.ValueOf(v1))
	}
	return nil
}

func parseInt(s, format string) (any, error) {
	return strconv.Atoi(s)
}

func parseInt8(s, format string) (any, error) {
	i, err := strconv.ParseInt(s, 10, 8)
	return int8(i), err
}

func parseInt16(s, format string) (any, error) {
	i, err := strconv.ParseInt(s, 10, 16)
	return int16(i), err
}

func parseInt32(s, format string) (any, error) {
	i, err := strconv.ParseInt(s, 10, 32)
	return int32(i), err
}

func parseInt64(s, format string) (any, error) {
	return strconv.ParseInt(s, 10, 64)
}

func parseUint(s, format string) (any, error) {
	u, err := strconv.ParseUint(s, 10, 64)
	return uint(u), err
}

func parseUint8(s, format string) (any, error) {
	u, err := strconv.ParseUint(s, 10, 8)
	return uint8(u), err
}

func parseUint16(s, format string) (any, error) {
	u, err := strconv.ParseUint(s, 10, 16)
	return uint16(u), err
}

func parseUint32(s, format string) (any, error) {
	u, err := strconv.ParseUint(s, 10, 32)
	return uint32(u), err
}

func parseUint64(s, format string) (any, error) {
	return strconv.ParseUint(s, 10, 64)
}

func parseBool(s, format string) (any, error) {
	if len(format) > 0 {
		return s == format, nil
	}
	return strconv.ParseBool(s)
}

func parseString(s, format string) (any, error) {
	return s, nil
}

// func parseTime(s, format string) (any, error) {
// 	return time.ParseInLocation(format, s, time.Local)
// }

// func parseDuration(s, format string) (any, error) {
// 	return time.ParseDuration(s)
// }
