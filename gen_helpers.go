// DO NOT MODIFY! AUTO GENERATED BY igo v1.0.2 (https://github.com/rocketlaunchr/igo)

// Copyright 2019-20 PJ Engineering and Business Solutions Pty. Ltd. All rights reserved.

package dbq

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/civil"
)

// Database is used to set the Database.
// Different databases have different syntax for placeholders etc.
type Database int

const (
	// MySQL database
	MySQL Database = 0
	// PostgreSQL database
	PostgreSQL Database = 1
)

// INSERT will generate an INSERT statement.
func INSERT(tableName string, columns []string, rows int, dbtype ...Database) string {
	return fmt.Sprintf("INSERT INTO %s ( %s ) VALUES %s", tableName, strings.Join(columns, ","), Ph(len(columns), rows, 0, dbtype...))
}

// Ph generates the placeholders for SQL queries.
// For a bulk insert operation, rows is the number of rows you intend
// to insert, and columnsN is the number of fields per row.
// For the IN function, set rows to 1.
// For PostgreSQL, you can use incr to increment the placeholder starting count.
//
// NOTE: The function panics if either columnsN or rows is 0.
func Ph(columnsN, rows int, incr int, dbtype ...Database) string {

	var typ Database
	if len(dbtype) > 0 {
		typ = dbtype[0]
	}

	if columnsN == 0 {
		panic(errors.New("columnsN must not be 0"))
	}

	if rows == 0 {
		panic(errors.New("rows must not be 0"))
	}

	if typ == MySQL {
		inner := "( " + strings.TrimSuffix(strings.Repeat("?,", columnsN), ",") + " ),"
		return strings.TrimSuffix(strings.Repeat(inner, rows), ",")
	}

	var singleValuesStr string

	varCount := 1 + incr
	for i := 1; i <= rows; i++ {
		singleValuesStr = singleValuesStr + "("
		for j := 1; j <= columnsN; j++ {
			singleValuesStr = singleValuesStr + fmt.Sprintf("$%d,", varCount)
			varCount++
		}
		singleValuesStr = strings.TrimSuffix(singleValuesStr, ",") + "),"
	}

	return strings.TrimSuffix(singleValuesStr, ",")
}

func sliceConv(arg reflect.Value) []interface{} {
	out := []interface{}{}

	if arg.Kind() == reflect.Slice {
		for i := 0; i < arg.Len(); i++ {
			out = append(out, sliceConv(reflect.ValueOf(arg.Index(i).Interface()))...)
		}
	} else {
		out = append(out, arg.Interface())
	}

	return out
}

// StdTimeConversionConfig provides a standard configuration for unmarshaling to
// time-related fields in a struct. It properly converts timestamps and datetime columns into
// time.Time objects. It assumes a MySQL database as default.
func StdTimeConversionConfig(dbtype ...Database) *StructorConfig {

	layouts := []string{
		"2006-01-02 15:04:05",
		time.RFC3339,
	}

	if len(dbtype) > 0 && dbtype[0] == PostgreSQL {

		layouts[0], layouts[1] = layouts[1], layouts[0]
	}

	return &StructorConfig{
		WeaklyTypedInput: true,
		DecodeHook: func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
			if f.Kind() != reflect.String {
				return data, nil
			}

			switch t {
			case reflect.TypeOf(civil.Date{}):
				return civil.ParseDate(data.(string))
			case reflect.TypeOf(civil.DateTime{}):
				t, err := time.Parse(layouts[0], data.(string))
				if err != nil {
					t, err = time.Parse(layouts[1], data.(string))
					if err != nil {
						return nil, err
					}
				}
				return civil.DateTime{
					Date: civil.DateOf(t),
					Time: civil.TimeOf(t),
				}, nil
			case reflect.TypeOf(civil.Time{}):
				return civil.ParseTime(data.(string))
			case reflect.TypeOf(time.Time{}):
				t, err := time.Parse(layouts[0], data.(string))
				if err != nil {
					t, err := time.Parse(layouts[1], data.(string))
					if err != nil {
						return nil, err
					}
					return t, nil
				}
				return t, nil
			default:
				return data, nil
			}

			return data, nil
		},
	}
}

// Struct converts the fields of the struct into a slice of values.
// You can use it to convert a struct into the placeholder arguments required by
// the Q and E function. tagName is used to indicate the struct tag (default is "dbq").
// The function panics if strct is not an actual struct.
func Struct(strct interface{}, tagName ...string) []interface{} {

	tg := "dbq"

	if len(tagName) > 0 {
		tg = tagName[0]
	}

	out := []interface{}{}

	if strct == nil {
		panic(errors.New("strct must be a struct"))
	}

	s := reflect.ValueOf(strct)

	if s.Kind() == reflect.Ptr {
		s = reflect.Indirect(s)
	}
	typeOfT := s.Type()

	for i := 0; i < s.NumField(); i++ {
		f := typeOfT.Field(i)

		if f.PkgPath != "" {

			continue
		}

		fieldTag := f.Tag.Get(tg)
		fieldValRaw := s.Field(i)
		fieldVal := fieldValRaw.Interface()

		if fieldValRaw.Kind() == reflect.Map {
			continue
		}

		if fieldTag == "-" || (strings.HasSuffix(fieldTag, ",omitempty") && reflect.DeepEqual(fieldVal, reflect.Zero(reflect.TypeOf(fieldVal)).Interface())) {
			continue
		}

		if fieldValRaw.Kind() == reflect.Slice {
			out = append(out, sliceConv(fieldValRaw)...)
			continue
		}

		out = append(out, fieldVal)
	}

	return out
}

func parseUintP(s string) *uint {
	n, _ := strconv.ParseUint(s, 10, 0)
	return &[]uint{uint(n)}[0]
}

func parseUint8P(s string) *uint8 {
	n, _ := strconv.ParseUint(s, 10, 8)
	return &[]uint8{uint8(n)}[0]
}

func parseUint16P(s string) *uint16 {
	n, _ := strconv.ParseUint(s, 10, 16)
	return &[]uint16{uint16(n)}[0]
}

func parseUint32P(s string) *uint32 {
	n, _ := strconv.ParseUint(s, 10, 32)
	return &[]uint32{uint32(n)}[0]
}

func parseUint64P(s string) *uint64 {
	n, _ := strconv.ParseUint(s, 10, 64)
	return &[]uint64{uint64(n)}[0]
}

func parseIntP(s string) *int {
	n, _ := strconv.ParseInt(s, 10, 0)
	return &[]int{int(n)}[0]
}

func parseInt8P(s string) *int8 {
	n, _ := strconv.ParseInt(s, 10, 8)
	return &[]int8{int8(n)}[0]
}

func parseInt16P(s string) *int16 {
	n, _ := strconv.ParseInt(s, 10, 16)
	return &[]int16{int16(n)}[0]
}

func parseInt32P(s string) *int32 {
	n, _ := strconv.ParseInt(s, 10, 32)
	return &[]int32{int32(n)}[0]
}

func parseInt64P(s string) *int64 {
	n, _ := strconv.ParseInt(s, 10, 64)
	return &[]int64{int64(n)}[0]
}

func parseUint(s string) uint {
	n, _ := strconv.ParseUint(s, 10, 0)
	return uint(n)
}

func parseUint8(s string) uint8 {
	n, _ := strconv.ParseUint(s, 10, 8)
	return uint8(n)
}

func parseUint16(s string) uint16 {
	n, _ := strconv.ParseUint(s, 10, 16)
	return uint16(n)
}

func parseUint32(s string) uint32 {
	n, _ := strconv.ParseUint(s, 10, 32)
	return uint32(n)
}

func parseUint64(s string) uint64 {
	n, _ := strconv.ParseUint(s, 10, 64)
	return n
}

func parseInt(s string) int {
	n, _ := strconv.ParseInt(s, 10, 0)
	return int(n)
}

func parseInt8(s string) int8 {
	n, _ := strconv.ParseInt(s, 10, 8)
	return int8(n)
}

func parseInt16(s string) int16 {
	n, _ := strconv.ParseInt(s, 10, 16)
	return int16(n)
}

func parseInt32(s string) int32 {
	n, _ := strconv.ParseInt(s, 10, 32)
	return int32(n)
}

func parseInt64(s string) int64 {
	n, _ := strconv.ParseInt(s, 10, 64)
	return n
}
