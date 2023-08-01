/*
Copyright 2023 jiayu.qu@chinacreator.com.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"c2dk-operator/internal/resources"
	"context"
	json "encoding/json"
	"errors"
	"fmt"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"time"

	c2dkv1 "c2dk-operator/api/v1"
)

// C2appReconciler reconciles a C2app object
type C2appReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

type C2app interface {
	isExist(c2app *c2dkv1.C2app, ctx context.Context, req ctrl.Result) (error, bool)
}

// check c2app exist
func (r *C2appReconciler) isExist(c2app *c2dkv1.C2app, ctx context.Context, req ctrl.Request) (error, bool) {
	if err := r.Client.Get(ctx, req.NamespacedName, c2app); err != nil {
		// crd not exist
		if client.IgnoreNotFound(err) == nil {
			return nil, false
		} else {
			klog.Fatal("c2app exist but fetch failed")
			return err, true
		}
	}
	return nil, true
}

// check crd resource whether need to be updated
func (r *C2appReconciler) isNeedToUpdateByAnnotaion(c2app c2dkv1.C2app, ctx context.Context, req ctrl.Request) bool {
	oldSpec := c2dkv1.C2appSpec{}
	if c2app.Annotations["spec"] != "" {
		if err := json.Unmarshal([]byte(c2app.Annotations["spec"]), &oldSpec); err != nil {
			klog.Error("c2app annotations unmarshal failed")
			return false
		}
		if reflect.DeepEqual(c2app.Spec, oldSpec) {
			klog.Info("c2app not update")
			return false
		} else {
			klog.Info("c2app yaml changed")
			return true
		}
	}
	return true
}

// set crd sub resource status
func (r *C2appReconciler) createOrUpdateSubResourceStatus(c2app *c2dkv1.C2app, resourceStatus c2dkv1.C2ResourceStatus) error {
	var isexist bool
	var index int
	for i, status := range c2app.Status.ResourceStatus {
		if status.Name == resourceStatus.Name && status.Namespace == resourceStatus.Namespace {
			isexist = true
			index = i
		}
	}

	// update status
	if isexist == true {
		c2app.Status.ResourceStatus[index] = resourceStatus
	} else {
		c2app.Status.ResourceStatus = append(c2app.Status.ResourceStatus, resourceStatus)
	}

	// create or update
	if err := r.Status().Update(context.TODO(), c2app); err != nil {
		return err
	}

	return nil
}

// set crd status
func (r *C2appReconciler) createOrUpdateCrdStatus(c2app *c2dkv1.C2app, status string) error {
	c2app.Status.Status = status
	if err := r.Status().Update(context.TODO(), c2app); err != nil {
		return err
	}
	return nil
}

// create configmap resource
func (r *C2appReconciler) createOrUpdateConfigMap(c2app *c2dkv1.C2app) error {
	configmapList, err := resources.GenerateConfigMapByC2app(c2app)
	if err != nil {
		klog.Errorf("c2app: %s crd generate configmap list failed reason: %s", c2app.Name, err.Error())
		return err
	} else {
		klog.Infof("c2app: %s crd generate configmap list success", c2app.Name)
	}

	for _, configmap := range configmapList {
		configmap := configmap
		var resourceStatus c2dkv1.C2ResourceStatus = c2dkv1.C2ResourceStatus{
			Name:      configmap.Name,
			Namespace: configmap.Namespace,
			Type:      "configmap",
			Status:    resources.STATUS_SUCCESS,
		}
		if err := resources.CreateConfigmapWithPolicy(r.Client, &configmap); err != nil {
			resourceStatus.Status = resources.STATUS_SUCCESS
			// update crd status
			_ = r.createOrUpdateSubResourceStatus(c2app, resourceStatus)
			klog.Errorf("c2app: %s crd %s configmap namespace: %s name: %s failed reason: %s", c2app.Name, resources.UPDATE_POLICY, configmap.Namespace, configmap.Name, err.Error())

			return err
		} else {
			// update crd status
			_ = r.createOrUpdateSubResourceStatus(c2app, resourceStatus)
			klog.Infof("c2app: %s crd %s configmap namespace: %s name: %s success", c2app.Name, resources.UPDATE_POLICY, configmap.Namespace, configmap.Name)
		}
	}
	return nil
}

func (r *C2appReconciler) createOrUpdateService(c2app *c2dkv1.C2app) error {
	serviceList, err := resources.GenerateServiceByC2app(c2app)
	if err != nil {
		klog.Errorf("c2app: %s crd generate service list failed reason: %s", c2app.Name, err.Error())
		return err
	} else {
		klog.Infof("c2app: %s crd generate service list success", c2app.Name)
	}

	for _, service := range serviceList {
		service := service
		var resourceStatus c2dkv1.C2ResourceStatus = c2dkv1.C2ResourceStatus{
			Name:      service.Name,
			Namespace: service.Namespace,
			Type:      "service",
			Status:    resources.STATUS_SUCCESS,
		}
		if err := resources.CreateServiceWithPolicy(r.Client, &service); err != nil {
			resourceStatus.Status = resources.STATUS_FAILED
			_ = r.createOrUpdateSubResourceStatus(c2app, resourceStatus)
			klog.Errorf("c2app: %s crd %s service namespace: %s name: %s failed reason: %s", c2app.Name, resources.UPDATE_POLICY, service.Namespace, service.Name, err.Error())
			return err
		} else {
			_ = r.createOrUpdateSubResourceStatus(c2app, resourceStatus)
			klog.Infof("c2app: %s crd %s service namespace: %s name: %s success", c2app.Name, resources.UPDATE_POLICY, service.Namespace, service.Name)
		}
	}
	return nil
}

// create deployment resource
func (r *C2appReconciler) createOrUpdateDeployment(c2app *c2dkv1.C2app) error {
	// todo: need to check deployment status and pods's status
	deploymentList, err := resources.GenerateDeploymentByC2app(c2app)
	//data, _ := jsonf.MarshalIndent(deploymentList, "", "  ")
	//fmt.Println(string(data))
	if err != nil {
		klog.Errorf("c2app: %s crd generate deployment list failed reason: %s", c2app.Name, err.Error())
		return err
	} else {
		klog.Infof("c2app: %s crd generate deployment list success", c2app.Name)
	}

	for _, deployment := range deploymentList {
		deployment := deployment
		var resourceStatus c2dkv1.C2ResourceStatus = c2dkv1.C2ResourceStatus{
			Name:      deployment.Name,
			Namespace: deployment.Namespace,
			Type:      "deployment",
			Status:    resources.STATUS_SUCCESS,
		}
		if err := resources.CreateDeploymentWithPolicy(r.Client, &deployment); err != nil {
			resourceStatus.Status = resources.STATUS_FAILED
			_ = r.createOrUpdateSubResourceStatus(c2app, resourceStatus)
			klog.Errorf("c2app: %s crd %s deployment namespace: %s name: %s failed reason: %s", c2app.Name, resources.UPDATE_POLICY, deployment.Namespace, deployment.Name, err.Error())
			return err
		} else {
			_ = r.createOrUpdateSubResourceStatus(c2app, resourceStatus)
			klog.Infof("c2app: %s crd %s deployment namespace: %s name: %s success", c2app.Name, resources.UPDATE_POLICY, deployment.Namespace, deployment.Name)

			// deployment health check
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(resources.CheckTimeOut))

		FOR:
			for {
				select {
				case <-ctx.Done():
					cancel()
					klog.Errorf("c2app: %s crd %s deployment namespace: %s name: %s health check failed", c2app.Name, resources.UPDATE_POLICY, deployment.Namespace, deployment.Name)
					return errors.New(fmt.Sprintf("current timeout value: %s health check timeout", resources.CheckTimeOut))
				default:
					_, err := resources.DeploymentStatusQuery(r.Client, client.ObjectKeyFromObject(&deployment))
					if err != nil {
						time.Sleep(time.Second * 2)
						//klog.Errorf("c2app: %s crd %s deployment namespace: %s name: %s health check failed will retry", c2app.Name, resources.UPDATE_POLICY, deployment.Namespace, deployment.Name)
					} else {
						klog.Infof("c2app: %s crd %s deployment namespace: %s name: %s health check success", c2app.Name, resources.UPDATE_POLICY, deployment.Namespace, deployment.Name)
						break FOR
					}
				}
			}
			cancel()
		}
	}

	return nil
}

// create secret resources
func (r *C2appReconciler) createOrUpdateSecret(c2app *c2dkv1.C2app) error {
	return nil
}

// create persistent volume
func (r *C2appReconciler) createStorageVolume(c2app *c2dkv1.C2app) error {
	pvcList, err := resources.GenerateStorageByC2app(c2app)
	if err != nil {
		klog.Errorf("c2app: %s crd generate pvc list failed reason: %s", c2app.Name, err.Error())
		return err
	} else {
		klog.Infof("c2app: %s crd generate pvc list success", c2app.Name)
	}

	for _, pvc := range pvcList {
		pvc := pvc
		if err := resources.CreatePvcWithNoPolicy(r.Client, &pvc); err != nil {
			klog.Errorf("c2app: %s crd %s pvc namespace: %s name: %s failed reason: %s", c2app.Name, resources.UPDATE_POLICY, pvc.Namespace, pvc.Name, err.Error())
			return err
		} else {
			klog.Infof("c2app: %s crd %s pvc namespace: %s name: %s success", c2app.Name, resources.UPDATE_POLICY, pvc.Namespace, pvc.Name)
		}
	}
	return nil
}

// create namespace
func (r *C2appReconciler) createNamespace(c2app *c2dkv1.C2app) error {
	namespaceList, err := resources.GenerateNamespaceByC2app(c2app)
	if err != nil {
		klog.Errorf("c2app: %s crd generate namespace list failed reason: %s", c2app.Name, err.Error())
	} else {
		klog.Infof("c2app: %s crd generate namespace list success", c2app.Name)
	}

	for _, namespace := range namespaceList {
		namespace := namespace
		if err := resources.CreateNamespaceWithNoPolicy(r.Client, &namespace); err != nil {
			klog.Errorf("c2app: %s crd %s namespace: %s failed reason: %s", c2app.Name, resources.UPDATE_POLICY, namespace.Name, err.Error())
			return err
		} else {
			klog.Infof("c2app: %s crd %s namespace: %s success", c2app.Name, resources.UPDATE_POLICY, namespace.Name)
		}
	}
	return nil
}

func (r *C2appReconciler) operateResource(c2app *c2dkv1.C2app, req ctrl.Request) error {
	// set crd status failed
	_ = r.createOrUpdateCrdStatus(c2app, resources.STATUS_SUCCESS)

	// create namespace
	if err := r.createNamespace(c2app); err != nil {
		klog.Errorf("namespace %s failed", resources.UPDATE_POLICY)
		return err
	}

	// create or update configmap
	if err := r.createOrUpdateConfigMap(c2app); err != nil {
		klog.Errorf("configmap %s failed", resources.UPDATE_POLICY)
		return err
	}

	// create or update service
	if err := r.createOrUpdateService(c2app); err != nil {
		klog.Errorf("service %s failed", resources.UPDATE_POLICY)
		return err
	}

	// create pvc nfs provide , doesn't support update
	if err := r.createStorageVolume(c2app); err != nil {
		klog.Errorf("pvc %s failed", resources.UPDATE_POLICY)
		return err
	}

	// create or updateDeployment
	if err := r.createOrUpdateDeployment(c2app); err != nil {
		klog.Errorf("deployment %s failed", resources.UPDATE_POLICY)
		return err
	}

	// update crd status ok
	_ = r.createOrUpdateCrdStatus(c2app, resources.STATUS_SUCCESS)
	klog.Infof("c2app: %s crd resource all %s it", c2app.Name, resources.UPDATE_POLICY)

	return nil
}

//+kubebuilder:rbac:groups=c2dk.c2cloud.cn,resources=c2apps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=c2dk.c2cloud.cn,resources=c2apps/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=c2dk.c2cloud.cn,resources=c2apps/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the C2app object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.4/pkg/reconcile
func (r *C2appReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// TODO(user): your logic here
	var C2app c2dkv1.C2app

	err, flag := r.isExist(&C2app, ctx, req)

	// crd fetch failed
	if err != nil {
		klog.Fatal(err, "unable to fetch c2app")
		return ctrl.Result{}, err
	}

	// crd not exist
	if err == nil && flag == false {
		klog.Info("c2app not exist")
		return ctrl.Result{}, nil
	}

	// crd exist

	// crd is being delete
	if C2app.DeletionTimestamp != nil {
		klog.Info("c2app is being delete")
		return ctrl.Result{}, nil
	}

	// operate other resource base on crd resource
	if err := r.operateResource(&C2app, req); err != nil {
		// need to delete crd
		//if err := resources.DeleteC2app(r.Client, &C2app); err != nil {
		//	klog.Errorf("c2app: %s create resource failed and delete failed", C2app.Name)
		//	return ctrl.Result{}, err
		//} else {
		//	klog.Infof("c2app: %s create resource failed and delete success", C2app.Name)
		//	return ctrl.Result{}, nil
		//}

		// don't delete it, please use api to delete it
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *C2appReconciler) SetupWithManager(mgr ctrl.Manager) error {
	// will be not watch other resource because of just create it
	return ctrl.NewControllerManagedBy(mgr).
		For(&c2dkv1.C2app{}). // just watch c2app crd, don't watch other resource update or create
		WithEventFilter(predicate.GenerationChangedPredicate{
			Funcs: predicate.Funcs{
				UpdateFunc: func(e event.UpdateEvent) bool {
					return e.ObjectOld.GetGeneration() != e.ObjectNew.GetGeneration()
				},
			},
		}). // set not watch crd status , fix loop condition
		Complete(r)
}
