BIN_DIR=_output/bin
RELEASE_VER=0.2

kube-arbitrator: init
	go build -o ${BIN_DIR}/kube-batchd ./cmd/kube-batchd/

verify: generate-code
	hack/verify-gofmt.sh
	hack/verify-golint.sh
	hack/verify-gencode.sh

init:
	mkdir -p ${BIN_DIR}

generate-code:
	go build -o ${BIN_DIR}/deepcopy-gen ./cmd/deepcopy-gen/
	${BIN_DIR}/deepcopy-gen -i ./pkg/apis/scheduling/v1alpha1/ -O zz_generated.deepcopy

images: kube-arbitrator
	cp ./_output/bin/kube-batchd ./deployment/
	docker build ./deployment/ -f ./deployment/Dockerfile -t kubearbitrator/kube-batchd:${RELEASE_VER}
	rm -f ./deployment/kube-batchd

run-test:
	hack/make-rules/test.sh $(WHAT) $(TESTS)

e2e: kube-arbitrator
	hack/run-e2e.sh

coverage:
	KUBE_COVER=y hack/make-rules/test.sh $(WHAT) $(TESTS)

clean:
	rm -rf _output/
	rm -f kube-arbitrator
