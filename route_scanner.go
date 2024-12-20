package mbr

import (
	"log"
	"reflect"
)

func scanRoutes(ctrl Controller) (routes []Route) {
	ptrType := reflect.TypeOf(ctrl)
	elementType := ptrType.Elem()

	log.Println("scanRoutes: " + elementType.String())

	for i := 0; i < ptrType.NumMethod(); i++ {
		m := ptrType.Method(i)
		methodType := m.Type
		log.Printf("  method %s: %+v", m.Name, methodType)

		if methodType.Kind() == reflect.Func && //it is a function
			methodType.NumIn() == 1 && methodType.In(0) == ptrType && // with one arg which is pointer receiver to struct
			methodType.NumOut() == 1 && methodType.Out(0) == reflect.TypeFor[Route]() { // returning one value and this value is Route
			//} COMMENT TO MARK if conditions end [crazy go formatting. easier to accept rather then fight]
			route := m.Func.Call([]reflect.Value{reflect.ValueOf(ctrl)})[0].Interface().(Route)

			route.name = elementType.String() + "." + m.Name

			routes = append(routes, route)
		}
	}

	return routes
}
