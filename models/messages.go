package models

import (
	"juno-go/utils/request_types"
)

type BaseMessage interface {
	GetType() uint64
	GetRequestId() string
}

type RegisterModuleRequest struct {
	RequestId string `json:"requestId"`
	ModuleId string `json:"moduleId"`
	Version string `json:"version"`
	Dependencies map[string]string `json:"dependencies"`
}

func (message RegisterModuleRequest) GetType() uint64 {
	return request_types.RegisterModuleRequest
}

func (message RegisterModuleRequest) GetRequestId() string {
	return message.RequestId
}

type RegisterModuleResponse struct {
	RequestId string `json:"requestId"`
}

func (message RegisterModuleResponse) GetType() uint64 {
	return request_types.RegisterModuleResponse
}

func (message RegisterModuleResponse) GetRequestId() string {
	return message.RequestId
}

type FunctionCallRequest struct {
	RequestId string `json:"requestId"`
	Function string `json:"function"`
	Arguments map[string]interface{} `json:"arguments"`
}

func (message FunctionCallRequest) GetType() uint64 {
	return request_types.FunctionCallRequest
}

func (message FunctionCallRequest) GetRequestId() string {
	return message.RequestId
}

type FunctionCallResponse struct {
	RequestId string `json:"requestId"`
	Data interface{} `json:"data"`
}

func (message FunctionCallResponse) GetType() uint64 {
	return request_types.FunctionCallResponse
}

func (message FunctionCallResponse) GetRequestId() string {
	return message.RequestId
}

type RegisterHookRequest struct {
	RequestId string `json:"requestId"`
	Hook string `json:"hook"`
}

func (message RegisterHookRequest) GetType() uint64 {
	return request_types.RegisterHookRequest
}

func (message RegisterHookRequest) GetRequestId() string {
	return message.RequestId
}

type RegisterHookResponse struct {
	RequestId string `json:"requestId"`
}

func (message RegisterHookResponse) GetType() uint64 {
	return request_types.RegisterHookResponse
}

func (message RegisterHookResponse) GetRequestId() string {
	return message.RequestId
}

type TriggerHookRequest struct {
	RequestId string `json:"requestId"`
	Hook string `json:"hook"`
}

func (message TriggerHookRequest) GetType() uint64 {
	return request_types.TriggerHookRequest
}

func (message TriggerHookRequest) GetRequestId() string {
	return message.RequestId
}

type TriggerHookResponse struct {
	RequestId string `json:"requestId"`
	Hook string `json:"hook"`
}

func (message TriggerHookResponse) GetType() uint64 {
	return request_types.TriggerHookResponse
}

func (message TriggerHookResponse) GetRequestId() string {
	return message.RequestId
}

type DeclareFunctionRequest struct {
	RequestId string `json:"requestId"`
	Function string `json:"function"`
}

func (message DeclareFunctionRequest) GetType() uint64 {
	return request_types.DeclareFunctionRequest
}

func (message DeclareFunctionRequest) GetRequestId() string {
	return message.RequestId
}

type DeclareFunctionResponse struct {
	RequestId string `json:"requestId"`
	Function string `json:"function"`
}

func (message DeclareFunctionResponse) GetType() uint64 {
	return request_types.DeclareFunctionResponse
}

func (message DeclareFunctionResponse) GetRequestId() string {
	return message.RequestId
}

type ErrorMessage struct {
	RequestId string `json:"requestId"`
	Error uint32 `json:"error"`
}

func (message ErrorMessage) GetType() uint64 {
	return request_types.Error
}

func (message ErrorMessage) GetRequestId() string {
	return message.RequestId
}

type UnknownMessage struct {
	RequestId string `json:"requestId"`
}

func (message UnknownMessage) GetType() uint64 {
	return request_types.Error
}

func (message UnknownMessage) GetRequestId() string {
	return message.RequestId
}
