package corehttp

import (
	"fmt"
	"github.com/gin-contrib/static"
	"github.com/hamster-shared/hamster-gateway/core/context"
)

func StartApi(ctx *context.CoreContext) error {
	r := NewMyServer(ctx)
	// router
	v1 := r.Group("/api/v1")
	{

		// basic configuration
		config := v1.Group("/config")
		{
			config.GET("/settting", getConfig)
			config.POST("/settting", setConfig)
			config.POST("/boot", setBootState)
			config.GET("/boot", getBootState)
		}
		p2p := v1.Group("/p2p")
		{
			p2p.GET("/bw", getP2pBW)
		}
	}
	//r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Use(static.Serve("/", static.LocalFile("./frontend/dist", false)))
	// listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	port := ctx.GetConfig().ApiPort
	return r.Run(fmt.Sprintf("0.0.0.0:%d", port))
}
