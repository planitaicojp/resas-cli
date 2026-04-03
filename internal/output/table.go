package output

import (
	"fmt"
	"io"
	"reflect"
	"strings"
	"text/tabwriter"
)

type TableFormatter struct{}

func (f *TableFormatter) Format(w io.Writer, data any) error {
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Slice {
		return fmt.Errorf("tableフォーマットはスライス型のみ対応")
	}
	if v.Len() == 0 {
		return nil
	}

	tw := tabwriter.NewWriter(w, 0, 4, 2, ' ', 0)

	elemType := v.Type().Elem()
	headers := make([]string, 0, elemType.NumField())
	for i := 0; i < elemType.NumField(); i++ {
		tag := elemType.Field(i).Tag.Get("json")
		if tag == "" || tag == "-" {
			continue
		}
		name := strings.Split(tag, ",")[0]
		headers = append(headers, strings.ToUpper(name))
	}
	if _, err := fmt.Fprintln(tw, strings.Join(headers, "\t")); err != nil {
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
		if _, err := fmt.Fprintln(tw, strings.Join(vals, "\t")); err != nil {
			return err
		}
	}

	return tw.Flush()
}
