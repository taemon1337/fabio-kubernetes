package kubernetes

import (
	"errors"
	"log"
  "net/http"
  "time"

	"github.com/fabiolb/fabio/config"
	"github.com/fabiolb/fabio/registry"
)

// be is an implementation of a registry backend for kubernetes.
type be struct {
	c     *http.Client
	cfg   *config.Kubernetes
	dereg map[string](chan bool)
}

func NewBackend(cfg *config.Kubernetes) (registry.Backend, error) {
  client := &http.Client{Transport: &http.Transport{cfg.Transport}}

  req := client.NewRequest("GET", cfg.Scheme + cfg.Addr, nil)
  req.Header.Add("Bearer", cfg.Token)

  res, err := client.Do(req)
  if err != nil {
    return nil, err
  }

	// we're good
	log.Printf("[INFO] kubernetes: Connecting to Kubernetes at %s%s", cfs.Scheme, cfg.Addr)
	return &be{c: client, cfg: cfg}, nil
}

func (b *be) Register(services []string) error {
	log.Printf("[INFO] Kubernetes default registration of Fabio is not supported.")
	return nil
}

func (b *be) Deregister(service string) error {
	dereg := b.dereg[service]
	if dereg == nil {
		log.Printf("[WARN]: Attempted to deregister unknown service %q", service)
		return nil
	}
	dereg <- true // trigger deregistration
	<-dereg       // wait for completion
	delete(b.dereg, service)

	return nil
}

func (b *be) DeregisterAll() error {
	log.Printf("[DEBUG]: kubernetes: Deregistering all registered aliases.")
	for name, dereg := range b.dereg {
		if dereg == nil {
			continue
		}
		log.Printf("[INFO] kubernetes: Deregistering %q", name)
		dereg <- true // trigger deregistration
		<-dereg       // wait for completion
	}
	return nil
}

func (b *be) ManualPaths() ([]string, error) {
	log.Printf("[INFO] Kubernetes No implementation for ManualPaths.")
	return [], err
}

func (b *be) ReadManual(path string) (value string, version uint64, err error) {
	log.Printf("[INFO] Kubernetes No implementation for ReadPaths.")
	return (b.c, path, 0)
}

func (b *be) WriteManual(path string, value string, version uint64) (ok bool, err error) {
	log.Printf("[INFO] Kubernetes No implementation for WriteManual.")
  return (b.c, path, 0)
}

func (b *be) WatchServices() chan string {
	log.Printf("[INFO] kubernetes: Watching all Kubernetes services")

  /*
   * Kubernetes Api : http://kubernetes-master/api/v1/services
   *   - request all services from k8 api
   *   - filter by some special label (like 'fabio')
   *   - add fabio route for '{service}.{namespace}.my-domain.com' to '{service}.{namespace}.svc.cluster.local'
   *   - delete any existing routes not found in k8 services (or a complete overwrite of config)
   */

	svc := make(chan string)
	go watchServices(b.c, b.cfg)
	return svc
}

func (b *be) WatchManual() chan string {
	log.Printf("[INFO] kubernetes: No implementation for WatchManual")

	kv := make(chan string)
	return kv
}

func (b *be) WatchNoRouteHTML() chan string {
	log.Printf("[INFO] kubernetes: No implementation for WatchNoRouteHTML")

	html := make(chan string)
	return html
}

