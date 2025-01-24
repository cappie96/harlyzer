package harlyzer

import (
	"fmt"
	"reflect"
	"strings"
)

func formatUrl(url string) string {
	return strings.Split(url, "/")[2]
}

func formatTimings(s interface{}) string {
	v := reflect.ValueOf(s)
	t := reflect.TypeOf(s)
	var formattedTimings string
	formattedTimings += fmt.Sprintf("%-10s | %s\n", "Phase", "Time (ms)")
	formattedTimings += "------------------------\n"
	for i := 0; i < v.NumField(); i++ {
		fieldName := t.Field(i).Name
		fieldValue := v.Field(i).Float()
		formattedTimings += fmt.Sprintf("%-10s | %.3f"+
			"ms\n", fieldName, fieldValue)
	}
	return formattedTimings
}

func formatHeaders(headers []Header) string {
	var formattedHeaders string
	for _, header := range headers {
		formattedHeaders += fmt.Sprintf("%s: %s\n", header.Name, header.Value)
	}
	return formattedHeaders
}
