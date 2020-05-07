package protocol

import (
	"fmt"
	"juno-go/models"
	"time"
)

type BaseProtocol interface {
	Encode(models.BaseMessage) ([]byte, error)
	Decode([]byte) models.BaseMessage
	SetModuleId(string)
	GetModuleId() string
}

func GenerateRequestId(moduleId string) string {
	return fmt.Sprintf("%s-%d", moduleId, time.Now().UnixNano())
}

func Initialize(protocol BaseProtocol, moduleId, version string, dependencies map[string]string) models.BaseMessage {
	protocol.SetModuleId(moduleId)
	return models.RegisterModuleRequest{
		RequestId:    GenerateRequestId(moduleId),
		ModuleId:     moduleId,
		Version:      version,
		Dependencies: dependencies,
	}
}

func RegisterHook(protocol BaseProtocol, hook string) models.BaseMessage {
	return models.RegisterHookRequest{
		RequestId: GenerateRequestId(protocol.GetModuleId()),
		Hook:      hook,
	}
}

func TriggerHook(protocol BaseProtocol, hook string, data interface{}) models.BaseMessage {
	return models.TriggerHookRequest{
		RequestId: GenerateRequestId(protocol.GetModuleId()),
		Hook:      hook,
		Data:      data,
	}
}

func DeclareFunction(protocol BaseProtocol, function string) models.BaseMessage {
	return models.DeclareFunctionRequest{
		RequestId: GenerateRequestId(protocol.GetModuleId()),
		Function:  function,
	}
}

func CallFunction(protocol BaseProtocol, function string, arguments map[string]interface{}) models.BaseMessage {
	return models.FunctionCallRequest{
		RequestId: GenerateRequestId(protocol.GetModuleId()),
		Function:  function,
		Arguments: arguments,
	}
}
