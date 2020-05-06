package juno_go

import (
	"errors"
	"juno-go/connection"
	"juno-go/models"
	"juno-go/protocol"
	"juno-go/utils/request_types"
	"net"
)

type RequestListType map[string]chan interface{}
type FunctionListType map[string]func(map[string]interface{}) interface{}
type HookListType map[string][]func(interface{})

type JunoModule struct {
	connection    connection.BaseConnection
	protocol      protocol.BaseProtocol
	requests      RequestListType
	functions     FunctionListType
	hookListeners HookListType
	messageBuffer []byte
	registered    bool
}

func Default(connectionPath string) JunoModule {
	addr, err := net.ResolveTCPAddr("tcp", connectionPath)
	if err != nil {
		return FromUnixSocket(connectionPath)
	}
	return FromInetSocket(addr.IP, uint16(addr.Port))
}

func FromInetSocket(ip net.IP, port uint16) JunoModule {
	return NewJunoModule(protocol.NewJsonProtocol(), connection.NewInetSocketConnection(ip.String(), port))
}

func FromUnixSocket(socketPath string) JunoModule {
	return NewJunoModule(protocol.NewJsonProtocol(), connection.NewUnixSocketConnection(socketPath))
}

func NewJunoModule(protocol protocol.BaseProtocol, connection connection.BaseConnection) JunoModule {
	return JunoModule{
		connection:    connection,
		protocol:      protocol,
		requests:      RequestListType{},
		functions:     FunctionListType{},
		hookListeners: HookListType{},
		messageBuffer: []byte{},
		registered:    false,
	}
}

func (module JunoModule) Initialize(moduleId, version string, dependencies map[string]string) error {
	err := module.connection.SetupConnection()
	if err != nil {
		return err
	}
	module.connection.SetOnDataHandler(module.onDataHandler)

	request := protocol.Initialize(module.protocol, moduleId, version, dependencies)
	err = module.sendRequest(request)
	if err != nil {
		return err
	}

	module.registered = true
	return nil
}

func (module JunoModule) sendRequest(message models.BaseMessage) error {
	if message.GetType() == request_types.RegisterModuleRequest && module.registered {
		return errors.New("module already registered")
	}

	encoded, err := module.protocol.Encode(message)
	if err != nil {
		return err
	}
	if module.registered || message.GetType() == request_types.RegisterModuleRequest {
		return module.connection.Send(encoded)
	} else {
		module.messageBuffer = append(module.messageBuffer, encoded...)
	}
	return nil
}

func (module JunoModule) onDataHandler(data []byte) {
	response := module.protocol.Decode(data)
	var value interface{}
	switch response.GetType() {
	case request_types.RegisterModuleRequest:
		{
			value = true
			break
		}
	case request_types.FunctionCallResponse:
		{
			value = response.(models.FunctionCallResponse).Data
			break
		}
	case request_types.DeclareFunctionResponse:
		{
			value = true
			break
		}
	case request_types.RegisterHookResponse:
		{
			value = true
			break
		}
	case request_types.TriggerHookResponse:
		{
			value = module.executeHookTriggered(response.(models.TriggerHookRequest))
			break
		}
	case request_types.FunctionCallRequest:
		{
			value = module.executeFunctionCall(response.(models.FunctionCallRequest))
			break
		}
	default:
		{
			value = false
			break
		}
	}
	module.requests[response.GetRequestId()] <- value
	delete(module.requests, response.GetRequestId())
}

func (module JunoModule) executeFunctionCall(request models.FunctionCallRequest) bool {
	if module.functions[request.Function] != nil {
		res := module.functions[request.Function](request.Arguments)
		if res.(chan interface{}) != nil {
			res = <-res.(chan interface{})
		}
		err := module.sendRequest(models.FunctionCallResponse{
			RequestId: request.RequestId,
			Data:      res,
		})
		return err == nil
	} else {
		// Function wasn't found in the module.
		return false
	}
}

func (module JunoModule) executeHookTriggered(request models.TriggerHookRequest) bool {
	if &request.Hook != nil {
		// Hook triggered by another module.
		if request.Hook == `juno.activated` {
			module.registered = true
			if module.messageBuffer != nil {
				_ = module.connection.Send(module.messageBuffer)
			}
		} else if module.hookListeners[request.Hook] != nil {
			for _, listener := range module.hookListeners[request.Hook] {
				listener(nil)
			}
		}
		return true
	} else {
		// This module triggered the hook.
		return true
	}
}

