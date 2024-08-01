package utils

// Get config path by app environment
func GetConfigPath(configPath string) string {
	if configPath == "qc" {
		return "./config/qc"
	}
	if configPath == "staging" {
		return "./config/staging"
	}
	if configPath == "prod" {
		return "./config/prod"
	}
	return "./config/local"
}
