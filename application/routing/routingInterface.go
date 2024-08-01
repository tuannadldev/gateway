package routing

import (
	"gateway/application/model"
)

type RoutingUseCase interface {
	Forward(routingData *model.RoutingData) (interface{}, error)
}
