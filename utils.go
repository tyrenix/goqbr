package qbr

import (
	"reflect"
	"strings"
	"time"

	"github.com/tyrenix/qbr/domain"
)

// isZero check on zero value
func isZero(value any) bool {
	// check is nil
	if value == nil {
		return true
	}

	// get value by reflect
	v := reflect.ValueOf(value)

	// for pointer types
	if v.Kind() == reflect.Ptr ||
		v.Kind() == reflect.Slice ||
		v.Kind() == reflect.Map ||
		v.Kind() == reflect.Chan ||
		v.Kind() == reflect.Interface {
		return v.IsNil()
	}

	// is string check on empty
	if v.Kind() == reflect.String {
		return v.String() == ""
	}

	// is number check on zero value
	if (v.Kind() >= reflect.Int && v.Kind() <= reflect.Int64) ||
		(v.Kind() >= reflect.Float32 && v.Kind() <= reflect.Float64) {
		return v.IsZero()
	}

	// for Time types
	if v.Kind() == reflect.Struct && v.Type() == reflect.TypeOf(time.Time{}) {
		return v.Interface().(time.Time).IsZero()
	}

	// for other types return false
	return false
}

// isFieldIgnored checks if a field is ignored for a given query type.
//
// The function checks if the query type is in the field's list of ignored operations.
// If it is, the function returns true, indicating that the field is ignored. Otherwise,
// it returns false.
func isFieldIgnored(field *domain.Field, queryType domain.OperationType) bool {
	// check is ignored
	for _, ignoreOp := range field.IgnoreOn {
		if ignoreOp == queryType {
			return true
		}
	}

	// not ignored
	return false
}

// extractFieldFromStruct extracts a Field object from a given struct field.
//
// The function retrieves the "db" tag from the field annotation and uses it to
// initialize a Field object. If the "db" tag is empty, the function returns nil.
// Additionally, the function checks for a "qbr" tag and parses any annotations
// it contains. If the "qbr" tag includes an "ignore_on" annotation, the function
// extracts the ignored operations and adds them to the Field's IgnoredOperations
// slice.
//
// The resulting Field object is returned, representing a database field with
// optional ignored operations based on the struct field's annotations.
func extractFieldFromStruct(ft reflect.StructField) *domain.Field {
	// get tags from field annotation
	db := ft.Tag.Get(string(domain.QueryDB))

	// check is not empty
	if db == "" {
		return nil
	}

	// create field
	field := &domain.Field{
		DB: db,
	}

	// query builder tag
	qbr := ft.Tag.Get(string(domain.QueryQbr))

	// check is not empty
	if qbr == "" {
		return field
	}

	// get annotations from query builder annotation
	for _, block := range strings.Split(qbr, " ") {
		// check is not empty
		if block == "" {
			continue
		}

		// get annotation
		switch {
		case strings.HasPrefix(block, string(domain.QueryIgnoreOn)+"="):
			field.IgnoreOn = append(
				field.IgnoreOn,
				extractIgnoredOperationOnAnnotations(block)...,
			)
		default:
			continue
		}
	}

	// return fields
	return field
}

// extractIgnoredOperationOnAnnotations extracts the ignored operations from the given block string.
//
// The block string is expected to be in the format "ignore_on=<operation1>,<operation2>,...".
//
// The function splits the block by comma, trims the resulting strings, and adds them to a slice of
// ignored operations. The operation types are converted to lower case to ensure consistency.
//
// The function returns the slice of ignored operations.
func extractIgnoredOperationOnAnnotations(block string) []domain.OperationType {
	// delete from block annotation type
	block = strings.TrimPrefix(block, string(domain.QueryIgnoreOn)+"=")

	// split by comma
	ops := strings.Split(block, ",")

	// slice of ignored operations
	ignOps := make([]domain.OperationType, 0, len(ops))

	// add ignored operations
	for _, op := range ops {
		// check is not empty
		if op == "" {
			continue
		}

		// get operation type
		ignOps = append(ignOps, domain.OperationType(strings.ToLower(op)))
	}

	// return ignored operations
	return ignOps
}
