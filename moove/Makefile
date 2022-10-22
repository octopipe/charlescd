BUF_VERSION:=1.1.0
PROTO_FILES := $(shell find ../butler -name '*.proto' | sort)

proto:
	for i in ${PROTO_FILES}; do \
		mkdir -p pb/v1 && cp $$i ./pb/v1 ; \
	done

	docker run -v $(shell pwd):/pb -w /pb --rm bufbuild/buf:$(BUF_VERSION) generate