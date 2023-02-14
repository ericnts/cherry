package mate

type ResourceType int

const (
	RTController ResourceType = iota
	RTApplication
	RTService
	RTFactory
	RTRepository
	RTInfrastructure
)


