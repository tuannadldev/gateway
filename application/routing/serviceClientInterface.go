package routing

import (
	"gateway/application/model"
)

type ServiceClient interface {
	Invoke(routingData *model.RoutingData) (interface{}, error)
}
