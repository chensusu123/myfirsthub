package io

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

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
	Set(ctx context.Context, key string, data []byte) error
	Get(ctx context.Context, key string) ([]byte, error)
	Del(ctx context.Context, key string) error
}

type MemBackend struct {
	Data map[string][]byte
}

func (m *MemBackend) Set(ctx context.Context, key string, data []byte) error {
	m.Data[key] = data
	return nil
}

func (m *MemBackend) Get(ctx context.Context, key string) ([]byte, error) {
	return m.Data[key], nil
}

func (m *MemBackend) Del(ctx context.Context, key string) error { delete(m.Data, key); return nil }

var defaultCoder Coder = JsonCoder{}
var defaultBackend Backend = &MemBackend{Data: make(map[string][]byte)}

func InitBackendCoder(backend Backend, coder Coder) {
	//defaultCoder = coder
	if coder != nil {
		defaultCoder = coder
	}
	//defaultBackend = backend
	if backend != nil {
		defaultBackend = backend
	}
}

func LoadData(ctx context.Context, key string, value interface{}) error {
	key = paddingKey(key, "0")
	return loadData(ctx, key, value)
}

func SaveData(ctx context.Context, key string, value interface{}) error {
	key = paddingKey(key, "0")
	return saveData(ctx, key, value)
}

func DeleteData(ctx context.Context, key string) error {
	key = paddingKey(key, "0")
	return deleteData(ctx, key)
}

func paddingKey(key string, svr string) string {
	return fmt.Sprintf("s:%s:%s", svr, key)
}

func LoadSvrData(ctx context.Context, key string, data interface{}) error {
	key = paddingKey(key, appconfig.GlobalConfig().Global.SectionID)
	return loadData(ctx, key, data)
}

func SaveSvrData(ctx context.Context, key string, value interface{}) error {
	key = paddingKey(key, appconfig.GlobalConfig().Global.SectionID)
	return saveData(ctx, key, value)
}

func DeleteSvrData(ctx context.Context, key string) error {
	key = paddingKey(key, appconfig.GlobalConfig().Global.SectionID)
	return deleteData(ctx, key)
}

func loadData(ctx context.Context, key string, value interface{}) error {
	// check data is a pointer
	if reflect.TypeOf(value).Kind() != reflect.Ptr {
		return errors.New("data is not a pointer")
	}

	data, err := defaultBackend.Get(ctx, key)
	if err != nil {
		return err
	}

	if data == nil {
		return nil
	}
	if len(data) == 0 {
		return nil
	}
	err = defaultCoder.Unmarshal(data, value)
	if err != nil {
		return err
	}

	return nil
}

func saveData(ctx context.Context, key string, value interface{}) error {
	data, err := defaultCoder.Marshal(value)
	if err != nil {
		return err
	}
	err = defaultBackend.Set(ctx, key, data)
	if err != nil {
		return err
	}

	return nil
}

func deleteData(ctx context.Context, key string) error {
	return defaultBackend.Del(ctx, key)
}
