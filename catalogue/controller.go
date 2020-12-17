package catalogue

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pilillo/mastro/abstract"
	"github.com/pilillo/mastro/utils/conf"
	"github.com/pilillo/mastro/utils/errors"
)

const (
	assetRestEndpoint string = "asset"
	assetIDParam      string = "asset_id"
)

// Ping ... replies to a ping message for healthcheck purposes
func Ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

// CreateAsset ... creates an asset description entry
func CreateAsset(c *gin.Context) {
	asset := abstract.Asset{}
	if err := c.ShouldBindJSON(&asset); err != nil {
		restErr := errors.GetBadRequestError("Invalid JSON Body")
		c.JSON(restErr.Status, restErr)
	} else {
		result, saveErr := assetService.CreateAsset(asset)
		if saveErr != nil {
			c.JSON(saveErr.Status, saveErr)
		} else {
			c.JSON(http.StatusCreated, result)
		}
	}

}

// GetAssetByID ... retrieves an asset description by its Unique Name ID
func GetAssetByID(c *gin.Context) {
	nameID := c.Param(assetIDParam)
	asset, getErr := assetService.GetAssetByID(nameID)
	if getErr != nil {
		c.JSON(getErr.Status, getErr)
	} else {
		c.JSON(http.StatusOK, asset)
	}
}

// ListAllAssets ... returns all assets
func ListAllAssets(c *gin.Context) {
	assets, err := assetService.ListAllAssets()
	if err != nil {
		c.JSON(err.Status, err)
	} else {
		c.JSON(http.StatusOK, assets)
	}
}

var router = gin.Default()

// StartEndpoint ... starts the service endpoint
func StartEndpoint(cfg *conf.Config) {
	// init service
	assetService.Init(cfg)

	// add an healthcheck for the endpoint
	router.GET(fmt.Sprintf("healthcheck/%s", assetRestEndpoint), Ping)

	// get asset set as asset/id
	router.GET(fmt.Sprintf("%s/:%s", assetRestEndpoint, assetIDParam), GetAssetByID)

	// put asset set as asset/
	router.PUT(fmt.Sprintf("%s/", assetRestEndpoint), CreateAsset)

	// list all assets
	router.GET(fmt.Sprintf("%s/", assetRestEndpoint), ListAllAssets)

	// run router as standalone service
	// todo: do we need to run multiple endpoints from the main?
	router.Run(fmt.Sprintf(":%s", cfg.Details.Values["port"]))
}
