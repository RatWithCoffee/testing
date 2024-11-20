package lab6

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"go.uber.org/zap"
)

type TestCase struct {
	Name           string
	ContentType    string
	Body           getBody
	StatusCode     int
	ResponseBody   string
	ReadTo         interface{}
	IsReturningErr bool
}

const contTypeJson = "application/json"

type Item struct {
	ID          int    `json:"id" validate:"required"`
	Description string `json:"description" validate:"required"`
}

type RequestBody struct {
	Items []Item `json:"items"`
}

var cases = []TestCase{
	{
		Name:        "PositiveCase",
		ContentType: contTypeJson,
		Body: func() []byte {
			return []byte("{\"id\": 1 , \"description\": \"dsfdgf\"}")
		},
		StatusCode:     http.StatusOK,
		ReadTo:         &Item{},
		IsReturningErr: false,
	},
	{
		Name: "LargeBody",
		Body: func() []byte {
			return jsonFromStruct(RequestBody{Items: generateItems()})
		},
		ContentType:    contTypeJson,
		StatusCode:     http.StatusRequestEntityTooLarge,
		ResponseBody:   bodyTooLarge,
		IsReturningErr: true,
	},
	{
		Name:           "WrongContentType",
		ContentType:    "text/plain",
		StatusCode:     http.StatusUnsupportedMediaType,
		ResponseBody:   unsupportedMediaTypeMsg,
		IsReturningErr: true,
	},
	{
		Name:           "NoContentTypeHeader",
		ContentType:    "",
		StatusCode:     http.StatusUnsupportedMediaType,
		ResponseBody:   unsupportedMediaTypeMsg,
		IsReturningErr: true,
	},
	{
		Name:        "NoRequiredField",
		ContentType: contTypeJson,
		Body: func() []byte {
			return jsonFromStruct(Item{ID: 1})
		},
		StatusCode:     http.StatusBadRequest,
		ReadTo:         &Item{},
		ResponseBody:   "Key: 'Item.Description' Error:Field validation for 'Description' failed on the 'required' tag",
		IsReturningErr: true,
	},
	{
		Name:        "InvalidType",
		ContentType: contTypeJson,
		Body: func() []byte {
			return []byte("{\"id\": \"asdf\" , \"description\": \"dsfdgf\"}")
		},
		StatusCode:     http.StatusBadRequest,
		ReadTo:         &Item{},
		ResponseBody:   "request body contains an invalid value for the \"id\" field (at position 13)",
		IsReturningErr: true,
	},
	{
		Name:        "BadJson(ErrUnexpectedEOF)",
		ContentType: contTypeJson,
		Body: func() []byte {
			return []byte("{\"id\": \"asdf\" , \"description\": \"dsfdgf\"")
		},
		StatusCode:     http.StatusBadRequest,
		ReadTo:         &Item{},
		ResponseBody:   badJson,
		IsReturningErr: true,
	},
	{
		Name:        "BadJson",
		ContentType: contTypeJson,
		Body: func() []byte {
			return []byte("{\"id\" \"asdf\" , \"description\": \"dsfdgf\"}")
		},
		StatusCode:     http.StatusBadRequest,
		ReadTo:         &Item{},
		ResponseBody:   "request body contains badly-formed JSON (at position 7)",
		IsReturningErr: true,
	},
	{
		Name:        "BadJson",
		ContentType: contTypeJson,
		Body: func() []byte {
			return nil
		},
		StatusCode:     http.StatusBadRequest,
		ReadTo:         &Item{},
		ResponseBody:   emptyBody,
		IsReturningErr: true,
	},
	{
		Name:        "BadJson",
		ContentType: contTypeJson,
		Body: func() []byte {
			return []byte("{\"id\": 1 , \"description\": \"dsfdgf\"}")
		},
		StatusCode:     http.StatusInternalServerError,
		ReadTo:         Item{},
		ResponseBody:   "",
		IsReturningErr: true,
	},
	{
		Name:        "SeveralJsonObjects",
		ContentType: contTypeJson,
		Body: func() []byte {
			return []byte("{\"id\": 1 , \"description\": \"dsfdgf\"}{\"id\": 1 , \"description\": \"dsfdgf\"}")
		},
		StatusCode:     http.StatusBadRequest,
		ReadTo:         &Item{},
		ResponseBody:   severalJsonObjects,
		IsReturningErr: true,
	},
}

var logTest = GetLoggerForTest()

func GetLoggerForTest() zap.Logger {
	return *zap.Must(zap.NewProduction())
}

var respUtils = HttpHandlerUtils{Log: logTest}

type getBody func() []byte

func jsonFromStruct(body interface{}) []byte {
	jsonData, err := json.Marshal(body)
	if err != nil {
		logTest.Fatal("error marshalling json:", zap.Error(err))
	}
	return jsonData
}

func TestReadJsonBody(t *testing.T) {
	for i, testCase := range cases {
		t.Logf("Run test case %d, name: %s", i, testCase.Name)

		var jsonData []byte
		if testCase.Body != nil {
			jsonData = testCase.Body()
		}

		r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(jsonData))
		if testCase.ContentType != "" {
			r.Header.Set("Content-Type", testCase.ContentType)
		}
		w := httptest.NewRecorder()

		err := respUtils.ReadJsonBody(w, r, testCase.ReadTo)
		if testCase.IsReturningErr && err == nil {
			t.Fatalf("expected not nil error")
		}
		if !testCase.IsReturningErr && err != nil {
			t.Fatalf("expected err: %v, got nil", err)
		}

		res := w.Result()

		if res.StatusCode != testCase.StatusCode {
			t.Fatalf("got status code: %d, expected: %d", res.StatusCode, testCase.StatusCode)
		}

		if testCase.ResponseBody == "" {
			continue
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			t.Fatalf("error reading response body: %v", err)
		}
		defer res.Body.Close()
		var respBody HttpErr
		err = json.Unmarshal(body, &respBody)
		if err != nil {
			t.Fatalf("unmarshal response body error, err: %v, body: %s", err, string(body))
		}
		expErr := HttpErr{Error: testCase.ResponseBody}
		if !reflect.DeepEqual(expErr, respBody) {
			t.Fatalf("wrong response body, expected: %v, got: %v", expErr, respBody)
		}
	}
}

func generateItems() []Item {
	items := make([]Item, 0)
	currSize := 0
	stringSize := 100
	i := 0
	for currSize <= maxBodySize {
		items = append(items, Item{
			ID:          i,
			Description: generateRandomString(stringSize),
		})
		currSize += stringSize
		i++
	}
	return items
}

func generateRandomString(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = 'a'
	}
	return string(b)
}
