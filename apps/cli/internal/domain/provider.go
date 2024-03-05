package domain

type ProviderTemplate struct {
	Name ProviderName
	Pull PullProvider
	Push PushProvider
}

type ProviderName string
type Provider interface {
	GetName() ProviderName
}
type PullProvider interface {
	GetName() ProviderName
	PullRemoteEnvValues() (EnvString, error)
}
type PushProvider interface {
	GetName() ProviderName
	PushLocalEnvValues(EnvString) error
}
type PushPullProvider interface {
	Provider
	PullProvider
	PushProvider
}
