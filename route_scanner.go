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
		log.Printf("  scanRoutes %s: %+v", m.Name, methodType)

		if methodType.Kind() == reflect.Func && methodType.NumIn() == 1 && methodType.In(0) == ptrType &&
			methodType.NumOut() == 1 && methodType.Out(0) == reflect.TypeFor[Route]() {
			log.Println("We Ha!")
		}
	}

	return routes
}
