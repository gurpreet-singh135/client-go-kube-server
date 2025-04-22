
package main

import (
	"myapp/handlers"
	"myapp/util"
)

func main() {
	util.Initialize_client()
	handlers.Create_handlers()
}

