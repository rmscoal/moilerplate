package middleware

import (
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/delivery/v1/model"
)

// Middleware struct server as the parent class of all
// middleware. This enables any middlewares to have the
// inherited methods of a GiinBaseController.
type Middleware struct {
	model.BaseControllerV1
}

var (
	once                sync.Once
	middlewareSingleton *Middleware
)

// NewMiddleware creates a new middleware if and only
// if there were no existing middleware instance. It
// follows the singleton creational pattern for resource
// effectiveness.
func NewMiddleware() *Middleware {
	if middlewareSingleton == nil {
		once.Do(func() {
			middlewareSingleton = &Middleware{}
		})
	}
	return middlewareSingleton
}

func (m *Middleware) addToContext(c *gin.Context, key string, value any) {
	// Passing value by context is the best practice, however let's try gin's feature
	// Pass value to context code:
	//
	// c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), key, value))
	//
	// and receive it like:
	//
	// pquery := c.Request.Context().Value(keyAsString).(dto.PaginationDTORequest)

	if c.Keys == nil {
		c.Keys = make(map[string]any)
	}
	c.Keys[key] = value
}
