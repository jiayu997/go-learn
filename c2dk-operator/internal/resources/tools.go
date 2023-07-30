package resources

import (
	c2dkv1 "c2dk-operator/api/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"time"
)

const (
	CREATE = "CREATE"
	UPDATE = "UPDATE"

	CheckTimeOut      = 90 * time.Second
	STORAGE_CLASS_NFS = "nfs-client"

	STATUS_FAILED  = "failed"
	STATUS_SUCCESS = "success"
)

var UPDATE_POLICY string

func init() {
	UPDATE_POLICY = os.Getenv("UPDATE_POLICY")
	if UPDATE_POLICY != UPDATE && UPDATE_POLICY != CREATE {
		UPDATE_POLICY = UPDATE
	}
}

func NewOwnerReference(c2app *c2dkv1.C2app) []metav1.OwnerReference {
	return []metav1.OwnerReference{
		*metav1.NewControllerRef(c2app, schema.GroupVersionKind{
			Group:   c2dkv1.GroupVersion.Group,
			Version: c2dkv1.GroupVersion.Version,
			Kind:    c2app.Kind,
		}),
	}
}

func removeDuplicateObjectList(objectList []interface{}) ([]interface{}, error) {
	var resultList []interface{}

	if len(objectList) <= 0 {
		return nil, nil
	}

	for i := range objectList {
		flag := true
		for j := range resultList {
			switch objectList[i].(type) {
			case appsv1.Deployment:
			case corev1.Namespace:
				object, _ := objectList[i].(corev1.Namespace)
				result, _ := resultList[j].(corev1.Namespace)
				if object.Name == result.Name {
					flag = false
					break
				}
			case corev1.Service:
				object, _ := objectList[i].(corev1.Service)
				result, _ := resultList[j].(corev1.Service)
				if object.Name == result.Name && object.Namespace == object.Namespace {
					flag = false
					break
				}
			case corev1.ConfigMap:
			case corev1.Secret:
			case corev1.PersistentVolumeClaim:
			}
		}

	}
	return nil, nil
}

// get new client
func NewClient() (client.Client, error) {
	// get client cfg
	// default ~/.kube/config or cluster
	cfg, err := config.GetConfig()
	if err != nil {
		return nil, err
	}

	cli, err := client.New(cfg, client.Options{})

	if err != nil {
		return nil, err
	}
	return cli, nil
}
