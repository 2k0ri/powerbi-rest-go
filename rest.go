package powerbi

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"time"
)

type DataSet struct {
	Name   string  `json:"name"`
	Tables []Table `json:"tables"`
}

type Table struct {
	Name    string   `json:"name"`
	Columns []Column `json:"columns"`
}

type Column struct {
	Name     string `json:"name"`
	DataType string `json:"dataType"`
}

type Rows struct {
	Rows []interface{} `json:"rows"`
}

func StructToColumns(o interface{}) []Column {
	v := reflect.ValueOf(o)
	vt := v.Type()
	c := make([]Column, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		t := v.Field(i).Type()
		r, err := oDataReflect(t)
		if err != nil {
			log.Fatal(err)
		}
		c[i] = Column{Name: vt.Field(i).Name, DataType: r}
	}
	return c
}

// oDataReflect returns PowerBI supported OData type from reflect.Type
// @see https://msdn.microsoft.com/en-us/library/mt203569.aspx
func oDataReflect(t reflect.Type) (string, error) {
	switch t.Kind() {
	case reflect.Bool:
		return "Boolean", nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "Int64", nil
	case reflect.Float32, reflect.Float64:
		return "Double", nil
	case reflect.String:
		return "String", nil
	case reflect.Struct:
		// only passed when it is time.Time
		if t == reflect.TypeOf(time.Time{}) {
			return "Datetime", nil
		}
	}
	return "", errors.New(fmt.Sprintf("%v is not matched in Int64, Double, bool, string and DateTime", t.String()))
}
