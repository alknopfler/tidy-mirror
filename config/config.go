package config

import (
	"fmt"
	"github.com/TwiN/go-color"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

//ZTPConfig is the global configuration data model
type MirrorConfig struct {
	PullSecretTempFile     string
	ConfigFile             string
	Kubeconfig             string
	RegistryCertPath       string
	PullSecretNS           string
	PullSecretName         string
	RegistryOCPDestIndexNS string
	RegistryOLMSourceIndex string
	RegistryOLMDestIndexNS string
	MarketplaceNS          string
	OwnCatalogName         string
	RegistryOCPRelease     string   `yaml:"ocp_release_version"`
	RegistryURL            string   `yaml:"registry_url"`
	RegistryUser           string   `yaml:"registry_username"`
	RegistryPass           string   `yaml:"registry_password"`
	ListPackages           []string `yaml:"list_packages"`
	ExtraImagesToMirror    []string `yaml:"extra_images_to_mirror"`
}

//fmt.Println(e.Spokes[0].Name, e.Spokes[0].Master0.NicExtDhcp)

//Constructor new config file from file
func NewConfig(configPath string, kubeconfig string) (MirrorConfig, error) {
	//Read main config from the config file
	var conf = MirrorConfig{
		Kubeconfig:             kubeconfig,
		PullSecretTempFile:     "/tmp/pull-secret-temp.json",
		RegistryCertPath:       "/etc/pki/ca-trust/source/anchors",
		PullSecretNS:           "openshift-config",
		PullSecretName:         "pull-secret",
		RegistryOCPDestIndexNS: "ocp4/openshift4",
		RegistryOLMSourceIndex: "registry.redhat.io/redhat/redhat-operator-index:v",
		RegistryOLMDestIndexNS: "olm/redhat-operator-index",
		MarketplaceNS:          "openshift-marketplace",
		OwnCatalogName:         "Tmirror Catalog",
	}

	if configPath == "" {
		return conf, fmt.Errorf(color.InRed("configFile param is empty"), "")
	}
	//Read config from file
	f, err := os.Open(configPath)
	if err != nil {
		return conf, fmt.Errorf(color.InRed("opening config file %s: %v"), conf.ConfigFile, err)
	}

	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(conf); err != nil {
		return conf, fmt.Errorf(color.InRed("decoding config file %s: %v"), conf.ConfigFile, err)
	}
	if err != nil {
		fmt.Println(color.InRed(">>>> [ERROR] Error reading config file: " + err.Error()))
		return conf, err
	}
	fmt.Println(color.InYellow(">>>> [INFO] ConfigFile param is not empty. Reading file from this path: " + configPath))

	// Set the rest of config from env
	if kubeconfig == "" {
		return conf, fmt.Errorf(color.InRed("Kubeconfig param is empty"), "")
	}
	fmt.Println(color.InYellow(">>>> [INFO] KUBECONFIG env is not empty. Reading file from this path: " + kubeconfig))
	conf.Kubeconfig = kubeconfig

	//modify config for source index depending on the config read from file
	conf.RegistryOLMSourceIndex += strings.Join(strings.Split(conf.RegistryOCPRelease, ".")[:2], ".")

	return conf, nil
}

//ReadFromConfigFile reads the config file
func (c *MirrorConfig) readFromConfigFile() error {

	f, err := os.Open(c.ConfigFile)
	if err != nil {
		return fmt.Errorf(color.InRed("opening config file %s: %v"), c.ConfigFile, err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(c); err != nil {
		return fmt.Errorf(color.InRed("decoding config file %s: %v"), c.ConfigFile, err)
	}
	return nil
}
