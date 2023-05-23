package main

import (
	"context"

	"terraform-provider-excel/excel"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {
	providerserver.Serve(context.Background(), excel.New, providerserver.ServeOpts{
		Address: "deathfireofdoom.com/edu/excel",
	})
}
