include .env

export
AWS_DEFAULT_PROFILE=${env}-string
API=string-api
ECS_CLUSTER=string-core
SERVICE_TAG=${tag}
AWS_REGION=us-west-2
ECR=${AWS_ACCT}.dkr.ecr.us-west-2.amazonaws.com
ECS_API_REPO=${ECR}/${API}

all: build push deploy

test-envvars:
	@[ "${env}" ] || ( echo "env var is not set"; exit 1 )
	@[ "${tag}" ] || ( echo "env tag is not set"; exit 1 )

build: test-envvars
	GOOS=linux GOARCH=amd64 go build -o ./tmp/app ./cmd/app/main.go
	docker build --platform linux/amd64 -t $(ECS_API_REPO):${SERVICE_TAG} .

push: test-envvars
	aws ecr get-login-password --region $(AWS_REGION) | docker login --username AWS --password-stdin $(ECS_API_REPO)
	docker push $(ECS_API_REPO):${SERVICE_TAG}

deploy: test-envvars
	aws ecs --region $(AWS_REGION) update-service --cluster $(ECS_CLUSTER) --service ${API} --force-new-deployment