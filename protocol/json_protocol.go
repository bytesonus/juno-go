package protocol

import (
	"encoding/json"
	"juno-go/models"
)

type JsonProtocol struct {
	moduleId string
}

func (protocol JsonProtocol) Encode(message models.BaseMessage) ([]byte, error) {
	return json.Marshal(message)
}

func (protocol JsonProtocol) Decode(data []byte) models.BaseMessage {
	var message models.BaseMessage
	err := json.Unmarshal(data, &message)
	if err != nil {
		return models.UnknownMessage{RequestId: "undefined"}
	}
	return message
}

func (protocol JsonProtocol) SetModuleId(moduleId string) {
	protocol.moduleId = moduleId
}

func (protocol JsonProtocol) GetModuleId() string {
	return protocol.moduleId
}

func NewJsonProtocol() JsonProtocol {
	return JsonProtocol{moduleId: ""}
}

