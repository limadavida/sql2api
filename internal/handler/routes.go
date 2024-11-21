package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/limadavida/sql2api/internal/config"
)

func Router() {

	handler := NewHandler(*config.ConfigData)
	r := gin.Default()

	routeMap := map[string]func(*gin.Context){
		"POST":   handler.Post(),
		"GET":    handler.Get(),
		"PUT":    handler.Put(),
		"DELETE": handler.Del(),
	}

	for method, handler := range routeMap {
		basicRoute := "/" + config.ConfigData.Project
		switch method {
		case "POST":
			r.POST(basicRoute, handler)
		case "GET":
			r.GET(basicRoute, handler)
		case "PUT":
			r.PUT(basicRoute, handler)
		case "DELETE":
			r.DELETE(basicRoute, handler)
		}
	}

	port := config.ConfigData.Servers[0]
	r.Run(fmt.Sprintf(":%d", port))
}
