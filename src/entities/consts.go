package entities

const (
	serviceModeEnv     = "SERVICE_MODE"
	serviceModeProd    = "PROD"
	serviceModeDevelop = "DEVELOP"
	serviceModeTest    = "TEST"
)

const (
	ServiceRequestTimeout = 3
)

// Viper config const's
const (
	ConfigPathWd         = "core.wd"
	ConfigServiceVersion = "ver"

	ServiceName     = "name"
	ServicesVersion = "services.version"
)

// RabbitMq
const (
	PhotosQueue = "photos"
)
