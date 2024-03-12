package helpers

import (
	"fmt"
	"reflect"
)

func PrettyPrintStruct(v interface{}) {
	prettyPrintStructRecursive(v, "")
}

func prettyPrintStructRecursive(v interface{}, indent string) {
	indentIncrementation := "    "

	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// Check if the value is a struct
	if val.Kind() != reflect.Struct {
		fmt.Printf("%s%v\n", indent, val.Interface())
		return
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		typeField := val.Type().Field(i)
		fieldName := typeField.Name

		// Skip empty strings
		if field.Kind() == reflect.String && field.String() == "" {
			continue
		}

		// Add conditions to format specific types if needed, like *big.Int or common.Address
		switch field.Kind() {
		case reflect.Struct:
			fmt.Println(indent + fieldName + ":")
			prettyPrintStructRecursive(field.Interface(), indent+indentIncrementation)
		case reflect.Ptr:
			if !field.IsNil() {
				fmt.Println(indent + fieldName + ":")
				prettyPrintStructRecursive(field.Elem().Interface(), indent+indentIncrementation)
			}
		default:
			fmt.Printf("%s%s: %v\n", indent, fieldName, field.Interface())
		}
	}
}
