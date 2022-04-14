package corehttp

import (
	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/hamster-shared/hamster-gateway/core/modules/config"
	"net/http"
)

func getConfig(gin *MyContext) {
	cfg := gin.CoreContext.GetConfig()
	gin.JSON(http.StatusOK, Success(cfg))
}

func setConfig(gin *MyContext) {
	cfg := gin.CoreContext.GetConfig()
	reqBody := config.Config{}
	if err := gin.BindJSON(&reqBody); err != nil {
		gin.JSON(http.StatusBadRequest, BadRequest())
		return
	}

	cfg.ChainApi = reqBody.ChainApi
	// 校验seed 是否合法
	_, err := signature.KeyringPairFromSecret(reqBody.SeedOrPhrase, 42)
	if err != nil {
		gin.JSON(http.StatusBadRequest, BadRequest("seed not invalid"))
		return
	}

	cfg.SeedOrPhrase = reqBody.SeedOrPhrase
	cfg.PublicIp = reqBody.PublicIp
	cfg.PublicPort = reqBody.PublicPort

	err = gin.CoreContext.Cm.Save(cfg)
	if err != nil {
		gin.JSON(http.StatusBadRequest, BadRequest("save config fail"))
		return
	}

	gin.JSON(http.StatusOK, Success(""))
}

func setBootState(gin *MyContext) {

	var op BootParam

	if err := gin.BindJSON(&op); err != nil {
		gin.JSON(http.StatusBadRequest, BadRequest())
		return
	}

	if op.Option {
		gin.CoreContext.StateService.Start()
	} else {
		gin.CoreContext.StateService.Stop()
	}

	gin.JSON(http.StatusOK, Success(""))
}

func getBootState(gin *MyContext) {
	gin.JSON(http.StatusOK, Success(gin.CoreContext.StateService.Running()))
}

func getP2pBW(gin *MyContext) {
	nd := gin.CoreContext.StateService.Node
	if !nd.IsOnline {
		gin.JSON(http.StatusBadRequest, BadRequest("Not Online"))
		return
	}

	stats := nd.Reporter.GetBandwidthTotals()

	gin.JSON(http.StatusOK, Success(stats))
}
