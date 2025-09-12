#/bin/bash

sudo nitro-cli terminate-enclave --all

docker build -f Dockerfile.enclave -t safe-code-enclave .

nitro-cli build-enclave --docker-uri safe-code-enclave --output-file safe-code-enclave.eif

sudo nitro-cli run-enclave --cpu-count 1 --memory 512 --eif-path safe-code-enclave.eif --debug-mode