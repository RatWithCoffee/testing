package lab6

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

type HttpHandlerUtils struct {
	Log zap.Logger
}

type Response struct {
	Status string `json:"status,omitempty"`
}

type HttpErr struct {
	Error string `json:"err,omitempty"`
}

func (utils *HttpHandlerUtils) Ok(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	respJson, err := json.Marshal(Response{Status: "ok"})

	if err != nil {
		utils.LogMarshalErr(err)
		return
	}
	_, err = w.Write(respJson)
	if err != nil {
		utils.LogWriteRespErr(err)
	}
}

func (utils *HttpHandlerUtils) BadRequest(w http.ResponseWriter, errToSend string) {
	utils.writeErrorToResp(w, errToSend, http.StatusBadRequest)
}

func (utils *HttpHandlerUtils) Unauthorized(w http.ResponseWriter) {
	errToSend := "authorization error"
	utils.writeErrorToResp(w, errToSend, http.StatusUnauthorized)
}

func (utils *HttpHandlerUtils) InternalServerError(w http.ResponseWriter, errToSend string) {
	utils.writeErrorToResp(w, errToSend, http.StatusInternalServerError)
}

func (utils *HttpHandlerUtils) Forbidden(w http.ResponseWriter, errToSend string) {
	utils.writeErrorToResp(w, errToSend, http.StatusForbidden)
}

func (utils *HttpHandlerUtils) NotFound(w http.ResponseWriter, errToSend string) {
	utils.writeErrorToResp(w, errToSend, http.StatusNotFound)
}

func (utils *HttpHandlerUtils) PayloadTooLarge(w http.ResponseWriter, errToSend string) {
	utils.writeErrorToResp(w, errToSend, http.StatusRequestEntityTooLarge)
}

func (utils *HttpHandlerUtils) StatusUnsupportedMediaType(w http.ResponseWriter, errToSend string) {
	utils.writeErrorToResp(w, errToSend, http.StatusUnsupportedMediaType)
}

func (utils *HttpHandlerUtils) writeErrorToResp(w http.ResponseWriter, errToSend string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	respJson, err := json.Marshal(HttpErr{Error: errToSend})
	if err != nil {
		utils.LogMarshalErr(err)
		return
	}
	_, err = w.Write(respJson)
	if err != nil {
		utils.LogWriteRespErr(err)
	}
}

func (utils *HttpHandlerUtils) WriteJsonToResp(w http.ResponseWriter, resp interface{}) {
	w.Header().Set("Content-Type", "application/json")

	respJson, err := json.Marshal(resp)
	if err != nil {
		utils.LogMarshalErr(err)
		return
	}
	_, err = w.Write(respJson)
	if err != nil {
		utils.LogWriteRespErr(err)
	}
}

func (utils *HttpHandlerUtils) LogMarshalErr(err error) {
	utils.Log.Error("marshal error", zap.Error(err))
}

func (utils *HttpHandlerUtils) LogWriteRespErr(err error) {
	utils.Log.Error("error writing response", zap.Error(err))
}
