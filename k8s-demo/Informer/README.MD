## Informer创建
client-go 中提供了几种不同的 Informer：
- 通过调用 NewInformer 函数创建一个简单的不带 indexer 的 Informer。
- 通过调用 NewIndexerInformer 函数创建一个简单的带 indexer 的 Informer。
- 通过调用 NewSharedIndexInformer 函数创建一个 Shared 的 Informer。
- 通过调用 NewDynamicSharedInformerFactory 函数创建一个为 Dynamic 客户端的 Shared 的 Informer

## SharedInformer
- 我们平时说的 Informer 其实就是 SharedInformer，它是可以共享使用的。如果同一个资源的 Informer 被实例化多次，那么就会运行多个 ListAndWatch 操作，这会加大 APIServer 的压力。而 SharedInformer 通过一个 map 来让同一类资源的 Informer 实现共享一个 Refelctor，这样就不会出现上面这个问题了
- Informer通过Local Store缓存目标资源对象，且仅为自己所用。但是在K8s中，一个Controller可以关心不止一种资源，使得多个Controller所关心的资源彼此会存在交集。如果几个Controller都用自己的Informer来缓存同一个目标资源，显然会导致不小的空间开销，因此K8s引入了SharedInformer来解决这个问题
- SharedInformer拥有为多个Controller提供一个共享cache的能力，从而避免资源缓存的重复、减小空间开销。除此之外，一个SharedInformer对一种资源只建立一个与API Server的Watch监听，且能够将监听得到的事件分发给下游所有感兴趣的Controller，这也显著地减少了API Server的负载压力。实际上，K8s中广泛使用的都是SharedInformer，Informer则出场甚少。