package juno_go

import (
	"errors"
	"juno-go/connection"
	"juno-go/models"
	"juno-go/protocol"
	"juno-go/utils/request_types"
	"net"
	"sync"
)

type RequestListType struct {
	sync.RWMutex
	m map[string]chan interface{}
}
type FunctionListType struct {
	sync.RWMutex
	m map[string]func(map[string]interface{}) interface{}
}
type HookListType struct {
	sync.RWMutex
	m map[string][]func(interface{})
}
type MutexBool struct {
	sync.RWMutex
	value bool
}

type JunoModule struct {
	connection    connection.BaseConnection
	protocol      protocol.BaseProtocol
	requests      RequestListType
	functions     FunctionListType
	hookListeners HookListType
	messageBuffer []byte
	registered    MutexBool
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
		connection: connection,
		protocol:   protocol,
		requests: RequestListType{
			m: make(map[string]chan interface{}),
		},
		functions: FunctionListType{
			m: make(map[string]func(map[string]interface{}) interface{}),
		},
		hookListeners: HookListType{
			m: make(map[string][]func(interface{})),
		},
		messageBuffer: []byte{},
		registered: MutexBool{
			value: false,
		},
	}
}

func (module *JunoModule) Initialize(moduleId, version string, dependencies map[string]string) (chan interface{}, error) {
	err := module.connection.SetupConnection()
	if err != nil {
		return nil, err
	}
	module.connection.SetOnDataHandler(module.onDataHandler)

	request := protocol.Initialize(module.protocol, moduleId, version, dependencies)
	return module.sendRequest(request)
}

func (module *JunoModule) DeclareFunction(fnName string, fn func(map[string]interface{}) interface{}) (chan interface{}, error) {
	module.functions.Lock()
	module.functions.m[fnName] = fn
	module.functions.Unlock()
	return module.sendRequest(
		protocol.DeclareFunction(module.protocol, fnName),
	)
}

func (module *JunoModule) CallFunction(fnName string, args map[string]interface{}) (chan interface{}, error) {
	return module.sendRequest(
		protocol.CallFunction(module.protocol, fnName, args),
	)
}

func (module *JunoModule) RegisterHook(hook string, cb func(interface{})) (chan interface{}, error) {
	module.hookListeners.Lock()
	defer module.hookListeners.Unlock()
	if module.hookListeners.m[hook] != nil {
		module.hookListeners.m[hook] = append(module.hookListeners.m[hook], cb)
	} else {
		module.hookListeners.m[hook] = []func(interface{}){cb}
	}
	return module.sendRequest(
		protocol.RegisterHook(module.protocol, hook),
	)
}

func (module *JunoModule) TriggerHook(hook string) (chan interface{}, error) {
	return module.sendRequest(
		protocol.TriggerHook(module.protocol, hook),
	)
}

func (module *JunoModule) Close() error {
	return module.connection.CloseConnection()
}

func (module *JunoModule) sendRequest(message models.BaseMessage) (chan interface{}, error) {
	module.registered.RLock()
	if message.GetType() == request_types.RegisterModuleRequest && module.registered.value {
		return nil, errors.New("module already registered")
	}

	encoded, err := module.protocol.Encode(message)
	if err != nil {
		return nil, err
	}
	if module.registered.value || message.GetType() == request_types.RegisterModuleRequest {
		err = module.connection.Send(encoded)
		if err != nil {
			return nil, err
		}
	} else {
		module.messageBuffer = append(module.messageBuffer, encoded...)
	}
	module.registered.RUnlock()

	module.requests.Lock()
	channel := make(chan interface{})
	module.requests.m[message.GetRequestId()] = channel
	module.requests.Unlock()
	return channel, nil
}

func (module *JunoModule) onDataHandler(data []byte) {
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
			value = module.executeHookTriggered(response.(models.TriggerHookResponse))
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

	module.requests.Lock()
	if module.requests.m[response.GetRequestId()] != nil {
		module.requests.m[response.GetRequestId()] <- value
		delete(module.requests.m, response.GetRequestId())
	}
	module.requests.Unlock()
}

func (module *JunoModule) executeFunctionCall(request models.FunctionCallRequest) bool {
	module.functions.RLock()
	defer module.functions.RUnlock()
	if module.functions.m[request.Function] != nil {
		res := module.functions.m[request.Function](request.Arguments)
		if res.(chan interface{}) != nil {
			res = <-res.(chan interface{})
		}
		channel, err := module.sendRequest(models.FunctionCallResponse{
			RequestId: request.RequestId,
			Data:      res,
		})
		if err != nil {
			return false
		}
		_ = <-channel
		return true
	} else {
		// Function wasn't found in the module.
		return false
	}
}

func (module *JunoModule) executeHookTriggered(request models.TriggerHookResponse) bool {
	if &request.Hook != nil {
		module.hookListeners.RLock()
		defer module.hookListeners.RUnlock()
		// Hook triggered by another module.
		if request.Hook == `juno.activated` {
			module.registered.Lock()
			module.registered.value = true
			if module.messageBuffer != nil {
				_ = module.connection.Send(module.messageBuffer)
			}
			module.registered.Unlock()
		} else if request.Hook == `juno.deactivated` {
			module.registered.Lock()
			module.registered.value = false
			module.registered.Unlock()
		} else if module.hookListeners.m[request.Hook] != nil {
			for _, listener := range module.hookListeners.m[request.Hook] {
				listener(nil)
			}
		}
		return true
	} else {
		// This module triggered the hook.
		return true
	}
}
