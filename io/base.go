package io

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver/appconfig"
)

type Coder interface {
	Marshal(v interface{}) ([]byte, error)
	Unmarshal(data []byte, v interface{}) error
}

type JsonCoder struct{}

func (JsonCoder) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (JsonCoder) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

type Backend interface {
	Set(key string, data []byte) error
	Get(key string) ([]byte, error)
	Del(key string) error
}

type MemBackend struct {
	Data map[string][]byte
}

func (m *MemBackend) Set(key string, data []byte) error {
	m.Data[key] = data
	return nil
}

func (m *MemBackend) Get(key string) ([]byte, error) {
	return m.Data[key], nil
}

func (m *MemBackend) Del(key string) error { delete(m.Data, key); return nil }

var defaultCoder Coder = JsonCoder{}
var defaultBackend Backend = &MemBackend{Data: make(map[string][]byte)}

func InitBackendCoder(backend Backend, coder Coder) {
	//defaultCoder = coder
	if coder == nil {
		defaultCoder = coder
	}
	//defaultBackend = backend
	if backend == nil {
		defaultBackend = backend
	}
}

func LoadData(logger fklog.FKLogI, key string, value interface{}) error {
	// check data is a pointer
	if reflect.TypeOf(value).Kind() != reflect.Ptr {
		return errors.New("data is not a pointer")
	}

	data, err := defaultBackend.Get(key)
	if err != nil {
		return err
	}

	if data == nil {
		return nil
	}

	err = defaultCoder.Unmarshal(data, value)
	if err != nil {
		return err
	}

	return nil
}

func SaveData(logger fklog.FKLogI, key string, value interface{}) error {
	data, err := defaultCoder.Marshal(value)
	if err != nil {
		return err
	}

	err = defaultBackend.Set(key, data)
	if err != nil {
		return err
	}

	return nil
}

func DeleteData(logger fklog.FKLogI, key string) error {
	return defaultBackend.Del(key)
}

func paddingKey(key string) string {
	appConfig := appconfig.GlobalConfig()
	return fmt.Sprintf("svr%s:%s", appConfig.Global.SectionID, key)
}

func LoadSvrData(logger fklog.FKLogI, key string, data interface{}) error {
	return LoadData(logger, paddingKey(key), data)
}

func SaveSvrData(logger fklog.FKLogI, key string, value interface{}) error {
	return SaveData(logger, paddingKey(key), value)
}

func DeleteSvrData(logger fklog.FKLogI, key string) error {
	return DeleteData(logger, paddingKey(key))
}
