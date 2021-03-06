package clientset

import (
	glog "github.com/golang/glog"
	quotav1 "github.com/openshift/origin/pkg/quota/generated/clientset/typed/quota/v1"
	discovery "k8s.io/client-go/discovery"
	rest "k8s.io/client-go/rest"
	flowcontrol "k8s.io/client-go/util/flowcontrol"
)

type Interface interface {
	Discovery() discovery.DiscoveryInterface
	QuotaV1() quotav1.QuotaV1Interface
	// Deprecated: please explicitly pick a version if possible.
	Quota() quotav1.QuotaV1Interface
}

// Clientset contains the clients for groups. Each group has exactly one
// version included in a Clientset.
type Clientset struct {
	*discovery.DiscoveryClient
	*quotav1.QuotaV1Client
}

// QuotaV1 retrieves the QuotaV1Client
func (c *Clientset) QuotaV1() quotav1.QuotaV1Interface {
	if c == nil {
		return nil
	}
	return c.QuotaV1Client
}

// Deprecated: Quota retrieves the default version of QuotaClient.
// Please explicitly pick a version.
func (c *Clientset) Quota() quotav1.QuotaV1Interface {
	if c == nil {
		return nil
	}
	return c.QuotaV1Client
}

// Discovery retrieves the DiscoveryClient
func (c *Clientset) Discovery() discovery.DiscoveryInterface {
	if c == nil {
		return nil
	}
	return c.DiscoveryClient
}

// NewForConfig creates a new Clientset for the given config.
func NewForConfig(c *rest.Config) (*Clientset, error) {
	configShallowCopy := *c
	if configShallowCopy.RateLimiter == nil && configShallowCopy.QPS > 0 {
		configShallowCopy.RateLimiter = flowcontrol.NewTokenBucketRateLimiter(configShallowCopy.QPS, configShallowCopy.Burst)
	}
	var cs Clientset
	var err error
	cs.QuotaV1Client, err = quotav1.NewForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}

	cs.DiscoveryClient, err = discovery.NewDiscoveryClientForConfig(&configShallowCopy)
	if err != nil {
		glog.Errorf("failed to create the DiscoveryClient: %v", err)
		return nil, err
	}
	return &cs, nil
}

// NewForConfigOrDie creates a new Clientset for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *Clientset {
	var cs Clientset
	cs.QuotaV1Client = quotav1.NewForConfigOrDie(c)

	cs.DiscoveryClient = discovery.NewDiscoveryClientForConfigOrDie(c)
	return &cs
}

// New creates a new Clientset for the given RESTClient.
func New(c rest.Interface) *Clientset {
	var cs Clientset
	cs.QuotaV1Client = quotav1.New(c)

	cs.DiscoveryClient = discovery.NewDiscoveryClient(c)
	return &cs
}
