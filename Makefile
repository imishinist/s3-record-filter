
.PHONY: build
build:
	GOOS=linux GOARCH=amd64	go build -o build/print/print ./handler/print
	GOOS=linux GOARCH=amd64	go build -o build/filter/filter ./handler/filter

.PHONY: package
package:
	aws cloudformation package \
		--template-file simple-pipeline.yaml \
		--s3-bucket $(BUCKET) \
		--s3-prefix s3-record-filter \
		--output-template-file .template.yaml

.PHONY: deploy
deploy:
	./cloudformation.py deploy \
		--stack-name s3-record-filter \
		--template-path .template.yaml \
		--params-path simple-pipeline.jsonnet
