MAKEFILE_THIS := $(lastword $(MAKEFILE_LIST))
SRC_PATH      := $(patsubst %/,%,$(dir $(abspath $(MAKEFILE_THIS))))
APP_FILE      := mapper_app.json
APP_NAME      := mapper-app
LOG_LEVEL     := INFO
HTTP_PORT     := 8080

.PHONY: all
all: create build run

.PHONY: clean
clean:
	rm -Rf $(APP_NAME)

.PHONY: create
create: $(APP_FILE) clean
	flogo create -f $(APP_FILE) $(APP_NAME)

.PHONY: build
build: $(APP_NAME)
	cd $(APP_NAME) && flogo build -e

.PHONY: run
run: $(APP_NAME)/bin
	cd $(APP_NAME)/bin && FLOGO_LOG_LEVEL=$(LOG_LEVEL) HTTPPORT=$(HTTP_PORT) ./$(APP_NAME) 

.PHONY: tests
tests: first list

.PHONY: list
list:
	curl http://localhost:$(HTTP_PORT)/demo/list -H "Content-Type: applicaiton/json" -d "{\"company\":\"TIBCO\",\"address\":{\"street\":\"3033 Hillsview Ave\",\"city\":\"Palo Alto\",\"state\":\"CA\"} }"

.PHONY: first
first:
	curl http://localhost:$(HTTP_PORT)/demo/first -H "Content-Type: applicaiton/json" -d "{\"company\":\"TIBCO\",\"address\":{\"street\":\"3033 Hillsview Ave\",\"city\":\"Palo Alto\",\"state\":\"CA\"} }"
