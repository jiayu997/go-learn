package resources

import (
	c2dkv1 "c2dk-operator/api/v1"
	"context"
	"errors"
	"fmt"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// create or update c2pp
func CreateOrUpdateC2app(cli client.Client, c2app *c2dkv1.C2app) error {
	var objectKey client.ObjectKey = client.ObjectKeyFromObject(c2app)
	var oldC2app c2dkv1.C2app

	// get old c2app
	err := cli.Get(context.TODO(), objectKey, &oldC2app)
	if err != nil {
		if err := client.IgnoreNotFound(err); err != nil {
			return err
		} else {
			// c2app not exist
			if err := cli.Create(context.TODO(), c2app); err != nil {
				return err
			} else {
				return nil
			}
		}
	}

	// crd resource need to update(already exist)
	//c2app.ResourceVersion = oldC2app.ResourceVersion     # must add resourceversion otherwise update c2app crd resource failed
	c2app.ObjectMeta = oldC2app.ObjectMeta

	if err := cli.Update(context.TODO(), c2app); err != nil {
		return err
	} else {
		return nil
	}
}

func C2appStatusQuery(cli client.Client, objectKey client.ObjectKey) ([]c2dkv1.C2ResourceStatus, error) {
	var c2app c2dkv1.C2app

	// query c2app
	if err := cli.Get(context.TODO(), objectKey, &c2app); err != nil {
		if err := client.IgnoreNotFound(err); err != nil {
			return nil, errors.New(fmt.Sprintf("c2app: %s get failed", objectKey.Name))
		} else {
			return nil, errors.New(fmt.Sprintf("c2app: %s not found", objectKey.Name))
		}
	}

	// check resource status
	resourceList := make([]c2dkv1.C2ResourceStatus, 0)
	if c2app.Status.Status == STATUS_SUCCESS {
		return nil, nil
	} else {
		for _, resource := range c2app.Status.ResourceStatus {
			if resource.Status != STATUS_SUCCESS {
				resourceList = append(resourceList, resource)
			}
		}
	}
	return resourceList, nil
}
