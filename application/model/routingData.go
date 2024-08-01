package model

type RoutingData struct {
	ServiceName   string
	ServiceMethod string
	Payload       interface{}
	Metadata      map[string]string
}
