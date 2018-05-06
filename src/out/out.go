package out

import (
	"bufio"
	"fmt"
	"github.com/itsdalmo/packer-resource/src/models"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"
)

// Run (business logic)
func Run(request models.PutRequest, sourceDir string) (models.PutResponse, error) {
	var response models.PutResponse

	if err := request.Source.Validate(); err != nil {
		return response, fmt.Errorf("invalid configuration: %s", err)
	}
	if err := request.Params.Validate(); err != nil {
		return response, fmt.Errorf("invalid parameters: %s", err)
	}

	if request.Source.AWSAccessKeyID != "" {
		os.Setenv("AWS_ACCESS_KEY_ID", request.Source.AWSAccessKeyID)
	}
	if request.Source.AWSSecretAccessKey != "" {
		os.Setenv("AWS_SECRET_ACCESS_KEY", request.Source.AWSSecretAccessKey)
	}
	if request.Source.AWSSessionToken != "" {
		os.Setenv("AWS_SESSION_TOKEN", request.Source.AWSSessionToken)
	}
	if request.Source.AWSRegion != "" {
		os.Setenv("AWS_DEFAULT_REGION", request.Source.AWSRegion)
	}

	p := &packer{
		Dir:    sourceDir,
		Params: request.Params,
		Wrt:    os.Stderr,
	}

	if err := p.Validate(); err != nil {
		return response, fmt.Errorf("packer validate failed: %s", err)
	}

	ami, err := p.Build()
	if err != nil {
		return response, fmt.Errorf("packer build failed: %s", err)
	}

	response.Version = models.Version{ImageID: ami}
	response.Metadata = make([]models.Metadata, 0)
	return response, nil
}

type packer struct {
	Dir    string
	Params models.PutParameters
	Wrt    io.Writer
}

func (p *packer) Arguments() []string {
	var args []string

	for k, v := range p.Params.Variables {
		args = append(args, fmt.Sprintf("-var=%s=%s", k, v))
	}
	if v := p.Params.VarFile; v != "" {
		args = append(args, "-var-file="+path.Join(p.Dir, v))
	}
	args = append(args, path.Join(p.Dir, p.Params.Template))

	return args
}

func (p *packer) Start(command string, rawOutput bool) (*exec.Cmd, io.Reader, error) {
	args := []string{command}
	if rawOutput {
		args = append(args, "-machine-readable")
	}

	// Append user arguments
	args = append(args, p.Arguments()...)

	// Run the command
	cmd := exec.Command("packer", args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return cmd, nil, fmt.Errorf("failed to open stdout pipe: %s", err)
	}
	if err := cmd.Start(); err != nil {
		return cmd, stdout, fmt.Errorf("failed to start cmd: %s", err)
	}
	return cmd, stdout, nil
}

func (p *packer) Validate() error {
	cmd, stdout, err := p.Start("validate", false)
	if err != nil {
		return err
	}

	go func() {
		s := bufio.NewScanner(stdout)
		for s.Scan() {
			fmt.Fprintln(p.Wrt, s.Text())
		}
	}()

	return cmd.Wait()
}

func (p *packer) Build() (string, error) {
	var ami string

	cmd, stdout, err := p.Start("build", true)
	if err != nil {
		return "", err
	}

	go func() {
		s := bufio.NewScanner(stdout)
		for s.Scan() {
			line := s.Text()
			if strings.Contains(line, ",ui,") {
				printUI(p.Wrt, line)
				continue
			}
			if strings.Contains(line, ",artifact,0,id,") {
				ami = parseAMI(line)
				continue
			}
		}
	}()

	if err := cmd.Wait(); err != nil {
		return "", err
	}

	return ami, nil
}

func parseAMI(s string) string {
	// 1525617757,amazon-ebs,artifact,0,id,eu-west-1:ami-2f6c5359
	p := strings.Split(s, ",")
	o := p[len(p)-1]
	// eu-west-1:ami-2f6c5359
	p = strings.Split(o, ":")
	o = p[len(p)-1]
	return o
}

// Based on: https://github.com/hashicorp/packer/blob/997f8e4a2ac3403446b46b5456c122684ce41210/packer/ui.go#L279-L307
func printUI(w io.Writer, s string) {
	p := strings.Split(s, ",")
	if len(p) != 5 {
		fmt.Fprintf(w, "[WARN] Failed to parse output with length %d: %v\n", len(p), p)
	}
	s = p[4]
	s = strings.Replace(s, "#!(PACKER_COMMA)", ",", -1)
	s = strings.Replace(s, "\\r", "\r", -1)
	s = strings.Replace(s, "\\n", "\n", -1)
	fmt.Fprintln(w, s)
}
