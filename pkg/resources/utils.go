package resources

import (
	"context"
	"github.com/TwiN/go-color"
	apiroutev1 "github.com/openshift/api/route/v1"
	routev1 "github.com/openshift/client-go/route/clientset/versioned/typed/route/v1"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"log"
	"time"
)

//func to wait for all replicas in a deployment to be ready
func WaitForDeployment(ctx context.Context, deployment *appsv1.Deployment, client *kubernetes.Clientset) error {
	err := wait.Poll(time.Second, time.Minute, func() (bool, error) {
		deployment, err := client.AppsV1().Deployments(deployment.Namespace).Get(ctx, deployment.Name, metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		if deployment.Status.ReadyReplicas == deployment.Status.Replicas {
			return true, nil
		}
		return false, nil
	})
	return err
}

//Func to wait for route to be ready
func WaitForRoute(ctx context.Context, client *routev1.RouteV1Client, route *apiroutev1.Route) error {
	err := wait.Poll(time.Second, time.Minute, func() (bool, error) {
		res, err := client.Routes(route.Namespace).Get(ctx, route.Name, metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		if res.Status.Ingress[0].Host == route.Spec.Host {
			return true, nil
		}
		return false, nil
	})
	return err
}

//Func to wait for service to be ready
func WaitForService(ctx context.Context, client *kubernetes.Clientset, service *v1.Service) error {
	err := wait.Poll(time.Second, time.Minute, func() (bool, error) {
		res, err := client.CoreV1().Services(service.Namespace).Get(ctx, service.Name, metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		if res.Spec.ClusterIP == service.Spec.ClusterIP {
			return true, nil
		}
		return false, nil
	})
	return err
}

//Func to wait for PVC to be ready
func WaitForPVC(ctx context.Context, client *kubernetes.Clientset, pvc *v1.PersistentVolumeClaim) error {
	err := wait.Poll(time.Second, time.Minute, func() (bool, error) {
		res, err := client.CoreV1().PersistentVolumeClaims(pvc.Namespace).Get(ctx, pvc.Name, metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		if res.Status.Phase == v1.ClaimBound {
			return true, nil
		}
		return false, nil
	})
	return err
}

func Retry(attempts int, sleep time.Duration, f func() error) (err error) {
	for i := 0; i < attempts; i++ {
		if i > 0 {
			log.Printf(color.InYellow("[INFO] Doing a new attempt. Attempt number: %d"), i)
			time.Sleep(sleep)
			sleep *= 2
		}
		err = f()
		if err == nil {
			return nil
		}
	}
	log.Printf(color.InYellow("after %d attempts, last error: %s"), attempts, err.Error())
	return err
}
