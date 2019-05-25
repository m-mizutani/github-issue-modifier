ifeq (,$(wildcard $(STACK_CONFIG)))
    $(error STACK_CONFIG ($(STACK_CONFIG)) is not found)
endif

CODE_S3_BUCKET := $(shell cat $(STACK_CONFIG) | jq '.["CodeS3Bucket"]' -r )
CODE_S3_PREFIX := $(shell cat $(STACK_CONFIG) | jq '.["CodeS3Prefix"]' -r )
STACK_NAME := $(shell cat $(STACK_CONFIG) | jq '.["StackName"]' -r )
REGION := $(shell cat $(STACK_CONFIG) | jq '.["Region"]' -r )

SecretArn := $(shell cat $(STACK_CONFIG) | jq '.["SecretArn"]' -r )
GithubEndpoint := $(shell cat $(STACK_CONFIG) | jq '.["GithubEndpoint"]' -r )
Action := $(shell cat $(STACK_CONFIG) | jq '.["Action"]' -r )

TEMPLATE_FILE=template.yml

all: deploy

clean:
	rm build/main

build/main: api/*.go
	env GOARCH=amd64 GOOS=linux go build -o build/main .

sam.yml: $(TEMPLATE_FILE) build/main
	aws --region $(REGION) cloudformation package \
		--template-file $(TEMPLATE_FILE) \
		--s3-bucket $(CODE_S3_BUCKET) \
		--s3-prefix $(CODE_S3_PREFIX) \
		--output-template-file sam.yml

deploy: sam.yml
	aws --region $(REGION) cloudformation deploy \
		--template-file sam.yml \
		--stack-name $(STACK_NAME) \
		--capabilities CAPABILITY_IAM \
		--parameter-overrides \
		  SecretArn='$(SecretArn)' \
		  GithubEndpoint='$(GithubEndpoint)' \
		  Action='$(Action)'
