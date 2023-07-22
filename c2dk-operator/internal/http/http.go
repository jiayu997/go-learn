package http

import (
	c2dkv1 "c2dk-operator/api/v1"
	"c2dk-operator/internal/resources"
	"errors"
	"fmt"
	"io"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/klog/v2"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

const (
	API_KEY = "c2"
)

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func initHttpRoute(mgr manager.Manager) {
	// /c2app/ method post
	http.HandleFunc("/c2app", func(writer http.ResponseWriter, request *http.Request) {
		createC2app(writer, request, mgr)
	})

	// /c2app/status?name=c2app-sample
	http.HandleFunc("/c2app/status", func(writer http.ResponseWriter, request *http.Request) {
		getC2appStatus(writer, request, mgr)
	})

	// /c2app/deployment/status?name=busybox-mysql&namespace=default
	http.HandleFunc("/c2app/deployment/status", func(writer http.ResponseWriter, request *http.Request) {
		queryDeploymentStatus(writer, request, mgr)
	})

	// /c2app/pod/status?name=pod-name&namespace=default
	http.HandleFunc("/c2app/pod/status", func(writer http.ResponseWriter, request *http.Request) {
		getPodStatus(writer, request, mgr)
	})
}

func InitHttp(mgr manager.Manager, httpPort string) {
	// initHttpRoute
	initHttpRoute(mgr)

	// start http
	klog.Fatal(http.ListenAndServe(httpPort, nil))
}

func httpMethodCompare(r *http.Request, targets []string) bool {
	for _, target := range targets {
		if r.Method == target {
			return true
		}
	}
	return false
}

// get http headers and check x-api-key value
func httpHeadAPIKeyCheck(r *http.Request) bool {
	//apiKeyValue := r.Header.Get("x-api-key")
	apiKeyValues := r.Header.Values("x-api-key")
	for _, keyValue := range apiKeyValues {
		if keyValue == API_KEY {
			return true
		}
	}
	return false
}

// sendJson data
func sendJsonResponse(w http.ResponseWriter, code int, message string) {
	jsonData, _ := json.Marshal(Response{
		Code:    code,
		Message: message,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = w.Write(jsonData)
}

func httpCheck(w http.ResponseWriter, r *http.Request, methods []string) (int, error) {
	// http Method check
	if !httpMethodCompare(r, methods) {
		//sendJsonResponse(w, http.StatusMethodNotAllowed, )
		return http.StatusMethodNotAllowed, errors.New(fmt.Sprintf("%s method is not allowed", r.Method))
	}

	// http apiKey check
	if !httpHeadAPIKeyCheck(r) {
		//sendJsonResponse(w, http.StatusForbidden, )
		return http.StatusForbidden, errors.New(fmt.Sprintf("x-api-key: %s is not allowed", r.Header.Get("x-api-key")))
	}
	return 0, nil
}

// /c2app/deployment/status?name=busybox-mysql&namespace=default
func queryDeploymentStatus(w http.ResponseWriter, r *http.Request, mgr manager.Manager) {
	//klog.Info("/c2app/deployment/status")
	klog.Info(r.Proto, "\t", r.UserAgent(), "\t", r.URL.RequestURI())

	// http check
	code, err := httpCheck(w, r, []string{"GET"})
	if err != nil {
		sendJsonResponse(w, code, err.Error())
		return
	}

	// deployment status check
	deploymentName := r.URL.Query().Get("name")
	deploymentNameSpace := r.URL.Query().Get("namespace")

	if deploymentName == "" || deploymentNameSpace == "" {
		sendJsonResponse(w, http.StatusBadRequest, fmt.Sprintf("url: %s is error, please check your url", r.URL.RequestURI()))
		return
	}

	objectKey := client.ObjectKey{
		Name:      deploymentName,
		Namespace: deploymentNameSpace,
	}

	// todo: need to optimize deploymentStatusquery fixed
	_, err = resources.DeploymentStatusQuery(newClient(mgr), objectKey)
	if err != nil {
		sendJsonResponse(w, http.StatusServiceUnavailable, fmt.Sprintf("namespace/%s -- deployment/%s Not Health", objectKey.Namespace, objectKey.Name))
		return
	}
	sendJsonResponse(w, http.StatusOK, fmt.Sprintf("namespace/%s -- deployment/%s Health", objectKey.Namespace, objectKey.Name))
}

// /c2app/status?name=c2app-sample
func getC2appStatus(w http.ResponseWriter, r *http.Request, mgr manager.Manager) {
	klog.Info(r.Proto, "\t", r.UserAgent(), "\t", r.URL.RequestURI())

	// http check
	code, err := httpCheck(w, r, []string{"GET"})
	if err != nil {
		sendJsonResponse(w, code, err.Error())
		return
	}

	// c2app name
	c2appName := r.URL.Query().Get("name")

	if c2appName == "" {
		sendJsonResponse(w, http.StatusBadRequest, fmt.Sprintf("url: %s is error, please check your url", r.URL.RequestURI()))
		return
	}

	objectKey := client.ObjectKey{Name: c2appName}

	_, err = resources.C2appStatusQuery(newClient(mgr), objectKey)
	if err != nil {
		sendJsonResponse(w, http.StatusServiceUnavailable, fmt.Sprintf("c2app: %s Not Health", objectKey.Name))
		return
	}
	sendJsonResponse(w, http.StatusOK, fmt.Sprintf("c2app: %s Health", objectKey.Name))
}

func getPodStatus(w http.ResponseWriter, r *http.Request, mgr manager.Manager) {
	//klog.Info("/c2app/pod/status")
	klog.Info(r.Proto, "\t", r.UserAgent(), "\t", r.URL.RequestURI())

	// http check
	code, err := httpCheck(w, r, []string{"GET"})
	if err != nil {
		sendJsonResponse(w, code, err.Error())
		return
	}

	// deployment status check
	podName := r.URL.Query().Get("name")
	podNameSpace := r.URL.Query().Get("namespace")

	if podNameSpace == "" || podName == "" {
		sendJsonResponse(w, http.StatusBadRequest, fmt.Sprintf("url: %s is error, please check your url", r.URL.RequestURI()))
		return
	}

	_, err = resources.PodStatusQueryByObjectKey(newClient(mgr), podName, podNameSpace)
	if err != nil {
		sendJsonResponse(w, http.StatusServiceUnavailable, err.Error())
		return
	}
	sendJsonResponse(w, http.StatusOK, fmt.Sprintf("namespace/%s -- pod/%s is health", podNameSpace, podName))
}

// /c2app/ method post for create crd resources
func createC2app(w http.ResponseWriter, r *http.Request, mgr manager.Manager) {
	klog.Info(r.Proto, "\t", r.UserAgent(), "\t", r.URL.RequestURI())

	// http check
	code, err := httpCheck(w, r, []string{"POST"})
	if err != nil {
		sendJsonResponse(w, code, err.Error())
		return
	}

	// get user post crd json
	body, err := io.ReadAll(r.Body)
	if err != nil {
		sendJsonResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// unmarshal get user post json data
	defer r.Body.Close()
	var template Template
	if err := json.Unmarshal(body, &template); err != nil {
		sendJsonResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// get user microservice <-> namespace and set
	query := r.URL.Query()
	for microservice_code, microservice_namespace := range query {
		for index, microservceInfo := range template.MicroServiceList {
			if microservceInfo.MicroServiceCode == microservice_code {
				if microservice_namespace[0] != "" {
					template.MicroServiceList[index].MicroServiceNamespace = microservice_namespace[0]
				} else {
					template.MicroServiceList[index].MicroServiceNamespace = "default"
				}
			}
		}
	}

	// user's json convert to c2app crd resource template
	var c2app c2dkv1.C2app
	// template change to c2app crd resource function
	if err := TemplateConvertToC2app(&c2app, &template); err != nil {
		sendJsonResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	// use c2app to create crd
	if err := resources.CreateOrUpdateC2app(mgr.GetClient(), &c2app); err != nil {
		sendJsonResponse(w, http.StatusServiceUnavailable, err.Error())
		return
	}

	sendJsonResponse(w, http.StatusOK, fmt.Sprintf("%s c2app create success", c2app.Name))
}

func newClient(mgr manager.Manager) client.Client {
	return mgr.GetClient()
}
