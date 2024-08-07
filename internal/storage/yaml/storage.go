package yaml

import (
	"bytes"
	"github.com/kirsle/configdir"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"os"
)

type Storage[T any] struct {
	path string
}

func NewStorage[T any](appName string, modelFile string) (*Storage[T], error) {
	path := configdir.LocalConfig(appName)
	if err := configdir.MakePath(path); err != nil {
		return nil, errors.Wrap(err, "failed to create model directory")
	}

	return &Storage[T]{
		path: path + modelFile,
	}, nil
}

func (s *Storage[T]) Get() (*T, error) {
	file, err := os.OpenFile(s.path, os.O_RDONLY|os.O_CREATE, 0755)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create model file")
	}
	defer func() {
		_ = file.Close()
	}()

	buf := bytes.NewBuffer(nil)
	if _, err = buf.ReadFrom(file); err != nil {
		return nil, errors.Wrap(err, "failed to read model file")
	}

	var model T
	if err = yaml.Unmarshal(buf.Bytes(), &model); err != nil {
		return nil, errors.Wrap(err, "failed to decode model file")
	}

	return &model, nil
}

func (s *Storage[T]) Save(model *T) error {
	file, err := os.OpenFile(s.path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return errors.Wrap(err, "failed to open model file")
	}
	defer func() {
		_ = file.Close()
	}()

	modelBytes, err := yaml.Marshal(model)
	if err != nil {
		return errors.Wrap(err, "failed to encode model file")
	}

	if _, err = file.Write(modelBytes); err != nil {
		return errors.Wrap(err, "failed to write model file")
	}

	return nil
}
