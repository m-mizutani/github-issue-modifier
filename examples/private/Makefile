ifeq (,$(wildcard $(STACK_CONFIG)))
    $(error STACK_CONFIG ($(STACK_CONFIG)) is not found)
endif

CodeS3Bucket := $(shell cat $(STACK_CONFIG) | jq '.["CodeS3Bucket"]' -r )
CodeS3Prefix := $(shell cat $(STACK_CONFIG) | jq '.["CodeS3Prefix"]' -r )
StackName := $(shell cat $(STACK_CONFIG) | jq '.["StackName"]' -r )
Region := $(shell cat $(STACK_CONFIG) | jq '.["Region"]' -r )
SecretArn := $(shell cat $(STACK_CONFIG) | jq '.["SecretArn"]' -r )
GithubEndpoint := $(shell cat $(STACK_CONFIG) | jq '.["GithubEndpoint"]' -r )
LambdaSubnetIds := $(shell cat $(STACK_CONFIG) | jq '.["LambdaSubnetIds"]' -r )
LambdaSecurityGroupIds := $(shell cat $(STACK_CONFIG) | jq '.["LambdaSecurityGroupIds"]' -r )
VpcEndpointId := $(shell cat $(STACK_CONFIG) | jq '.["VpcEndpointId"]' -r )


TEMPLATE_FILE=template.yml

all: deploy

clean:
	rm build/main

build/main: *.go
	env GOARCH=amd64 GOOS=linux go build -o build/main .

sam.yml: $(TEMPLATE_FILE) build/main
	aws --region $(Region) cloudformation package \
		--template-file $(TEMPLATE_FILE) \
		--s3-bucket $(CodeS3Bucket) \
		--s3-prefix $(CodeS3Prefix) \
		--output-template-file sam.yml

deploy: sam.yml
	aws --region $(Region) cloudformation deploy \
		--template-file sam.yml \
		--stack-name $(StackName) \
		--capabilities CAPABILITY_IAM \
		--parameter-overrides \
		  SecretArn='$(SecretArn)' \
		  GithubEndpoint='$(GithubEndpoint)' \
		  LambdaSubnetIds='$(LambdaSubnetIds)' \
		  LambdaSecurityGroupIds='$(LambdaSecurityGroupIds)'
