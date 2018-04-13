package gojenkins

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"os"
)

// Credential represents a Jenkins scoped credential
type Credential struct {
	// Id is optional. If not specified on creation then Jenkins will generate one.
	ID string `json:"id"`
	// Scope is one of {GLOBAL} TODO: where to find it
	Scope string `json:"scope"`

	// Username is required when JavaClass is set to classUserPassword.
	Username    string `json:"username"`
	Password    string `json:"password"`
	Description string `json:"description"`

	// Depending on the implementation class, fields are required.
	JavaClass string `json:"$class"`
}

const (
	classUserPassword = "com.cloudbees.plugins.credentials.impl.UsernamePasswordCredentialsImpl"
)

// NewUsernamePasswordCredential returns a Credential for user:password
// It's description and id is empty which can be set.
func NewUsernamePasswordCredential(username, password string) *Credential {
	return &Credential{
		Scope:     "GLOBAL",
		Username:  username,
		Password:  password,
		JavaClass: classUserPassword,
	}
}

// CreateCredential creates a credential unless it exists. It does not update it.
func (j *Jenkins) CreateCredential(c *Credential) error {
	returnValue := map[string]interface{}{}
	buffer := new(bytes.Buffer)
	envelope := struct {
		Empty       string      `json:""`
		Credentials *Credential `json:"credentials"`
	}{
		Empty:       "0",
		Credentials: c,
	}
	json.NewEncoder(os.Stdout).Encode(envelope)
	if err := json.NewEncoder(buffer).Encode(envelope); err != nil {
		return err
	}
	resp, err := j.Requester.PostJSON(
		"/credentials/store/system/domain/_/createCredentials",
		buffer,
		&returnValue,
		map[string]string{})
	if err != nil {
		return err
	}
	log.Printf("%#v", resp)
	if resp.StatusCode == 200 {
		return nil
	}
	return errors.New("Unable to create credential")
}
