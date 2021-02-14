
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
	aws cloudformation deploy \
		--stack-name s3-record-filter \
		--template-file .template.yaml \
		--capabilities CAPABILITY_IAM \
		--parameter-overrides BatchSize=1
