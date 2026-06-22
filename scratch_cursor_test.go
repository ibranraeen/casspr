package main

import (
	"fmt"
	"reflect"

	tea "charm.land/bubbletea/v2"
)

func main() {
	var cur tea.Cursor
	t := reflect.TypeOf(cur)
	fmt.Printf("Type: %s\n", t.String())
	fmt.Printf("Kind: %s\n", t.Kind().String())
	if t.Kind() == reflect.Struct {
		fmt.Printf("NumField: %d\n", t.NumField())
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			fmt.Printf("  Field %d: Name=%s, Type=%s, PkgPath=%s\n", i, field.Name, field.Type, field.PkgPath)
		}
	}
}
