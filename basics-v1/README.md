# client-go basics

![Go Logo](https://upload.wikimedia.org/wikipedia/commons/thumb/0/05/Go_Logo_Blue.svg/320px-Go_Logo_Blue.svg.png)

## Getting Started

client-go is a typical web service client library that supports all API types that are officially part of Kubernetes. 

## Creating and Using a Client

The above code contains the implementation details.

I will be attaching the code snippets for understanding.

```go
kubeconfig = flag.String("kubeconfig", "~/.kube/config", "kubeconfig file")
flag.Parse()
config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
clientset, err := kubernetes.NewForConfig(config)
```
This imports clientcmd from client-go in order to read and parse the kubeconfig. The above code different from the code in the main.go specified above.

- The default location for the kubeconfig file is in .kube/config in the user’s home
directory. This is also where kubectl gets the credentials for the Kubernetes clusters
- kubeconfig is then read and parsed using clientcmd.BuildConfigFromFlags
- From **clientcmd.BuildConfigFromFlags** we get a rest.Config, and this is passed to kubernetes.NewForConfig in order for creating the kubernetes clientset.
- clientset contains multiple clients for all native K8s resources.
- When running a binary inside of a pod in a cluster, the kubelet will automatically
mount a service account into the container at
/var/run/secrets/kubernetes.io/serviceaccount. It replaces the kubeconfig file just
mentioned and can easily be turned into a rest.Config via the
rest.InClusterConfig() method.

## K8s Object
Lets get deep dive into the K8s objects.

Kubernetes objects are instances of a kind and served as resource by API Server which is represented as structs.

Kubernetes objects fulfill a Go interface called
**runtime.Object** from the package k8s.io/apimachinery/pkg/runtime.

```go
// Object interface must be supported by all API types registered with Scheme.
// Since objects in a scheme are expected to be serialized to the wire, the
// interface an Object must provide to the Scheme allows serializers to set
// the kind, version, and group the object is represented as. An Object may
// choose to return a no-op ObjectKindAccessor in cases where it is not
// expected to be serialized.
type Object interface {
GetObjectKind() schema.ObjectKind
DeepCopyObject() Object
}
```
schema.ObjectKind (from the k8s.io/apimachinery/pkg/runtime/schema
package) is another simple interface:
```go
// All objects that are serialized from a Scheme encode their type information.
// This interface is used by serialization to set type information from the
// Scheme onto the serialized version of an object. For objects that cannot
// be serialized or have unique requirements, this interface may be a no-op.
type ObjectKind interface {
// SetGroupVersionKind sets or clears the intended serialized kind of an
// object. Passing kind nil should clear the current setting.
SetGroupVersionKind(kind GroupVersionKind)
// GroupVersionKind returns the stored group, version, and kind of an
// object, or nil if the object does not expose or provide these fields.
GroupVersionKind() GroupVersionKind
}
```

Kubernetes object in Go is a data structure that can:
- Return and set the GroupVersionKind
- Be deep-copied

Deep copy is used wherever code has to mutate an object without
modifying the original.

## TypeMeta
In Kubernetes, when you create or interact with objects like Pods, Services, or Deployments, these objects are represented in Go code using structs. The k8s.io/api package provides these structs that define the structure of these Kubernetes objects.

Now, some common functionalities are needed for these objects, like identifying their type and version. To support this, Kubernetes objects implement the schema.ObjectKind interface, which provides methods to get and set information about the type of the object.

Here's where it gets interesting:

1. The metav1.TypeMeta struct, from the k8s.io/apimachinery/meta/v1 package, contains fields for the type and version information. It's like a blueprint for adding metadata to Kubernetes objects.

2. Instead of duplicating this type and version information in every Kubernetes object struct from k8s.io/api, they use a concept called "embedding." In Go, you can embed one struct into another to reuse its fields and methods.

3. By embedding metav1.TypeMeta into the structs from k8s.io/api, those structs automatically get the fields and methods from TypeMeta. This means that every Kubernetes object, in addition to its specific fields, also has the type and version information thanks to TypeMeta.

4. The schema.ObjectKind interface requires methods like GetObjectKind() and SetObjectKind(obj ObjectKind) to be implemented. Since the structs from k8s.io/api already have the necessary type information due to the embedding of TypeMeta, they fulfill the requirements of the schema.ObjectKind interface without needing to write additional code.

```go
// TypeMeta describes an individual object in an API response or request
// with strings representing the type of the object and its API schema version.
// Structures that are versioned or persisted should inline TypeMeta.
//
// +k8s:deepcopy-gen=false
type TypeMeta struct {
// Kind is a string value representing the REST resource this object
// represents. Servers may infer this from the endpoint the client submits
// requests to.
// Cannot be updated.
// In CamelCase.
// +optional
Kind string `json:"kind,omitempty" protobuf:"bytes,1,opt,name=kind"`
// APIVersion defines the versioned schema of this representation of an
// object. Servers should convert recognized schemas to the latest internal
// value, and may reject unrecognized values.
// +optional
APIVersion string `json:"apiVersion,omitempty"`
}
```

With this, a pod declaration in Go looks like this:
```go
// Pod is a collection of containers that can run on a host. This resource is
// created by clients and scheduled onto hosts.
type Pod struct {
metav1.TypeMeta `json:",inline"`
// Standard object's metadata.
// +optional
metav1.ObjectMeta `json:"metadata,omitempty"`
// Specification of the desired behavior of the pod.
// +optional
Spec PodSpec `json:"spec,omitempty"`
// Most recently observed status of the pod.
// This data may not be up to date.
// Populated by the system.
// Read-only.
// +optional
Status PodStatus `json:"status,omitempty"`
}
```

## ObjectMeta
In Kubernetes, most top-level objects have a field called ObjectMeta. This field contains metadata about the object, and it comes from the k8s.io/apimachinery/pkg/meta/v1 package.

The ObjectMeta struct looks like this:

```go
type ObjectMeta struct {
    Name            string            `json:"name,omitempty"`
    Namespace       string            `json:"namespace,omitempty"`
    UID             types.UID         `json:"uid,omitempty"`
    ResourceVersion string            `json:"resourceVersion,omitempty"`
    CreationTimestamp Time              `json:"creationTimestamp,omitempty"`
    DeletionTimestamp *Time             `json:"deletionTimestamp,omitempty"`
    Labels          map[string]string `json:"labels,omitempty"`
    Annotations     map[string]string `json:"annotations,omitempty"`
    // ... other fields
}
```
In JSON or YAML representation, these fields are placed under the metadata key. For example, for a pod named "example" in the "default" namespace, ObjectMeta stores:
```yaml
metadata:
  namespace: default
  name: example
```
ObjectMeta contains metadata like the object's name, namespace, a unique identifier (UID), resource version, creation timestamp, deletion timestamp (if applicable), labels, annotations, and more.

The resource version is a crucial field in Kubernetes, although it's rarely directly manipulated by client code. It helps manage changes to objects in the system. The resourceVersion is part of ObjectMeta because each object with embedded ObjectMeta corresponds to a key in etcd (a distributed key-value store used by Kubernetes), where the resourceVersion value originated.

## ClientSets
When we use kubernetes.NewForConfig(config) from k8s.io/client-go/kubernetes, we gain access to almost all API groups and resources defined in k8s.io/api. This includes most resources served by the Kubernetes API server, with a few exceptions like APIServices and CustomResourceDefinition.
