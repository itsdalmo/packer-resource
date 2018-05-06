package models_test

import (
	"encoding/json"
	"github.com/itsdalmo/packer-resource/src/models"
	"reflect"
	"strings"
	"testing"
)

func TestSource(t *testing.T) {

	input := strings.TrimSpace(`
{
  "aws_access_key_id": "",
  "aws_secret_access_key": "",
  "aws_session_token": "",
  "aws_region": ""
}
`)

	t.Run("unmarshalling yields expected output", func(t *testing.T) {
		expected := models.Source{
			AWSAccessKeyID:     "",
			AWSSecretAccessKey: "",
			AWSSessionToken:    "",
			AWSRegion:          "",
		}

		var actual models.Source
		err := json.Unmarshal([]byte(input), &actual)

		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("expected: '%v' instead got: '%v'", expected, actual)
		}
	})
}
