package service

type RemoteServiceClient interface {
	Invoke(method string, data interface{}, md map[string]string) (interface{}, error)
}
