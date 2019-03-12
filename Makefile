DEPS=github.com/mpl/scgiclient \
	github.com/stretchr/testify \
	github.com/go-xmlfmt/xmlfmt

SRC_DIRS=pkg

deps:                 ## Get all dependencies
	@go get $(DEPS)

test: deps            ## Run tests
	@go test -v $(addprefix ./, $(addsuffix /..., $(SRC_DIRS)))

help:                 ## Print this help
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/##//'