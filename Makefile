.PHONY: test deploy

NAME=dnsimple-dk-hostmaster-ds-upload
ENTRY_POINT=Handle

env.yaml:
	lpass show 538627301036416249 --notes --quiet --color=never > $@

test: *.go
	go test ./...

deploy: env.yaml test
	gcloud functions deploy $(NAME) --entry-point=$(ENTRY_POINT) --runtime=go111 --trigger-http --memory=128M --region=europe-west1 --env-vars-file=env.yaml
