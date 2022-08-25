DIST_DIR=$(CURDIR)/dist

cli:
	go build -o ${DIST_DIR}/cndev main.go

clean:
	rm -rf ${DIST_DIR}

prod-build:
	hack/build $(DIST_DIR)