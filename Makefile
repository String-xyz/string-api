include .env.deploy

export
AWS_DEFAULT_PROFILE=${env}-string
ECR=${${env}_AWS_ACCT}.dkr.ecr.us-west-2.amazonaws.com

API=string-api
ECS_CLUSTER=string-core
ECS_API_REPO=${ECR}/${API}
SERVICE_TAG=${tag}

ECS_SANDBOX_CLUSTER=core-sandbox
SANDBOX_API=sandbox-string-api
ECS_SANDBOX_API_REPO=${ECR}/${SANDBOX_API}

all: build push deploy
all-sandbox: build-sandbox push-sandbox deploy-sandbox

test-envvars:
	@[ "${env}" ] || ( echo "env var is not set"; exit 1 )
	@[ "${tag}" ] || ( echo "env tag is not set"; exit 1 )

build: test-envvars
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./cmd/app/main ./cmd/app/main.go
	docker build --platform linux/amd64 -t $(ECS_API_REPO):${SERVICE_TAG} cmd/app/
	rm cmd/app/main

push: test-envvars
	aws ecr get-login-password --region $(AWS_REGION) | docker login --username AWS --password-stdin $(ECR)
	docker push $(ECS_API_REPO):${SERVICE_TAG}

deploy: test-envvars
	aws ecs --region $(AWS_REGION) update-service --cluster $(ECS_CLUSTER) --service ${API} --force-new-deployment

build-sandbox: test-envvars
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./cmd/app/main ./cmd/app/main.go
	docker build --platform linux/amd64 -t $(ECS_SANDBOX_API_REPO):${SERVICE_TAG} cmd/app/
	rm cmd/app/main

push-sandbox: test-envvars
	aws ecr get-login-password --region $(AWS_REGION) | docker login --username AWS --password-stdin $(ECR)
	docker push $(ECS_SANDBOX_API_REPO):${SERVICE_TAG}

deploy-sandbox: test-envvars
	aws ecs --region $(AWS_REGION) update-service --cluster $(ECS_SANDBOX_CLUSTER) --service ${SANDBOX_API} --force-new-deployment
