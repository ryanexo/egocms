package sqlmockutils

import (
    `database/sql/driver`
    `errors`
    `reflect`

    `github.com/DATA-DOG/go-sqlmock`
)

var ErrNonStruct = errors.New("the provided parameter is not a struct and is invalid")

type Rows struct {
    fields []reflect.StructField
    rows   *sqlmock.Rows
}

func (r Rows) Result() *sqlmock.Rows {
    return r.rows
}

func (r Rows) Add(data ...any) Rows {
    for _, s := range data {
        ref := reflect.ValueOf(s)
        if ref.Kind() == reflect.Ptr {
            ref = ref.Elem()
        }
        if ref.Kind() == reflect.Slice || ref.Kind() == reflect.Array {
            for i := 0; i < ref.Len(); i++ {
                r.Add(ref.Index(i).Interface())
            }
            continue
        }
        if ref.Kind() != reflect.Struct {
            continue
        }
        r.rows.AddRow(r.visibleValues(ref)...)
    }
    return r
}

func (r Rows) visibleValues(value reflect.Value) []driver.Value {
    var result []driver.Value
    for _, field := range r.fields {
        result = append(result, value.FieldByName(field.Name).Interface())
    }
    return result
}

func NewRows(definition any) Rows {
    ref := reflect.TypeOf(definition)
    if ref.Kind() == reflect.Ptr {
        ref = ref.Elem()
    }
    rawFields := reflect.VisibleFields(ref)
    var (
        actualFields []reflect.StructField
        stringFields []string
    )
    for _, field := range rawFields {
        switch field.Type.Kind() {
        case reflect.Ptr, reflect.Struct, reflect.Slice, reflect.Array, reflect.Chan, reflect.UnsafePointer:
            continue
        default:
            stringFields = append(stringFields, field.Name)
            actualFields = append(actualFields, field)
        }
    }

    return Rows{fields: actualFields, rows: sqlmock.NewRows(stringFields)}
}
