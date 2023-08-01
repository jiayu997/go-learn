package resources

import (
	c2dkv1 "c2dk-operator/api/v1"
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
