package registry

import (
	"context"
	"github.com/TwiN/go-color"
	"github.com/alknopfler/tidy-mirror/config"
	"github.com/alknopfler/tidy-mirror/pkg/auth"
	a "github.com/containers/common/pkg/auth"
	"github.com/containers/image/v5/types"
	"github.com/operator-framework/api/pkg/operators/v1alpha1"
	"io/ioutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	"os"
	"strings"
	"time"
)

//type FileServer
type Registry struct {
	config.MirrorConfig
}

func NewRegistry(conf config.MirrorConfig) *Registry {
	return &Registry{
		MirrorConfig: conf,
	}
}

/*Constructor NewFileServer
func NewRegistry(conf config.MirrorConfig) *Registry {
	return &Registry{

		PullSecretName:          "pull-secret",
		PullSecretNS:            "openshift-config",
		RegistrySrcPkg:          "kubernetes-nmstate-operator,metallb-operator,ocs-operator,local-storage-operator,advanced-cluster-management",
		RegistrySrcPkgFormatted: []string{"kubernetes-nmstate-operator", "metallb-operator ocs-operator", "local-storage-operator", "advanced-cluster-management"},
		RegistryExtraImages:     []string{"quay.io/jparrill/registry:3", "registry.access.redhat.com/rhscl/httpd-24-rhel7:latest", "quay.io/ztpfw/ui:latest"},
		RegistryUser:            "dummy",
		RegistryPass:            "dummy123",

		/*KubeframeNS:                 "kubeframe",
		MarketNS:                    "openshift-marketplace",
		OcDisCatalog:                "kubeframe-catalog",
		OcpReleaseFull:              config.Ztp.Config.OcOCPVersion + ".0",
		RegistryNS:                  "kubeframe-registry",
		RegistryRoute:               "",
		RegistryConfigFile:          "config.yml",
		RegistryOCPDestIndexNS:      "ocp4/openshift4",
		RegistryOLMDestIndexNS:      "olm/redhat-operator-index",
		RegistryOCPReleaseImage:     "quay.io/openshift-release-dev/ocp-release:" + config.Ztp.Config.OcOCPTag,
		RegistryOLMSourceIndex:      "registry.redhat.io/redhat/redhat-operator-index:v" + config.Ztp.Config.OcOCPVersion,
		RegistrySecretName:          "auth",
		RegistryConfigMapName:       "registry-conf",
		RegistryDeploymentName:      "kubeframe-registry",
		RegistryServiceName:         "kubeframe-registry",
		RegistryRouteName:           "kubeframe-registry",
		RegistryPVCName:             "data-pvc",
		RegistryDataMountPath:       "/var/lib/registry",
		RegistryCertMountPath:       "/certs",
		RegistryCertPath:            "/etc/pki/ca-trust/source/anchors",
		RegistryAutoSecretMountPath: "/auth",
		RegistryConfMountPath:       "/etc/docker/registry",
		RegistryPVMode:              "Filesystem",
		RegistryCaCertData:          []byte(""),
		RegistryPathCaCert:          "",
	}
}
*/

//Func Login to log into the new registry
func (r *Registry) Login(ctx context.Context) error {
	args := []string{r.RegistryURL}
	loginOpts := a.LoginOptions{
		AuthFile: r.PullSecretTempFile,
		//CertDir:       r.RegistryPathCaCert,
		CertDir:       r.RegistryCertPath,
		Password:      r.RegistryPass,
		Username:      r.RegistryUser,
		StdinPassword: false,
		GetLoginSet:   false,
		//Verbose:                   false,
		//AcceptRepositories:        true,
		Stdin:                     os.Stdin,
		Stdout:                    os.Stdout,
		AcceptUnspecifiedRegistry: true,
	}
	sysCtx := &types.SystemContext{
		AuthFilePath:                loginOpts.AuthFile,
		DockerCertPath:              loginOpts.CertDir,
		DockerInsecureSkipTLSVerify: types.NewOptionalBool(true),
	}
	return a.Login(ctx, sysCtx, &loginOpts, args)
}

func (r *Registry) CreateCatalogSource(ctx context.Context) error {
	//TODO create if not exists
	log.Println(color.InYellow(">>>> Creating catalog source."))
	olmclient := auth.NewZTPAuth(r.Kubeconfig).GetOlmAuth()

	catalogSource := &v1alpha1.CatalogSource{
		ObjectMeta: metav1.ObjectMeta{
			Name:      r.OwnCatalogName,
			Namespace: r.MarketplaceNS,
		},
		Spec: v1alpha1.CatalogSourceSpec{
			SourceType:  v1alpha1.SourceTypeGrpc,
			Image:       r.RegistryURL + "/" + r.RegistryOLMDestIndexNS + ":v" + strings.Join(strings.Split(r.RegistryOCPRelease, ".")[:2], "."),
			DisplayName: r.OwnCatalogName,
			Publisher:   r.OwnCatalogName,
			UpdateStrategy: &v1alpha1.UpdateStrategy{
				&v1alpha1.RegistryPoll{
					Interval: &metav1.Duration{Duration: time.Minute * 30},
				},
			},
		},
	}
	//create catalog source
	_, err := olmclient.CatalogSources(r.MarketplaceNS).Create(ctx, catalogSource, metav1.CreateOptions{})
	if err != nil {
		log.Printf(color.InRed(">>>> [ERROR] Error creating catalog source: %s"), err.Error())
		return err
	}

	return nil
}

func (r *Registry) GetPullSecretBase() string {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	//get client from kubeconfig extracted based on Mode (HUB or SPOKE)
	client := auth.NewZTPAuth(r.Kubeconfig).GetAuth()
	res, err := client.CoreV1().Secrets(r.PullSecretNS).Get(ctx, r.PullSecretName, metav1.GetOptions{})
	if err != nil {
		return ""
	}
	return string(res.Data[".dockerconfigjson"])
}

//Func to write the content of string to a temporal file
func (r *Registry) WritePullSecretBaseToTempFile(data string) error {
	err := ioutil.WriteFile(r.PullSecretTempFile, []byte(data), 0644)
	if err != nil {
		return err
	}
	// Defer done in the cmd cobra command in order to be available during the cmd execution and remove after program closed
	//defer os.Remove("/tmp/pull-secret-temp.json")
	return nil
}
