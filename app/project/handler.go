package project

import "service-base-go/pkg/handler"

type IProjectHandler interface {
	handler.HandlerInterface[handler.Request, handler.Response]
}
