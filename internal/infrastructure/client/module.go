package client

import (
	"go.uber.org/fx"

	catalogapi "github.com/Sokol111/ecommerce-catalog-service-api/gen/httpapi"
	httpclient "github.com/Sokol111/ecommerce-commons/pkg/http/client"
)

// AttributeClientModule provides AttributeClient with its dependencies
func AttributeClientModule() fx.Option {
	return fx.Module("attribute-client",
		fx.Provide(
			fx.Private,
			httpclient.ProvideHTTPClient("catalog-service"),
		),
		fx.Provide(catalogapi.ProvideClient),
		fx.Provide(newAttributeClient),
	)
}
