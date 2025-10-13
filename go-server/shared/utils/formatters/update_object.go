package formatters

import "reflect"

func UpdateObject(original, updates interface{}) {
	originalValue := reflect.ValueOf(original).Elem()
	updatesValue := reflect.ValueOf(updates).Elem()

	for i := 0; i < updatesValue.NumField(); i++ {
		field := updatesValue.Field(i)
		if field.IsValid() && field.Interface() != reflect.Zero(field.Type()).Interface() {
			originalValue.FieldByName(updatesValue.Type().Field(i).Name).Set(field)
		}
	}
}
