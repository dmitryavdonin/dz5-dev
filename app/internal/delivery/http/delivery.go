package delivery

import (
	"fmt"

	"github.com/dmitryavdonin/gtools/logger"
	"github.com/gin-gonic/gin"
)

type Delivery struct {
	router  *gin.Engine
	logger  logger.Interface
	port    int
	version string

	options Options
}

type Options struct{}

func New(verison string, port int, logger logger.Interface, options Options) (*Delivery, error) {

	var d = &Delivery{
		logger:  logger,
		port:    port,
		version: verison,
	}

	d.SetOptions(options)

	d.router = d.initRouter()
	return d, nil
}

func (d *Delivery) SetOptions(options Options) {
	if d.options != options {
		d.options = options
	}
}

func (d *Delivery) Run() error {
	return d.router.Run(fmt.Sprintf(":%d", d.port))
}
