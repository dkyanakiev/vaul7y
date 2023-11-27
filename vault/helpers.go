package vault

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/hashicorp/vault/api"
)

// extractListData reads the secret and returns a typed list of data and a
// boolean indicating whether the extraction was successful.
func extractListData(secret *api.Secret) ([]interface{}, bool) {
	if secret == nil || secret.Data == nil {
		return nil, false
	}

	k, ok := secret.Data["keys"]
	if !ok || k == nil {
		return nil, false
	}

	i, ok := k.([]interface{})
	return i, ok
}

// func OutputList(ui cli.Ui, data interface{}) int {
// 	switch data := data.(type) {
// 	case *api.Secret:
// 		secret := data
// 		return outputWithFormat(ui, secret, secret.Data["keys"])
// 	default:
// 		return outputWithFormat(ui, nil, data)
// 	}
// }

func ParseSecret(r io.Reader) (*api.Secret, error) {
	// First read the data into a buffer. Not super efficient but we want to
	// know if we actually have a body or not.
	var buf bytes.Buffer

	// io.Reader is treated like a stream and cannot be read
	// multiple times. Duplicating this stream using TeeReader
	// to use this data in case there is no top-level data from
	// api response
	var teebuf bytes.Buffer
	tee := io.TeeReader(r, &teebuf)

	_, err := buf.ReadFrom(tee)
	if err != nil {
		return nil, err
	}
	if buf.Len() == 0 {
		return nil, nil
	}

	// First decode the JSON into a map[string]interface{}
	var secret api.Secret
	dec := json.NewDecoder(&buf)
	dec.UseNumber()
	if err := dec.Decode(&secret); err != nil {
		return nil, err
	}

	// If the secret is null, add raw data to secret data if present
	if reflect.DeepEqual(secret, api.Secret{}) {
		data := make(map[string]interface{})
		dec := json.NewDecoder(&teebuf)
		dec.UseNumber()
		if err := dec.Decode(&data); err != nil {
			return nil, err
		}
		errRaw, errPresent := data["errors"]

		// if only errors are present in the resp.Body return nil
		// to return value not found as it does not have any raw data
		if len(data) == 1 && errPresent {
			return nil, nil
		}

		// if errors are present along with raw data return the error
		if errPresent {
			var errStrArray []string
			errBytes, err := json.Marshal(errRaw)
			if err != nil {
				return nil, err
			}
			if err := json.Unmarshal(errBytes, &errStrArray); err != nil {
				return nil, err
			}
			return nil, fmt.Errorf(strings.Join(errStrArray, " "))
		}

		// if any raw data is present in resp.Body, add it to secret
		if len(data) > 0 {
			secret.Data = data
		}
	}

	return &secret, nil
}
func DataIterator(t interface{}) {
	switch reflect.TypeOf(t).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(t)

		for i := 0; i < s.Len(); i++ {
			fmt.Println(s.Index(i))
		}
	}
}
