package config

import (
	"fmt"
	"github.com/TwiN/go-color"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v3"
)

//ZTPConfig is the global configuration data model
type MirrorConfig struct {
	Mirror struct {
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
	} `yaml:"mirror"`
}

//fmt.Println(e.Spokes[0].Name, e.Spokes[0].Master0.NicExtDhcp)

//Constructor new config file from file
func NewConfig(configPath string, kubeconfig string) (MirrorConfig, error) {

	//Read main config from the config file
	if configPath == "" {
		return MirrorConfig{}, fmt.Errorf(color.InRed("configFile param is empty"), "")
	}

	conf, err := readFromConfigFile(configPath)
	if err != nil {
		return MirrorConfig{}, err
	}
	fmt.Println("config---->", conf)
	conf.Mirror.Kubeconfig = kubeconfig
	conf.Mirror.PullSecretTempFile = "/tmp/pull-secret-temp.json"
	conf.Mirror.RegistryCertPath = "/etc/pki/ca-trust/source/anchors"
	conf.Mirror.PullSecretNS = "openshift-config"
	conf.Mirror.PullSecretName = "pull-secret"
	conf.Mirror.RegistryOCPDestIndexNS = "ocp4/openshift4"
	conf.Mirror.RegistryOLMSourceIndex = "registry.redhat.io/redhat/redhat-operator-index:v"
	conf.Mirror.RegistryOLMDestIndexNS = "olm/redhat-operator-index"
	conf.Mirror.MarketplaceNS = "openshift-marketplace"
	conf.Mirror.OwnCatalogName = "Tmirror Catalog"
	fmt.Println("config-post----->", conf)

	// Set the rest of config from param
	if kubeconfig == "" {
		return conf, fmt.Errorf(color.InRed("Kubeconfig param is empty"), "")
	}
	fmt.Println(color.InYellow(">>>> [INFO] KUBECONFIG env is not empty. Reading file from this path: " + kubeconfig))
	conf.Mirror.Kubeconfig = kubeconfig

	//modify config for source index depending on the config read from file
	conf.Mirror.RegistryOLMSourceIndex += strings.Join(strings.Split(conf.Mirror.RegistryOCPRelease, ".")[:2], ".")

	fmt.Println("final config---->", conf)
	return conf, nil
}

//ReadFromConfigFile reads the config file
func readFromConfigFile(configFile string) (MirrorConfig, error) {
	var conf MirrorConfig
	f, err := ioutil.ReadFile(configFile)
	if err != nil {
		return MirrorConfig{}, fmt.Errorf(color.InRed("opening config file %s: %v"), configFile, err)
	}

	err = yaml.Unmarshal(f, conf)
	if err != nil {
		return MirrorConfig{}, fmt.Errorf(color.InRed("decoding config file %s: %v"), configFile, err)
	}
	return conf, nil
}
