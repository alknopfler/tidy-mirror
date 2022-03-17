package registry

import (
	"context"
	"github.com/TwiN/go-color"
	"strings"

	"github.com/alknopfler/tidy-mirror/pkg/resources"

	adm "github.com/openshift/oc/pkg/cli/admin/release"
	"github.com/openshift/oc/pkg/cli/image/manifest"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"time"

	"log"
	"os"
)

func (r *Registry) RunMirrorOcp() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//Login to the registry to grab the authfile with the new registry credentials
	err := r.Login(ctx)
	if err != nil {
		log.Printf(color.InRed(">>>> [ERROR] login to registry: %s"), err.Error())
		return err
	}
	log.Println(color.InGreen(">>>> [INFO] login to registry successful"))

	//Mirror ocp with a retry strategic to avoid errors
	err = resources.Retry(4, 1*time.Minute, func() (err error) {
		return r.mirrorOcp()
	})
	if err != nil {
		log.Printf(color.InRed(">>>> [ERROR] mirroring the OCP image: %s"), err.Error())
		return err
	}
	log.Println(color.InGreen(">>>> [INFO] mirroring the OCP image successful"))
	return nil
}

func (r *Registry) mirrorOcp() error {

	opt := adm.MirrorOptions{
		IOStreams: genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr},
		SecurityOptions: manifest.SecurityOptions{
			RegistryConfig: r.PullSecretTempFile,
		},
		ParallelOptions: manifest.ParallelOptions{
			MaxPerRegistry: 100,
		},
		From:        strings.Join(strings.Split(r.RegistryOCPRelease, ".")[:2], "."),
		To:          r.RegistryURL + "/" + r.RegistryOCPDestIndexNS,
		ToRelease:   r.RegistryURL + "/" + r.RegistryOCPDestIndexNS + ":" + r.RegistryOCPRelease + "-x86_64",
		SkipRelease: false,
		DryRun:      false,
		ImageStream: nil,
		TargetFn:    nil,
	}

	return opt.Run()
}
