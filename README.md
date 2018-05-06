##  Packer resource (work in progress)

Concourse resource for Packer that is very much based on
[this resource](https://github.com/jdub/packer-resource) by the same name. A 
new resource was created because the old one has not been maintained in the
last 12 months, and the old resource had a tendency to hang.

## Source Configuration

- `aws_access_key_id`: *Optional*: Access key id if you are passing credentials.
- `aws_secret_access_key`: *Optional*: See above.
- `aws_session_token`: *Optional*: Use if your access/secret keys are temporary (assumed role/MFA authenticated).
- `aws_region`: Region where the images of interest live.

## Behaviour

#### `check`

Not implemented (will just return whatever version was passed to check).

#### `get`

Fetches additional metadata about the AMI, in addition to two files:

- `id`: Plain text file with the AMI ID.
- `packer.json`: Packer friendly variable file: `{"source_ami": "<ami-id>"}`.

(I.e. it has the same functionality as the ami-resource.)

#### `put`

Build an image using packer after specifing the following parameters:

- `template`: Path to the Packer template.
- `var_file`: *Optional*: Path to [external JSON variable file](https://www.packer.io/docs/templates/user-variables.html).
- `variables`: *Optional*: A map (name: value) of variables that will be passed to Packer.

## Example

The following (incomplete) example would build a new AMI using Packer:

```yaml
resource_types:
- name: packer
  type: docker-image
  source:
    repository: itsdalmo/packer-resource

resources:
- name: concourse-ami
  type: packer
  source:
    aws_access_key_id: ((aws-access-key))
    aws_secret_access_key: ((aws-secret-key))
    aws_session_token: ((aws-session-token))
    aws_region: eu-west-1

jobs:
- name: bake-concourse
  plan:
  - put: concourse-ami
    params:
      template: concourse.json
      var_file:
      - packer/variables.json
      variables:
        environment_tag: development
```
