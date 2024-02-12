package runtime

import (
	"encoding/json"
	"os"
)

type Runtime struct {
	cmds map[string]string
}

func NewRuntime(bindingFile string) (*Runtime, error) {
	jsonFile, err := os.Open(bindingFile)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()
	runtime := &Runtime{}
	err = json.NewDecoder(jsonFile).Decode(&runtime.cmds)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *Runtime) Invoke(cmd string) error {
	return nil
}
