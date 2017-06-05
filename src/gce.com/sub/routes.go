package sub

import (
	"net/http"
)

//Route struct is a
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

//Routes array is a
type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
	Route{
		"GetLastRecord",
		"GET",
		"/GetLastRecord/{host}/{app}",
		GetLastRecord,
	},
	Route{
		"GetLastRecordSet",
		"POST",
		"/GetLastRecordSet",
		GetLastRecordSet,
	},
}

/*Route{
      "TodoIndex",
      "GET",
      "/todos",
      TodoIndex,
  },
  Route{
      "TodoShow",
      "GET",
      "/todos/{todoId}",
      TodoShow,
  },*/
