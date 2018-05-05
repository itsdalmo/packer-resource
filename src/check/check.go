package check

import (
	"github.com/itsdalmo/packer-resource/src/models"
)

// Run (business logic)
func Run(request models.CheckRequest) (models.CheckResponse, error) {
	var response models.CheckResponse

	// We cannot know of any other versions since we build arbitrary templates.
	response = append(response, request.Version)

	return response, nil
}
