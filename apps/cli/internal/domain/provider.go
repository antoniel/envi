package domain

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
