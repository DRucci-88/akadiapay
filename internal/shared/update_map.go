package shared

import "reflect"

type UpdateMap map[string]any

func (u UpdateMap) SetIfNotNil(key string, value any) {
	if value == nil {
		return
	}

	rv := reflect.ValueOf(value)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return
	}

	u[key] = rv.Elem().Interface()
}

// func (u UpdateMap) SetIfNotNil(key string, value any) {
// 	if value == nil {
// 		return
// 	}

// 	switch v := value.(type) {
// 	case *string:
// 		u[key] = *v
// 	case *bool:
// 		u[key] = *v
// 	case *float64:
// 		u[key] = *v
// 	case *int:
// 		u[key] = *v
// 	}
// }
