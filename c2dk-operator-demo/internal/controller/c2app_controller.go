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
	c2cloudv1 "c2dk-operator/api/v1"
	"c2dk-operator/internal/resources"
	"context"
	"encoding/json"
	"errors"
	corev1 "k8s.io/api/core/v1"
	errorsv2 "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog/v2"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sort"
	"time"
)

// C2appReconciler reconciles a C2app object
type C2appReconciler struct {
	client.Client
	Scheme        *runtime.Scheme
	EventRecorder record.EventRecorder // 为crd产生事件
}

//+kubebuilder:rbac:groups=c2cloud.c2cloud.cn,resources=c2apps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=c2cloud.c2cloud.cn,resources=c2apps/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=c2cloud.c2cloud.cn,resources=c2apps/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the C2app object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.4/pkg/reconcile

func (r *C2appReconciler) isExist(c2app *c2cloudv1.C2app, ctx context.Context, req ctrl.Request) (error, bool) {
	err := r.Client.Get(ctx, req.NamespacedName, c2app)
	if err != nil {
		// crd存在
		if client.IgnoreNotFound(err) != nil {
			klog.Error("C2app CRD存在但是获取失败")
			return errors.New("C2app CRD存在但是获取失败"), true
		}
		// crd 不存在
		return nil, false
	}
	return nil, true
}

func (r *C2appReconciler) iSUpdate(C2app *c2cloudv1.C2app, ctx context.Context, req ctrl.Request) bool {
	// 获取crd状态，检查其是否发生变化，如果未发生变化就不需要更新svc/deployment
	// 更新逻辑(yaml 是否发生变化(将c2app的spec全部封装到 c2app.metadata.annotations中))
	oldSpec := c2cloudv1.C2appSpec{}
	if C2app.Annotations["spec"] != "" {
		if err := json.Unmarshal([]byte(C2app.Annotations["spec"]), &oldSpec); err != nil {
			klog.Error("oldSpec 序列化失败")
			return false
		}
		if reflect.DeepEqual(C2app.Spec, oldSpec) {
			//klog.Info("C2app CRD", C2app.Name, "未发生变化,忽略更新")
			return false
		}
		klog.Info("C2app CRD", C2app.Name, "内容发生变化，即将更新Deployment/Service/ConfigMap/Storage")
		return true
	}
	return true
}

// 创建deployment
func (r *C2appReconciler) createDeployment(app c2cloudv1.C2AppSpec, C2app *c2cloudv1.C2app, ctx context.Context, req ctrl.Request) error {
	// 创建deployment yaml
	deployment := resources.NewDeploymentYaml(app, C2app)
	objectKey := client.ObjectKey{
		Name:      app.Name,
		Namespace: app.Namespace,
	}
	if err := r.Client.Get(ctx, objectKey, &deployment); err != nil {
		//不存在但是获取失败
		if !errorsv2.IsNotFound(err) {
			klog.Error("deployment/", app.Name, "获取失败", err)
			return err
		}

		if err = r.Client.Create(ctx, &deployment); err != nil {
			klog.Error("deployment/", app.Name, "创建失败", err)
			return err
		}
		klog.Info("deployment/", app.Name, "创建成功")
	} else {
		// 本地已存在deployment，更新
		newDeployment := resources.NewDeploymentYaml(app, C2app)
		if err = r.Client.Update(ctx, &newDeployment); err != nil {
			klog.Info("deployment/", app.Name, "更新失败", err)
			return err
		}
		klog.Info("deployment/", app.Name, "更新成功")
	}

	// 对创建的deployment的状态进行判断
	count := 0
	for {
		time.Sleep(time.Second * 1)
		count++
		if err := r.Client.Get(ctx, objectKey, &deployment); err != nil {
			if !errorsv2.IsNotFound(err) {
				klog.Error("deployment/", app.Name, "状态获取失败", err)
				return err
			}
		}
		if count >= 120 {
			klog.Error("deployment/", app.Name, "状态异常部署失败")
			return errors.New("状态异常部署失败")
		}

		if deployment.Status.AvailableReplicas == app.Replicas {
			klog.Info("deployment/", app.Name, "运行正常")
			break
		}
	}
	return nil
}

// 创建svc
func (r *C2appReconciler) createService(app c2cloudv1.C2AppSpec, C2app *c2cloudv1.C2app, ctx context.Context, req ctrl.Request) error {
	service := resources.NewServiceYaml(app, C2app)
	objectKey := client.ObjectKey{
		Name:      app.Name,
		Namespace: app.Namespace,
	}
	if err := r.Client.Get(ctx, objectKey, &service); err != nil {
		//不存在但是获取失败
		if !errorsv2.IsNotFound(err) {
			klog.Error("service/", app.Name, "获取失败", err)
			return err
		}

		if err = r.Client.Create(ctx, &service); err != nil {
			klog.Error("service/", app.Name, "创建失败", err)
			return err
		}
		klog.Info("service/", app.Name, "创建成功")
	} else {
		// 本地已存在svc，更新
		newService := resources.NewServiceYaml(app, C2app)
		if err = r.Client.Update(ctx, &newService); err != nil {
			klog.Info("service/", app.Name, "更新失败", err)
			return err
		}
		klog.Info("service/", app.Name, "更新成功")
	}
	return nil
}

// 创建configmap
func (r *C2appReconciler) createConfigMap(app c2cloudv1.C2AppSpec, C2app *c2cloudv1.C2app, ctx context.Context, req ctrl.Request) error {
	for _, cm := range app.ConfigMap {
		configMap := resources.NewConfigMapYaml(cm, app.Namespace, C2app)
		objectKey := client.ObjectKey{
			Name:      cm.Name,
			Namespace: app.Namespace,
		}
		if err := r.Client.Get(ctx, objectKey, &configMap); err != nil {
			//不存在但是获取失败
			if !errorsv2.IsNotFound(err) {
				klog.Error("configmap/", cm.Name, "获取失败", err)
				return err
			}
			// 创建
			if err = r.Client.Create(ctx, &configMap); err != nil {
				klog.Error("configmap/", cm.Name, "创建失败", err)
				return err
			}
			klog.Info("configmap/", cm.Name, "创建成功")
		} else {
			newConfig := resources.NewConfigMapYaml(cm, app.Namespace, C2app)
			if err = r.Client.Update(ctx, &newConfig); err != nil {
				klog.Info("configmap/", cm.Name, "更新失败", err)
				return err
			}
			klog.Info("configmap/", cm.Name, "更新成功")
		}
	}
	return nil
}

// 创建pvc
func (r *C2appReconciler) createPvc(app c2cloudv1.C2AppSpec, C2app *c2cloudv1.C2app, ctx context.Context, req ctrl.Request) error {
	for _, pvc := range app.Storage {
		persistentVolumeClaim := resources.NewPersistentVolumeClaimYaml(pvc, app.Namespace, C2app)
		objectKey := client.ObjectKey{
			Name:      pvc.PvcName,
			Namespace: app.Namespace,
		}
		if err := r.Client.Get(ctx, objectKey, &persistentVolumeClaim); err != nil {
			//不存在但是获取失败
			if !errorsv2.IsNotFound(err) {
				klog.Error("pvc/", pvc.PvcName, "获取失败", err)
				return err
			}
			// 创建
			if err = r.Client.Create(ctx, &persistentVolumeClaim); err != nil {
				klog.Error("pvc/", pvc.PvcName, "创建失败", err)
				return err
			}
			klog.Info("pvc/", pvc.PvcName, "创建成功")
		} else {
			// pvc 已存在，无法更新
			klog.Error("pvc/", pvc.PvcName, "已存在,无法更新")
			return nil
		}
	}
	return nil
}

// 创建secret
func (r *C2appReconciler) createSecret(app c2cloudv1.C2AppSpec, C2app *c2cloudv1.C2app, ctx context.Context, req ctrl.Request) error {
	for _, se := range app.Secret {
		Secret := resources.NewSecretYaml(se, app.Namespace, C2app)
		objectKey := client.ObjectKey{
			Name:      se.Name,
			Namespace: app.Namespace,
		}
		if err := r.Client.Get(ctx, objectKey, &Secret); err != nil {
			if !errorsv2.IsNotFound(err) {
				klog.Error("secret/", se.Name, "获取失败", err)
				return err
			}
			if err = r.Client.Create(ctx, &Secret); err != nil {
				klog.Error("secret/", se.Name, "创建失败", err)
				return err
			}
			klog.Info("secret/", se.Name, "创建成功")
		} else {
			newSecret := resources.NewSecretYaml(se, app.Namespace, C2app)
			if err = r.Client.Update(ctx, &newSecret); err != nil {
				klog.Info("secret/", se.Name, "更新失败", err)
				return err
			}
			klog.Info("secret/", se.Name, "更新成功")
		}
	}
	return nil
}

func (r *C2appReconciler) createStatus(C2app *c2cloudv1.C2app, ctx context.Context, name, status, message, phase string) error {
	if phase != "" {
		C2app.Status.Phase = phase
	}

	C2app.Status.Condition = append(C2app.Status.Condition, c2cloudv1.C2AppStatus{
		Name:    name,
		Status:  status,
		Message: message,
	})
	if err := r.Status().Update(ctx, C2app); err != nil {
		return err
	}
	return nil
}

func (r *C2appReconciler) Operation(C2app *c2cloudv1.C2app, ctx context.Context, req ctrl.Request) error {
	// 检查是否需要删除这个CRD
	if C2app.Spec.Operation == "delete" {
		klog.Info("C2app CRD", C2app.Name, "即将被删除")
		if err := r.Client.Delete(ctx, C2app); err != nil {
			klog.Info("C2app CRD删除", C2app.Name, "删除失败")
			return err
		} else {
			klog.Info("C2app CRD删除", C2app.Name, "删除成功")
			return nil
		}
	}
	klog.Info("监测到集群内生成C2app CRD ", C2app.Name)

	// 不需要更新资源,crd未发生变化
	if !r.iSUpdate(C2app, ctx, req) {
		return nil
	}

	// 对资源优先级调整，按从大到小排序
	sort.Sort(C2app.Spec.C2AppList)

	// 更新状态
	C2app.Status.Name = C2app.Name
	C2app.Status.Phase = "Creating"
	if err := r.Status().Update(ctx, C2app); err != nil {
		return err
	}

	// 创建相关的资源
	for _, app := range C2app.Spec.C2AppList {
		// 创建pvc,仅支持nfs-provider
		if len(app.Storage) != 0 {
			err := r.createPvc(app, C2app, ctx, req)
			if err != nil {
				if er := r.createStatus(C2app, ctx, app.Name+"/pvc", "Failed", app.Name+"/pvc资源创建失败", "Failed"); err != nil {
					return er
				}
				return err
			}
			if er := r.createStatus(C2app, ctx, app.Name+"/pvc", "Success", app.Name+"/pvc资源创建成功", ""); err != nil {
				return er
			}
		}

		// 创建svc,调整svc label
		if !reflect.DeepEqual(corev1.ServiceSpec{}, app.ServiceSpec) { //判断是否为空
			app.ServiceSpec.Selector = app.Labels
			err := r.createService(app, C2app, ctx, req)
			if err != nil {
				if er := r.createStatus(C2app, ctx, app.Name+"/svc", "Failed", app.Name+"/svc资源创建失败", "Failed"); err != nil {
					return er
				}
				return err
			}
			if er := r.createStatus(C2app, ctx, app.Name+"/svc", "Success", app.Name+"/svc资源创建成功", "Failed"); err != nil {
				return er
			}
		}

		// 创建configmap
		if len(app.ConfigMap) != 0 {
			err := r.createConfigMap(app, C2app, ctx, req)
			if err != nil {
				if er := r.createStatus(C2app, ctx, app.Name+"/configmap", "Failed", app.Name+"/configmap资源创建失败", "Failed"); err != nil {
					return er
				}
				return err
			}
			if er := r.createStatus(C2app, ctx, app.Name+"/configmap", "Success", app.Name+"/configmap资源创建成功", ""); err != nil {
				return er
			}
		}

		// 创建secret
		if len(app.Secret) != 0 {
			err := r.createSecret(app, C2app, ctx, req)
			if err != nil {
				if er := r.createStatus(C2app, ctx, app.Name+"/secret", "Failed", app.Name+"/secret资源创建失败", "Failed"); err != nil {
					return er
				}
				return err
			}
			if er := r.createStatus(C2app, ctx, app.Name+"/secret", "Success", app.Name+"/secret资源创建成功", ""); err != nil {
				return er
			}
		}

		//创建deployment
		err := r.createDeployment(app, C2app, ctx, req)
		if err != nil {
			if er := r.createStatus(C2app, ctx, app.Name+"/deployment", "Failed", app.Name+"/deployment资源创建失败", "Failed"); err != nil {
				return er
			}
			return err
		}
		if er := r.createStatus(C2app, ctx, app.Name+"/deployment", "Success", app.Name+"/deployment资源创建成功", ""); err != nil {
			return er
		}
		//r.EventRecorder.Eventf(C2app, corev1.EventTypeNormal, "Create", "C2app CRD/%s resource create success", app.Name)
	}

	// 写入annotation，用于对比crd是否发生变化了
	data, _ := json.Marshal(C2app.Spec)
	if C2app.Annotations != nil {
		C2app.Annotations["spec"] = string(data)
	} else {
		C2app.Annotations = map[string]string{
			"spec": string(data),
		}
	}
	if err := r.Client.Update(ctx, C2app); err != nil {
		klog.Error("C2app CRD", C2app.Name, "Annotations更新失败")
		return err
	}
	klog.Info("C2app CRD", C2app.Name, "Annnotations更新成功")

	// 所有资源完成创建,更新状态
	if err := r.createStatus(C2app, ctx, C2app.Name, "Success", C2app.Name+"相关资源创建完成", "Success"); err != nil {
		return err
	}

	return nil
}

func (r *C2appReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// TODO(user): your logic here
	var C2app c2cloudv1.C2app

	err, flag := r.isExist(&C2app, ctx, req)
	// crd获取失败
	if err != nil {
		return ctrl.Result{}, err
	}

	// crd 资源不存在
	if err == nil && flag == false {
		return ctrl.Result{}, nil
	}

	// 当前crd资源正处在被删除阶段
	if C2app.DeletionTimestamp != nil {
		klog.Info("C2app CRD ", C2app.Name, " 正被删除")
		return ctrl.Result{}, nil
	}

	// 根据CRD操作资源
	if err := r.Operation(&C2app, ctx, req); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *C2appReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&c2cloudv1.C2app{}).
		Complete(r)
}
