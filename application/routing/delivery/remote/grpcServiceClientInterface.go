package remote

type GrpcServiceClient interface {
	GetMethodRegistry() map[string]func(interface{}, map[string]string) (interface{}, error)
}
