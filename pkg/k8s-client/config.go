package k8sclient

import (
	"github.com/pkg/errors"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"strings"
)

func ValidateK8sConfigFile(path string) error {
	stat, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return errors.New("file does not exist")
		}
		if os.IsPermission(err) {
			return errors.New("file path permission denied")
		}
		return err
	}
	if stat.IsDir() {
		return errors.New("file path must be a file, not a directory")
	}

	if !strings.HasSuffix(path, ".yaml") && !strings.HasSuffix(path, ".yml") {
		return errors.New("file path must be a yaml file")
	}

	if _, err = clientcmd.BuildConfigFromFlags("", path); err != nil {
		return errors.Wrap(err, "build k8s client config")
	}

	return nil
}
