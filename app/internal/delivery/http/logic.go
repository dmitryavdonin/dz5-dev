package delivery

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type AppVersion struct {
	Version string
}

func (d *Delivery) GetVersion(c *gin.Context) {

	//v := `{"version":"` + d.version + `"}`

	c.JSON(http.StatusOK, d.version)
}
