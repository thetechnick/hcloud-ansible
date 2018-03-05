package ansible

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/pflag"
	"github.com/thetechnick/hcloud-ansible/pkg/version"
)

// Module interface for Ansible modules
type Module interface {
	Args() interface{}
	Run() (ModuleResponse, error)
}

// RunModule executes the module
func RunModule(m Module, flags *pflag.FlagSet) {
	var resp ModuleResponse
	if err := flags.Parse(os.Args[1:]); err != nil {
		resp.Msg(fmt.Sprintf("Error parsing flags: %v", err)).
			Failed().
			exitJSON()
	}
	if v, _ := flags.GetBool("version"); v {
		version.PrintText()
	}

	if len(os.Args) != 2 {
		resp.Msg("No arguments file provided").
			Failed().
			exitJSON()
	}

	argsFile := os.Args[1]
	argsString, err := ioutil.ReadFile(argsFile)
	if err != nil {
		resp.Msg(fmt.Sprintf("Cannot read arguments file: %v", err)).
			Failed().
			exitJSON()
	}

	if err := json.Unmarshal(argsString, m.Args()); err != nil {
		resp.Msg(fmt.Sprintf("Cannot parse arguments file: %v", err)).
			Failed().
			exitJSON()
	}

	if resp, err := m.Run(); err != nil {
		resp.Msg(err.Error()).
			Failed().
			exitJSON()
	} else {
		resp.exitJSON()
	}
}

// ModuleResponse represents the reponse of the module
type ModuleResponse struct {
	msg     string
	changed bool
	failed  bool
	data    map[string]interface{}
}

// Msg sets the module message
func (r *ModuleResponse) Msg(msg string) *ModuleResponse {
	r.msg = msg
	return r
}

// Changed marks the the module as changed
func (r *ModuleResponse) Changed() *ModuleResponse {
	r.changed = true
	return r
}

// Failed marks the module as failed
func (r *ModuleResponse) Failed() *ModuleResponse {
	r.failed = true
	return r
}

// Set adds data to the reponse
func (r *ModuleResponse) Set(key string, value interface{}) *ModuleResponse {
	if r.data == nil {
		r.data = map[string]interface{}{}
	}
	r.data[key] = value
	return r
}

// Data returns the data of the response
func (r *ModuleResponse) Data() map[string]interface{} {
	return r.data
}

// HasChanged returns true if the module has changed
func (r *ModuleResponse) HasChanged() bool {
	return r.changed
}

// HasFailed returns true if the module has failed
func (r *ModuleResponse) HasFailed() bool {
	return r.failed
}

// MarshalJSON custom json marshalling
func (r *ModuleResponse) MarshalJSON() ([]byte, error) {
	data := map[string]interface{}{
		"changed": r.changed,
		"failed":  r.failed,
		"msg":     r.msg,
	}

	for key, value := range r.data {
		data[key] = value
	}
	return json.Marshal(data)
}

func (r *ModuleResponse) exitJSON() {
	response, _ := json.Marshal(r)
	fmt.Println(string(response))
	if r.failed {
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}
