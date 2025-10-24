package types

import "fmt"

func ToStr(config Configuration) string {
	return fmt.Sprintf("Config is %#v", config)
}

type Configuration struct {
	Environment environment `mapstructure:"ENVIRONMENT"`
	Identity    identity    `mapstructure:"IDENTITY"`
	Hosts       hosts       `mapstructure:"HOSTS"`
	Users       users       `mapstructure:"USERS"`
	Namespace   string      `mapstructure:"NAMESPACE"`
}

type hosts struct {
	RsBaseUri            string `mapstructure:"RS_BASE_URI"`
	IdentityPlatformFQDN string `mapstructure:"IDENTITY_PLATFORM_FQDN"`
	Scheme               string `mapstructure:"SCHEME"`
}

type identity struct {
	AmRealm string `mapstructure:"AM_REALM"`
}

type environment struct {
	Verbose   bool   `mapstructure:"VERBOSE"`
	Strict    bool   `mapstructure:"STRICT"`
	CloudType string `mapstructure:"CLOUDTYPE"`
	Paths     paths  `mapstructure:"PATHS"`
	SapigType string `mapstructure:"SAPIGTYPE"`
}

type users struct {
	CDKPlatformAdminUsername      string `mapstructure:"CDK_PLATFORM_ADMIN_USERNAME"`
	CDKPlatformAdminPassword      string `mapstructure:"CDK_PLATFORM_ADMIN_PASSWORD"`
	FIDCPlatformServiceAccountId  string `mapstructure:"FIDC_PLATFORM_SERVICE_ACCOUNT_ID"`
	FIDCPlatformServiceAccountKey string `mapstructure:"FIDC_PLATFORM_SERVICE_ACCOUNT_KEY"`
	PsuUserId                     string `mapstructure:"PSU_USER_ID"`
	PsuUsername                   string `mapstructure:"PSU_USERNAME"`
	PsuPassword                   string `mapstructure:"PSU_PASSWORD"`
	StudioUserId                  string `mapstructure:"STUDIO_USER_ID"`
	StudioUsername                string `mapstructure:"STUDIO_USERNAME"`
	StudioPassword                string `mapstructure:"STUDIO_PASSWORD"`
}

type paths struct {
	ConfigAuthHelper string `mapstructure:"CONFIG_AUTH_HELPER"`
}
