package di

import (
	"fmt"
	"reflect"
	"sync"
)

// Container manages all services and their dependencies
type Container struct {
	mu            sync.RWMutex
	services      map[string]any
	bootstrappers []func()
}

// NewContainer creates a new dependency injection container
func NewContainer(size int) *Container {
	return &Container{
		services:      make(map[string]any, size),
		mu:            sync.RWMutex{},
		bootstrappers: make([]func(), 0, size),
	}
}

// serviceProvider is a function that provides an instance of type T
type serviceProvider[T any] func(*Container) (T, error)

// containerService represents a service in the container
type containerService[T any] struct {
	mu       sync.Mutex
	instance T
	made     bool
	provider serviceProvider[T]
}

// serviceName returns a unique identifier for the type T
func serviceName[T any]() string {
	typeForT := reflect.TypeOf((*T)(nil)).Elem()

	if typeForT.Kind() == reflect.Ptr {
		typeForT = typeForT.Elem()
	}

	if typeForT.Name() != "" {
		return typeForT.PkgPath() + "." + typeForT.Name()
	}

	panic("unnamed type")
}

// Register registers a service that will be resolved during bootstrap
func Register[T any](c *Container, provider serviceProvider[T]) {
	fmt.Println(fmt.Sprintf("Registering service %s", serviceName[T]()))
	c.mu.Lock()
	defer c.mu.Unlock()

	name := serviceName[T]()

	c.services[name] = &containerService[T]{
		provider: provider,
	}

	c.bootstrappers = append(c.bootstrappers, func() {
		_ = Make[T](c)
	})
}

// RegisterDeferred registers a service that will be resolved only when requested
func RegisterDeferred[T any](c *Container, provider serviceProvider[T]) {
	c.mu.Lock()
	defer c.mu.Unlock()

	name := serviceName[T]()

	c.services[name] = &containerService[T]{
		provider: provider,
	}
}

// make resolves the service instance
func (s *containerService[T]) make(c *Container) (T, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.made {
		return s.instance, nil
	}

	instance, err := s.provider(c)
	if err != nil {
		var zero T
		return zero, err
	}

	s.instance = instance
	s.made = true

	return instance, nil
}

// Make resolves and returns an instance of type T
func Make[T any](c *Container) T {
	c.mu.RLock()
	defer c.mu.RUnlock()

	name := serviceName[T]()

	service, exists := c.services[name]
	if !exists {
		panic(fmt.Sprintf("service %s is not registered", name))
	}

	typedService, ok := service.(*containerService[T])
	if !ok {
		panic(fmt.Sprintf("service %s is not of the expected type", name))
	}

	instance, err := typedService.make(c)
	if err != nil {
		panic(fmt.Sprintf("failed to make service %s: %v", name, err))
	}

	return instance
}

// Bootstrap resolves all registered services
func (c *Container) Bootstrap() {
	for _, bootstrapper := range c.bootstrappers {
		bootstrapper()
	}

	// Free up memory
	c.bootstrappers = nil
}
