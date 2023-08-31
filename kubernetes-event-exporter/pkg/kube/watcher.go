package kube

import (
	"time"

	"github.com/resmoio/kubernetes-event-exporter/pkg/metrics"
	"github.com/rs/zerolog/log"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

var startUpTime = time.Now()

type EventHandler func(event *EnhancedEvent)

type EventWatcher struct {
	informer           cache.SharedInformer
	stopper            chan struct{}
	labelCache         *LabelCache
	omitLookup         bool
	annotationCache    *AnnotationCache
	fn                 EventHandler
	maxEventAgeSeconds time.Duration
	metricsStore       *metrics.Store
}

// 初始化event watcher
func NewEventWatcher(config *rest.Config, namespace string, MaxEventAgeSeconds int64, metricsStore *metrics.Store, fn EventHandler, omitLookup bool) *EventWatcher {
	clientset := kubernetes.NewForConfigOrDie(config)
	// 生成一个factory
	factory := informers.NewSharedInformerFactoryWithOptions(clientset, 0, informers.WithNamespace(namespace))

	// 生成一个event informer(原生的)
	informer := factory.Core().V1().Events().Informer()

	// 对watcher进行了一次封装，事件用的是原始的
	watcher := &EventWatcher{
		informer:           informer,
		stopper:            make(chan struct{}),
		labelCache:         NewLabelCache(config),
		omitLookup:         omitLookup,
		annotationCache:    NewAnnotationCache(config),
		fn:                 fn,
		maxEventAgeSeconds: time.Second * time.Duration(MaxEventAgeSeconds),
		metricsStore:       metricsStore,
	}

	// 添加事件处理
	informer.AddEventHandler(watcher)

	informer.SetWatchErrorHandler(func(r *cache.Reflector, err error) {
		watcher.metricsStore.WatchErrors.Inc()
	})

	return watcher
}

func (e *EventWatcher) OnAdd(obj interface{}) {
	// 获取k8s 原生event obj
	event := obj.(*corev1.Event)
	e.onEvent(event)
}

func (e *EventWatcher) OnUpdate(oldObj, newObj interface{}) {
	// Ignore updates
}

// Ignore events older than the maxEventAgeSeconds
func (e *EventWatcher) isEventDiscarded(event *corev1.Event) bool {
	timestamp := event.LastTimestamp.Time
	if timestamp.IsZero() {
		timestamp = event.EventTime.Time
	}
	eventAge := time.Since(timestamp)
	if eventAge > e.maxEventAgeSeconds {
		// Log discarded events if they were created after the watcher started
		// (to suppres warnings from initial synchrnization)
		if timestamp.After(startUpTime) {
			log.Warn().
				Str("event age", eventAge.String()).
				Str("event namespace", event.Namespace).
				Str("event name", event.Name).
				Msg("Event discarded as being older then maxEventAgeSeconds")
			e.metricsStore.EventsDiscarded.Inc()
		}
		return true
	}
	return false
}

// 对k8s原生event进行处理(这里似乎有版本区分了,event在1.27好像不属于core/v1下面了)
func (e *EventWatcher) onEvent(event *corev1.Event) {
	// 对部分事件丢弃
	if e.isEventDiscarded(event) {
		return
	}

	log.Debug().
		Str("msg", event.Message).
		Str("namespace", event.Namespace).
		Str("reason", event.Reason).
		Str("involvedObject", event.InvolvedObject.Name).
		Msg("Received event")

	e.metricsStore.EventsProcessed.Inc()

	ev := &EnhancedEvent{
		//封装事件
		Event: *event.DeepCopy(),
	}
	ev.Event.ManagedFields = nil

	if e.omitLookup {
		ev.InvolvedObject.ObjectReference = *event.InvolvedObject.DeepCopy()
	} else {
		labels, err := e.labelCache.GetLabelsWithCache(&event.InvolvedObject)
		if err != nil {
			if ev.InvolvedObject.Kind != "CustomResourceDefinition" {
				log.Error().Err(err).Msg("Cannot list labels of the object")
			} else {
				log.Debug().Err(err).Msg("Cannot list labels of the object (CRD)")
			}
			// Ignoring error, but log it anyways
		} else {
			ev.InvolvedObject.Labels = labels
			ev.InvolvedObject.ObjectReference = *event.InvolvedObject.DeepCopy()
		}

		annotations, err := e.annotationCache.GetAnnotationsWithCache(&event.InvolvedObject)
		if err != nil {
			if ev.InvolvedObject.Kind != "CustomResourceDefinition" {
				log.Error().Err(err).Msg("Cannot list annotations of the object")
			} else {
				log.Debug().Err(err).Msg("Cannot list annotations of the object (CRD)")
			}
		} else {
			ev.InvolvedObject.Annotations = annotations
			ev.InvolvedObject.ObjectReference = *event.InvolvedObject.DeepCopy()
		}
	}

	// 将事件丢给engine event处理了
	// func (e *Engine) OnEvent(event *kube.EnhancedEvent
	e.fn(ev)
}

func (e *EventWatcher) OnDelete(obj interface{}) {
	// Ignore deletes
}

// start informer
func (e *EventWatcher) Start() {
	go e.informer.Run(e.stopper)
}

func (e *EventWatcher) Stop() {
	e.stopper <- struct{}{}
	close(e.stopper)
}

func NewMockEventWatcher(MaxEventAgeSeconds int64, metricsStore *metrics.Store) *EventWatcher {
	watcher := &EventWatcher{
		labelCache:         NewMockLabelCache(),
		annotationCache:    NewMockAnnotationCache(),
		maxEventAgeSeconds: time.Second * time.Duration(MaxEventAgeSeconds),
		fn:                 func(event *EnhancedEvent) {},
		metricsStore:       metricsStore,
	}
	return watcher
}

func (e *EventWatcher) setStartUpTime(time time.Time) {
	startUpTime = time
}
