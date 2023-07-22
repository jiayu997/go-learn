package http

import (
	c2dkv1 "c2dk-operator/api/v1"
	"c2dk-operator/internal/resources"
	"errors"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// template convert to c2app crd
func TemplateConvertToC2app(c2app *c2dkv1.C2app, template *Template) error {
	if err := microServiceListConvert(c2app, template); err != nil {
		return err
	}
	//data, err := json.MarshalIndent(c2app.Spec.ApplicationList, "", "  ")
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(string(data))
	return nil
}

// change microservice to c2app crd resource template
func microServiceListConvert(c2app *c2dkv1.C2app, template *Template) error {
	c2app.Name = template.TemplateID
	c2app.Spec.ApplicationList = make([]c2dkv1.ApplicationSpec, len(template.MicroServiceList))
	for index, microservice := range template.MicroServiceList {
		var application c2dkv1.ApplicationSpec
		// todo need to calculate priority by depservice

		// c2dk application convert, this func will calculate pvc name, so don't adjust sequence
		if err := c2ApplicationConvert(&application, &microservice); err != nil {
			return err
		}

		// 	todo storage create -> pvc(nfs-provider)|configmap|secret , must after c2applicationConvert
		if err := c2StorageConvert(&application, &microservice); err != nil {
			return err
		}

		// todo service create

		if err := c2ServiceConvert(&application, &microservice); err != nil {
			return err
		}

		// init appliction
		c2app.Spec.ApplicationList[index] = application
	}

	return nil
}

// change microservice to applicationSpec crd
func c2ApplicationConvert(application *c2dkv1.ApplicationSpec, microservice *MicroServiceInfo) error {
	// select deployment/sts/ds, create deployment
	switch microservice.Kind {
	case "deployment":
		if err := deploymentConvert(application, microservice); err != nil {
			return err
		}
		// todo
	case "statefulset":
		if err := statefulsetConvert(application, microservice); err != nil {
			return nil
		}
		// todo
	case "daemonset":
		if err := daemonsetConvert(application, microservice); err != nil {
			return err
		}
	default:
		return errors.New("controller type just support deployment/statefulset/daemonset")
	}
	return nil
}

func deploymentConvert(application *c2dkv1.ApplicationSpec, microservice *MicroServiceInfo) error {
	application.Name = microservice.MicroServiceCode
	application.NameSpace = microservice.MicroServiceNamespace
	application.Labels = microservice.Labels
	application.Annotations = microservice.Annotations
	application.ControllerType = "deployment"
	// todo
	application.Priority = 0
	application.Replicas = microservice.InstanceCount

	if err := podSpecConvert(application, microservice); err != nil {
		return err
	}
	return nil
}
func podSpecConvert(application *c2dkv1.ApplicationSpec, microservice *MicroServiceInfo) error {
	// set pod spec
	// set pod restart policy
	application.PodSpec.RestartPolicy = corev1.RestartPolicyAlways

	// pod containers convert, this func must be first because of configmap's volume name is generate by some way
	if err := podContainerConvert(application, microservice); err != nil {
		return err
	}

	// pod volume convert
	if err := podVolumesConvert(application, microservice); err != nil {
		return err
	}

	return nil
}

// todo: Deployment.apps "c2dk-mysql" is invalid: [spec.template.spec.volumes[1].name: Duplicate value: "c2dk-mysql-my-cnf", spec.template.spec.volumes[3].name: Duplicate value
func podVolumesConvert(application *c2dkv1.ApplicationSpec, microservice *MicroServiceInfo) error {
	// set pod volumes from configmaps
	var volumes []corev1.Volume = make([]corev1.Volume, 0)
	var configDefaultMode = corev1.ConfigMapVolumeSourceDefaultMode
	for _, configInfo := range microservice.Storage.ConfigMapList {
		for key, configContext := range configInfo.ConfigMapContextList {
			volumes = append(volumes, corev1.Volume{
				Name: configInfo.ConfigMapContextList[key].Name,
				VolumeSource: corev1.VolumeSource{
					ConfigMap: &corev1.ConfigMapVolumeSource{
						LocalObjectReference: corev1.LocalObjectReference{Name: configInfo.VolumeName},
						DefaultMode:          &configDefaultMode,
						Items: []corev1.KeyToPath{
							corev1.KeyToPath{
								Key:  configContext.Key,
								Path: configContext.SubPath,
							},
						},
					},
				},
			})
		}
	}

	// set pod volumes from volume
	for _, volumeInfo := range microservice.Storage.VolumeList {
		volumes = append(volumes, corev1.Volume{
			Name: volumeInfo.VolumeName,
			VolumeSource: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: volumeInfo.VolumeName,
					ReadOnly:  false,
				},
			},
		})
	}
	application.PodSpec.Volumes = volumes
	return nil
}

func podContainerConvert(application *c2dkv1.ApplicationSpec, microservice *MicroServiceInfo) error {
	// todo: just use one container of pod, doesn't support multi container

	//var podSpec corev1.PodSpec
	var container corev1.Container
	// set container name
	container.Name = microservice.MicroServiceCode

	// set pod image , default just one container
	container.Image = microservice.Artifact.Image + ":" + microservice.Artifact.ImageVersion
	container.ImagePullPolicy = corev1.PullIfNotPresent

	// set pod command
	if len(microservice.Cmd) != 0 {
		container.Command = microservice.Cmd
	}

	// set pod args
	if len(microservice.Args) != 0 {
		container.Args = microservice.Args
	}

	// set pod container ports
	container.Ports = make([]corev1.ContainerPort, 0)
	for _, portInfo := range microservice.PortMapList {
		container.Ports = append(container.Ports, corev1.ContainerPort{
			Name:          portInfo.PortName,
			ContainerPort: portInfo.ContainerPort,
			Protocol:      corev1.Protocol(portInfo.Protocol),
		})
	}

	// set pod env
	container.Env = make([]corev1.EnvVar, 0)
	for _, envInfo := range microservice.EnvMapList {
		container.Env = append(container.Env, corev1.EnvVar{
			Name:  envInfo.AttrName,
			Value: envInfo.AttrValue,
		})
	}
	for _, connectInfo := range microservice.ConnectInfoMapList {
		container.Env = append(container.Env, corev1.EnvVar{
			Name:  connectInfo.AttrName,
			Value: connectInfo.AttrValue,
		})
	}

	// set pod request resource
	if (microservice.RequireResource != Resource{}) {
		container.Resources.Requests = corev1.ResourceList{
			corev1.ResourceMemory: resource.MustParse(microservice.RequireResource.Memory),
			corev1.ResourceCPU:    resource.MustParse(microservice.RequireResource.Cpu),
		}
	}

	// set pod limit resource
	if (microservice.LimitResource != Resource{}) {
		container.Resources.Limits = corev1.ResourceList{
			corev1.ResourceMemory: resource.MustParse(microservice.LimitResource.Memory),
			corev1.ResourceCPU:    resource.MustParse(microservice.LimitResource.Cpu),
		}
	}

	// set pod volumeMounts from configmaps
	container.VolumeMounts = make([]corev1.VolumeMount, 0)
	for i, configInfo := range microservice.Storage.ConfigMapList {
		for j, configContext := range configInfo.ConfigMapContextList {
			// generate random volume name or mount path name
			microservice.Storage.ConfigMapList[i].ConfigMapContextList[j].Name = fmt.Sprintf("%s-%d-%d", configInfo.VolumeName, i, j)
			container.VolumeMounts = append(container.VolumeMounts, corev1.VolumeMount{
				Name:      microservice.Storage.ConfigMapList[i].ConfigMapContextList[j].Name,
				ReadOnly:  false,
				MountPath: configContext.VolumePath,
				SubPath:   configContext.SubPath,
			})
		}
	}

	// set pod volumeMounts from volumes
	for _, volume := range microservice.Storage.VolumeList {
		container.VolumeMounts = append(container.VolumeMounts, corev1.VolumeMount{
			Name:      volume.VolumeName,
			ReadOnly:  false,
			MountPath: volume.VolumePath,
		})
	}

	// set pod liveness or readyness
	if len(microservice.ProbeList) > 2 {
		return errors.New("probe list length unvalidity")
	}
	if len(microservice.ProbeList) != 0 {
		container.LivenessProbe = new(corev1.Probe)
		container.ReadinessProbe = new(corev1.Probe)
		for _, probeInfo := range microservice.ProbeList {
			probe := new(corev1.Probe)
			probe.SuccessThreshold = probeInfo.SuccessThreshold
			probe.FailureThreshold = probeInfo.FailureThreshold
			probe.TimeoutSeconds = probeInfo.TimeoutSecond
			probe.PeriodSeconds = probeInfo.PeriodSecond
			probe.InitialDelaySeconds = probeInfo.InitialDelaySecond
			if probeInfo.Scheme == "TCP" {
				probe.TCPSocket = new(corev1.TCPSocketAction) // must new
				probe.TCPSocket.Host = "127.0.0.1"
				probe.TCPSocket.Port = intstr.IntOrString{IntVal: probeInfo.Port}
			} else if probeInfo.Scheme == "HTTP" {
				probe.HTTPGet = new(corev1.HTTPGetAction) // must new
				probe.HTTPGet.Path = probeInfo.Path
				probe.HTTPGet.Scheme = corev1.URISchemeHTTP
				probe.HTTPGet.Host = "127.0.0.1"
				probe.HTTPGet.Port = intstr.IntOrString{IntVal: probeInfo.Port}
			}

			// liveness or readiness
			if probeInfo.Mode == "liveness" {
				container.LivenessProbe = probe
			} else if probeInfo.Mode == "readiness" {
				container.ReadinessProbe = probe
			}
		}
	}

	application.PodSpec.Containers = []corev1.Container{container}

	return nil
}

func c2StorageConvert(application *c2dkv1.ApplicationSpec, microserviceInfo *MicroServiceInfo) error {
	// configmap convert
	application.ConfigMapSpec = make([]c2dkv1.ConfigMap, 0)
	for _, configInfo := range microserviceInfo.Storage.ConfigMapList {
		var configmap c2dkv1.ConfigMap
		configmap.Name = configInfo.VolumeName
		configmap.Data = make(map[string]string)
		configmap.BinaryData = make(map[string]string)
		//configmap.NameSpace = microserviceInfo.MicroServiceNamespace
		// set configmap sub key <-> value
		for _, configContext := range configInfo.ConfigMapContextList {
			configmap.Data[configContext.Key] = configContext.FileContent
		}
		application.ConfigMapSpec = append(application.ConfigMapSpec, configmap)
	}

	// todo: secret create

	// pvc convert
	application.StorageSpec = make([]c2dkv1.Storage, 0)
	for _, volume := range microserviceInfo.Storage.VolumeList {
		var storage c2dkv1.Storage
		//storage.NameSpace = microserviceInfo.MicroServiceNamespace
		storage.PvcName = volume.VolumeName
		storage.AccessMode = volume.VolumeAccessMode
		storage.StorageClassName = resources.STORAGE_CLASS_NFS
		application.StorageSpec = append(application.StorageSpec, storage)
	}
	return nil
}

func c2ServiceConvert(application *c2dkv1.ApplicationSpec, microserviceInfo *MicroServiceInfo) error {
	// no service
	if len(microserviceInfo.PortMapList) == 0 {
		return nil
	}
	application.ServiceSpec = make(map[string]c2dkv1.ServiceSpec, 0)
	var innerService, outService c2dkv1.ServiceSpec
	innerService.Name = fmt.Sprintf("%s-inner-svc", microserviceInfo.MicroServiceCode)
	outService.Name = fmt.Sprintf("%s-out-svc", microserviceInfo.MicroServiceCode)
	innerService.Ports = make([]corev1.ServicePort, 0)
	outService.Ports = make([]corev1.ServicePort, 0)
	innerService.Selector = microserviceInfo.Labels
	outService.Selector = microserviceInfo.Labels

	for _, portInfo := range microserviceInfo.PortMapList {
		var port corev1.ServicePort
		port.Name = portInfo.PortName
		port.Port = portInfo.ContainerPort
		port.Protocol = corev1.Protocol(portInfo.Protocol)
		port.TargetPort = intstr.IntOrString{IntVal: portInfo.ContainerPort}
		// out
		if (portInfo.IOuterService && portInfo.IsInnerService) || portInfo.IOuterService == true {
			if portInfo.NodePort != 0 {
				port.NodePort = portInfo.NodePort
			}
			outService.Type = corev1.ServiceTypeNodePort
			outService.Ports = append(outService.Ports, port)
		} else if portInfo.IsInnerService && portInfo.IOuterService == false {
			innerService.Type = corev1.ServiceTypeClusterIP
			innerService.Ports = append(innerService.Ports, port)
		}
	}
	if len(innerService.Ports) != 0 {
		application.ServiceSpec["inner"] = innerService
	}
	if len(outService.Ports) != 0 {
		application.ServiceSpec["out"] = outService
	}

	return nil
}

// todo
func daemonsetConvert(application *c2dkv1.ApplicationSpec, microservice *MicroServiceInfo) error {
	return nil
}

// todo
func statefulsetConvert(application *c2dkv1.ApplicationSpec, microservice *MicroServiceInfo) error {
	return nil
}
