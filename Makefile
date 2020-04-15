list-bucket: $(glob ./cmd/list-bucket/*.go)
	go build ./cmd/list-bucket

init:
	sqlite3 db.sqlite3 <./schema.sql

download: list-bucket
	./prowjob.sh >prowjob.json
	./filter-prowjob--failure-url.sh <./prowjob.json | while read -r url; do ./download-from-url.sh "$$url"; done

download-release: list-bucket
	./prowjob.sh >prowjob.json
	./filter-prowjob--failure-url.sh <./prowjob.json | grep /release | while read -r url; do ./download-from-url.sh "$$url"; done

download-aws-operator: list-bucket
	./prowjob.sh >prowjob.json
	./filter-prowjob--failure-url.sh <./prowjob.json | grep 'cluster-image-registry-operator.*e2e-aws-operator' | while read -r url; do ./download-from-url.sh "$$url"; done

update:
	go build ./cmd/record
	find output -type f | while read -r file; do ./denoise-2 <$$file | ./record $$file; done
