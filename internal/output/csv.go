package output

import (
	"encoding/csv"
	"fmt"
	"io"
	"reflect"
	"strings"
)

type CSVFormatter struct{}

func (f *CSVFormatter) Format(w io.Writer, data any) error {
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Slice {
		return fmt.Errorf("CSVフォーマットはスライス型のみ対応")
	}
	if v.Len() == 0 {
		return nil
	}

	cw := csv.NewWriter(w)

	elemType := v.Type().Elem()
	headers := make([]string, 0, elemType.NumField())
	for i := 0; i < elemType.NumField(); i++ {
		tag := elemType.Field(i).Tag.Get("json")
		if tag == "" || tag == "-" {
			continue
		}
		headers = append(headers, strings.Split(tag, ",")[0])
	}
	if err := cw.Write(headers); err != nil {
		return err
	}

	for i := 0; i < v.Len(); i++ {
		row := v.Index(i)
		vals := make([]string, 0, len(headers))
		for j := 0; j < row.NumField(); j++ {
			tag := row.Type().Field(j).Tag.Get("json")
			if tag == "" || tag == "-" {
				continue
			}
			vals = append(vals, fmt.Sprintf("%v", row.Field(j).Interface()))
		}
		if err := cw.Write(vals); err != nil {
			return err
		}
	}

	cw.Flush()
	return cw.Error()
}
