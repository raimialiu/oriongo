package constants

import "strings"

type Environment string

const (
	Production  Environment = "production"
	Development Environment = "development"
	STAGING     Environment = "staging"
	SANDBOXENV  Environment = "sandbox"
)

func (e Environment) String() string {
	return string(e)
}

func (e Environment) Name() string {
	return string(e)
}

func GetEnvName(env Environment) string {
	return strings.Join([]string{ENV_PREFIX, strings.ToUpper(env.String())}, "")
}
