package protocol

import (
	"encoding/json"

	"github.com/bytesonus/juno-go/utils/request_types"

	"github.com/bytesonus/juno-go/utils/request_keys"

	"github.com/bytesonus/juno-go/models"
)

type JsonProtocol struct {
	moduleId string
}

func (protocol *JsonProtocol) Encode(message models.BaseMessage) ([]byte, error) {
	var genericMap map[string]interface{}

	switch request := message.(type) {
	case models.RegisterModuleRequest:
		{
			genericMap = map[string]interface{}{
				request_keys.RequestId:    request.RequestId,
				request_keys.Type:         request_types.RegisterModuleRequest,
				request_keys.ModuleId:     request.ModuleId,
				request_keys.Version:      request.Version,
				request_keys.Dependencies: request.Dependencies,
			}
			break
		}
	case models.RegisterModuleResponse:
		{
			genericMap = map[string]interface{}{
				request_keys.RequestId: request.RequestId,
				request_keys.Type:      request_types.RegisterModuleResponse,
			}
			break
		}
	case models.FunctionCallRequest:
		{
			genericMap = map[string]interface{}{
				request_keys.RequestId: request.RequestId,
				request_keys.Type:      request_types.FunctionCallRequest,
				request_keys.Function:  request.Function,
				request_keys.Arguments: request.Arguments,
			}
			break
		}
	case models.FunctionCallResponse:
		{
			genericMap = map[string]interface{}{
				request_keys.RequestId: request.RequestId,
				request_keys.Type:      request_types.FunctionCallResponse,
				request_keys.Data:      request.Data,
			}
			break
		}
	case models.RegisterHookRequest:
		{
			genericMap = map[string]interface{}{
				request_keys.RequestId: request.RequestId,
				request_keys.Type:      request_types.RegisterHookRequest,
				request_keys.Hook:      request.Hook,
			}
			break
		}
	case models.RegisterHookResponse:
		{
			genericMap = map[string]interface{}{
				request_keys.RequestId: request.RequestId,
				request_keys.Type:      request_types.RegisterHookResponse,
			}
			break
		}
	case models.TriggerHookRequest:
		{
			genericMap = map[string]interface{}{
				request_keys.RequestId: request.RequestId,
				request_keys.Type:      request_types.TriggerHookRequest,
				request_keys.Hook:      request.Hook,
			}
			break
		}
	case models.TriggerHookResponse:
		{
			genericMap = map[string]interface{}{
				request_keys.RequestId: request.RequestId,
				request_keys.Type:      request_types.TriggerHookResponse,
			}
			break
		}
	case models.DeclareFunctionRequest:
		{
			genericMap = map[string]interface{}{
				request_keys.RequestId: request.RequestId,
				request_keys.Type:      request_types.DeclareFunctionRequest,
				request_keys.Function:  request.Function,
			}
			break
		}
	case models.DeclareFunctionResponse:
		{
			genericMap = map[string]interface{}{
				request_keys.RequestId: request.RequestId,
				request_keys.Type:      request_types.DeclareFunctionResponse,
				request_keys.Function:  request.Function,
			}
			break
		}
	case models.ErrorMessage:
		{
			genericMap = map[string]interface{}{
				request_keys.RequestId: request.RequestId,
				request_keys.Type:      request_types.Error,
				request_keys.Error:     request.Error,
			}
			break
		}
	case models.UnknownMessage:
	default:
		{
			genericMap = map[string]interface{}{
				request_keys.RequestId: "undefined",
				request_keys.Type:      request_types.Error,
				request_keys.Error:     0,
			}
			break
		}
	}
	data, err := json.Marshal(genericMap)
	if err != nil {
		return nil, err
	}
	return append(data, byte('\n')), nil
}

func (protocol *JsonProtocol) Decode(data []byte) models.BaseMessage {
	var message map[string]interface{}
	err := json.Unmarshal(data, &message)
	if err != nil {
		return models.UnknownMessage{RequestId: "undefined"}
	}

	switch message["type"].(float64) {
	case request_types.RegisterModuleRequest:
		{
			var request models.RegisterModuleRequest
			err := json.Unmarshal(data, &request)
			if err != nil {
				return models.UnknownMessage{RequestId: "undefined"}
			}
			return request
		}
	case request_types.RegisterModuleResponse:
		{
			var request models.RegisterModuleResponse
			err := json.Unmarshal(data, &request)
			if err != nil {
				return models.UnknownMessage{RequestId: "undefined"}
			}
			return request
		}
	case request_types.FunctionCallRequest:
		{
			var request models.FunctionCallRequest
			err := json.Unmarshal(data, &request)
			if err != nil {
				return models.UnknownMessage{RequestId: "undefined"}
			}
			return request
		}
	case request_types.FunctionCallResponse:
		{
			var request models.FunctionCallResponse
			err := json.Unmarshal(data, &request)
			if err != nil {
				return models.UnknownMessage{RequestId: "undefined"}
			}
			return request
		}
	case request_types.RegisterHookRequest:
		{
			var request models.RegisterHookRequest
			err := json.Unmarshal(data, &request)
			if err != nil {
				return models.UnknownMessage{RequestId: "undefined"}
			}
			return request
		}
	case request_types.RegisterHookResponse:
		{
			var request models.RegisterHookResponse
			err := json.Unmarshal(data, &request)
			if err != nil {
				return models.UnknownMessage{RequestId: "undefined"}
			}
			return request
		}
	case request_types.TriggerHookRequest:
		{
			var request models.TriggerHookRequest
			err := json.Unmarshal(data, &request)
			if err != nil {
				return models.UnknownMessage{RequestId: "undefined"}
			}
			return request
		}
	case request_types.TriggerHookResponse:
		{
			var request models.TriggerHookResponse
			err := json.Unmarshal(data, &request)
			if err != nil {
				return models.UnknownMessage{RequestId: "undefined"}
			}
			return request
		}
	case request_types.DeclareFunctionRequest:
		{
			var request models.DeclareFunctionRequest
			err := json.Unmarshal(data, &request)
			if err != nil {
				return models.UnknownMessage{RequestId: "undefined"}
			}
			return request
		}
	case request_types.DeclareFunctionResponse:
		{
			var request models.DeclareFunctionResponse
			err := json.Unmarshal(data, &request)
			if err != nil {
				return models.UnknownMessage{RequestId: "undefined"}
			}
			return request
		}
	case request_types.Error:
		{
			var request models.ErrorMessage
			err := json.Unmarshal(data, &request)
			if err != nil {
				return models.UnknownMessage{RequestId: "undefined"}
			}
			return request
		}
	default:
		{
			var request models.UnknownMessage
			err := json.Unmarshal(data, &request)
			if err != nil {
				return models.UnknownMessage{RequestId: "undefined"}
			}
			return request
		}
	}
}

func (protocol *JsonProtocol) SetModuleId(moduleId string) {
	protocol.moduleId = moduleId
}

func (protocol *JsonProtocol) GetModuleId() string {
	return protocol.moduleId
}

func NewJsonProtocol() *JsonProtocol {
	return &JsonProtocol{moduleId: ""}
}
