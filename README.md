# ranch53

Simple service that updates Route 53 based on labels on Rancher services.

## Installation

(not yet ready)

    $ go install github.com/flaccid/ranch53

## Usage

(todo)

### Build

#### Static Binary

To build a fully static 64-bit Linux binary (useful for docker):

    $ CGO_ENABLED=0 \
      GOOS=linux \
      GOARCH=amd64 \
        go build -a -installsuffix cgo -o bin/ranch53 .

#### Docker

Build the static binary per above, then the docker image with:

    $ docker build -t flaccid/ranch53 .

Example usage:

(todo: this requires needed env vars)

    $ docker run -it flaccid/ranch53

### Run

    $ go run main.go

License and Authors
-------------------
- Author: Chris Fordham (<chris@fordham-nagy.id.au>)

```text
Copyright 2017, Chris Fordham

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```
