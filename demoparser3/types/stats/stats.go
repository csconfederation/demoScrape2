package stats

import (
	"errors"
	"fmt"
	"reflect"
)

type Stats interface {
	Aggregate(newStats Stats) error
}

func Sum(current interface{}, new interface{}) error {
	statsVal := reflect.ValueOf(current)
	otherVal := reflect.ValueOf(new)

	// Must be pointers
	if statsVal.Kind() != reflect.Ptr || otherVal.Kind() != reflect.Ptr {
		return errors.New("aggregate called on non-pointer value")
	}

	statsElem := statsVal.Elem()
	otherElem := otherVal.Elem()

	// Must be same type
	if statsElem.Type() != otherElem.Type() {
		return fmt.Errorf("types must match: %v != %v", statsElem.Type(), otherElem.Type())
	}

	statsType := statsElem.Type()

	// Iterate through all fields
	for i := 0; i < statsElem.NumField(); i++ {
		field := statsType.Field(i)
		statsField := statsElem.Field(i)
		otherField := otherElem.Field(i)

		// Check if field has the aggregation tag
		tag := field.Tag.Get("end_of_match_sum")
		if tag != "true" {
			continue // Skip fields without the tag
		}

		// Skip unexported fields
		if !statsField.CanSet() {
			continue
		}

		// In your Sum() function, add this before the switch:
		if statsField.Kind() == reflect.Struct {
			// Recursively sum nested struct fields
			for j := 0; j < statsField.NumField(); j++ {
				nestedField := statsField.Type().Field(j)
				nestedTag := nestedField.Tag.Get("end_of_match_sum")

				if nestedTag == "true" {
					nestedStatsField := statsField.Field(j)
					nestedOtherField := otherField.Field(j)

					if nestedStatsField.CanSet() {
						switch nestedStatsField.Kind() {
						case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
							nestedStatsField.SetInt(nestedStatsField.Int() + nestedOtherField.Int())
						case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
							nestedStatsField.SetUint(nestedStatsField.Uint() + nestedOtherField.Uint())
						case reflect.Float32, reflect.Float64:
							nestedStatsField.SetFloat(nestedStatsField.Float() + nestedOtherField.Float())
						}
					}
				}
			}
			continue
		}

		// Sum based on field type
		switch statsField.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			statsField.SetInt(statsField.Int() + otherField.Int())

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			statsField.SetUint(statsField.Uint() + otherField.Uint())

		case reflect.Float32, reflect.Float64:
			statsField.SetFloat(statsField.Float() + otherField.Float())

		default:
			// Skip unsupported types (strings, maps, slices, etc.)
			continue
		}
	}

	return nil
}
