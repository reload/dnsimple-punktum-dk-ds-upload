.PHONY: test deploy logs clean doc

NAME=dnsimple-dk-hostmaster-ds-upload
ENTRY_POINT=Handle
REGION=europe-west1
RUNTIME=go111

export GO111MODULE=on

doc: README.md

README.md: *.go .godocdown.tmpl
	godocdown --output=README.md

env.yaml:
	lpass show 538627301036416249 --notes --quiet --color=never > $@

test: *.go
	go test ./...

deploy: env.yaml test
	gcloud functions deploy $(NAME) --entry-point=$(ENTRY_POINT) --runtime=$(RUNTIME) --trigger-http --memory=128M --region=$(REGION) --env-vars-file=env.yaml

logs:
	gcloud functions logs read $(NAME) --region=$(REGION) --limit=100

clean:
	$(RM) env.yaml
