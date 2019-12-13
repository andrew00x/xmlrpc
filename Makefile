SRC_DIRS=pkg

test:                 ## Run tests
	@go test -v $(addprefix ./, $(addsuffix /..., $(SRC_DIRS)))

help:                 ## Print this help
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/##//'