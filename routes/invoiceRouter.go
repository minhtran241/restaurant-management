package routes

import (
	"github.com/gin-gonic/gin"

	controller "github.com/minhtran241/restaurant-management/controllers"
)

func InvoiceRoutes(in *gin.Engine) {
	in.GET("/invoices", controller.GetInvoices())
	in.GET("/invoices/:invoice_id", controller.GetInvoice())
	in.POST("/invoices", controller.CreateInvoice())
	in.PATCH("/invoices/:invoice_id", controller.UpdateInvoice())
}