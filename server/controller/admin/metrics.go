package admin

import (
	"net/http"

	"server/common/metrics"
	"server/controller"

	"github.com/gin-gonic/gin"
)

type MetricsOverviewResponse struct {
	controller.Response
	Overview metrics.Overview `json:"overview"`
}

type RouteMetricsResponse struct {
	controller.Response
	Routes []metrics.RouteSnapshot `json:"routes"`
}

type UserMetricsResponse struct {
	controller.Response
	Users []metrics.UserSnapshot `json:"users"`
}

type BusinessCodeMetricsResponse struct {
	controller.Response
	BusinessCodes []metrics.BusinessCodeSnapshot `json:"businessCodes"`
}

type ModelMetricsResponse struct {
	controller.Response
	Models []metrics.ModelSnapshot `json:"models"`
}

type ModelFailureMetricsResponse struct {
	controller.Response
	Failures []metrics.ModelFailureSnapshot `json:"failures"`
}

func MetricsOverview(c *gin.Context) {
	res := new(MetricsOverviewResponse)
	res.Success()
	res.Overview = metrics.GetCollector().Overview()
	c.JSON(http.StatusOK, res)
}

func RouteMetrics(c *gin.Context) {
	res := new(RouteMetricsResponse)
	res.Success()
	res.Routes = metrics.GetCollector().RouteSnapshots()
	c.JSON(http.StatusOK, res)
}

func UserMetrics(c *gin.Context) {
	res := new(UserMetricsResponse)
	res.Success()
	res.Users = metrics.GetCollector().UserSnapshots()
	c.JSON(http.StatusOK, res)
}

func BusinessCodeMetrics(c *gin.Context) {
	res := new(BusinessCodeMetricsResponse)
	res.Success()
	res.BusinessCodes = metrics.GetCollector().BusinessCodeSnapshots()
	c.JSON(http.StatusOK, res)
}

func ModelMetrics(c *gin.Context) {
	res := new(ModelMetricsResponse)
	res.Success()
	res.Models = metrics.GetCollector().ModelSnapshots()
	c.JSON(http.StatusOK, res)
}

func ModelFailureMetrics(c *gin.Context) {
	res := new(ModelFailureMetricsResponse)
	res.Success()
	res.Failures = metrics.GetCollector().ModelFailureSnapshots()
	c.JSON(http.StatusOK, res)
}
