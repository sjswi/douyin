package config

type RPCSrvConfig struct {
	Name      string `mapstructure:"name" json:"name"`
	Namespace string `mapstructure:"namespace" json:"namespace"`
}
