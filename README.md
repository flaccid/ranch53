# ranch53

A simple service that updates Route 53 based on labels on Rancher services.

## How does it work?

Currently the following two synchronisations are actioned:

 1. Load Balancer Services that have the labels, `dns_alias` and `dns_target`
    will have a CNAME record created or updated in that R53 zone in the format of
    `<dns_alias>. CNAME <dns_target>`

 2. Hosts that have the labels, `tier` and `pool_zone` will update an A record
    with each of their private IP addresses under the name, `rancher-<tier>-pool.<pool_zone>`

Limited use cases are supported currently. Please raise an issue so we can look at
how to best add other uses.

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

    $ docker run \
        -it \
        --env-file .env \
          flaccid/ranch53

##### Publish

    $ docker push flaccid/ranch53

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
