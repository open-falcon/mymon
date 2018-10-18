# This how we want to name the binary output
BINARY=mymon
GO_VERSION_MIN=1.10

# Add mysql version for testing `MYSQL_VERSION=5.7 make docker`
# use mysql:latest as default
MYSQL_VERSION := $(or ${MYSQL_VERSION}, ${MYSQL_VERSION}, latest)

.PHONY: all
all: fmt build

# Format code
.PHONY: fmt
fmt:
	@echo "\033[92mRun gofmt on all source files ...\033[0m"
	@echo "gofmt -l -s -w ..."
	@ret=0 && for d in $$(go list -f '{{.Dir}}' ./... | grep -v /vendor/); do \
                gofmt -l -s -w $$d/*.go || ret=$$? ; \
        done ; exit $$ret

# Run all test cases
.PHONY: test
test:
	go test `go list ./... | grep -Ev '/fixtures|/vendor/'`

.PHONY: cover
cover:
	go test `go list ./... | grep -Ev '/fixtures|/vendor/'` \
	-coverpkg=./... -coverprofile=coverage.data ./... | column -t
	go tool cover -html=coverage.data -o coverage.html
	go tool cover -func=coverage.data -o coverage.txt
	@tail -n 1 coverage.txt | awk '{sub(/%/, "", $$NF); \
		if($$NF < 80) \
			{print "\033[91m"$$0"%\033[0m"} \
        else if ($$NF >= 90) \
        	{print "\033[92m"$$0"%\033[0m"} \
        else \
        	{print "\033[93m"$$0"%\033[0m"}}'

# Builds the project
.PHONY: build
build:
	@bash ./genver.sh $(GO_VERSION_MIN)
	go build -o ${BINARY}

# Installs our project: copies binaries
.PHONY: install
install:
	go install

# Cleans our projects: deletes binaries
.PHONY: clean
clean:
	$(info rm -f ${BINARY})
	@if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
	@rm -f coverage.* innodb_mymon.* ./fixtures/innodb_mymon.* ./fixtures/process_mymon.* ./fixtures/*.log
	@for GOOS in darwin linux windows; do \
		for GOARCH in 386 amd64; do \
			rm -f ./release/${BINARY}.$${GOOS}-$${GOARCH} ;\
		done ;\
	done
	@find . -name "innodb_*" -delete
	@find . -name "process_*" -delete
	@find . -name "*.log" -delete
	@find . -name "push_metric.txt" -delete
	@docker stop mymon-master 2>/dev/null || true
	@docker stop mymon-slave 2>/dev/null || true
	@kill -9 `lsof -ti tcp:1988` 2>/dev/null || true

.PHONY: lint
lint:
	gometalinter.v1 --config metalinter.json ./...

.PHONY: release
release:
	@echo "\033[92mCross platform building for release ...\033[0m"
	@bash ./genver.sh $(GO_VERSION_MIN)
	@for GOOS in darwin linux windows; do \
		for GOARCH in 386 amd64; do \
			GOOS=$${GOOS} GOARCH=$${GOARCH} go build -v -o ./release/${BINARY}.$${GOOS}-$${GOARCH} 2>/dev/null ; \
		done ; \
	done

.PHONY: test_env
test_env:
	@echo "\033[92mBuild mysql test enviorment\033[0m"
	@docker stop mymon-master 2>/dev/null || true
	@docker stop mymon-slave 2>/dev/null || true
	@kill -9 `lsof -ti tcp:1988` 2>/dev/null || true

	@docker run --name mymon-master --rm -d \
	-e MYSQL_ROOT_PASSWORD=1tIsB1g3rt \
	-e MYSQL_DATABASE=mysql \
	-p 3306:3306 \
	-v `pwd`/fixtures/setup_master.sql:/docker-entrypoint-initdb.d/setup_master.sql \
	mysql:$(MYSQL_VERSION) \
	--server_id=1
	@echo -n "waiting for master initializing "
	@while ! mysql -h 127.0.0.1 -uroot -P3306 -p1tIsB1g3rt -NBe "do 1;" 2>/dev/null; do \
	printf '.' ; sleep 1 ; done ; echo '.'
	@echo "mysql master environment is ready!"

	@echo "begin to set slave environment: "
	@docker run --name mymon-slave --rm -d --link=mymon-master:mymon-master \
	-e MYSQL_ROOT_PASSWORD=1tIsB1g3rt \
	-e MYSQL_DATABASE=mysql \
	-p 3308:3306 \
	-v `pwd`/fixtures/setup_slave.sql:/docker-entrypoint-initdb.d/setup_slave.sql \
	mysql:$(MYSQL_VERSION) \
	--server_id=2 \
	--read_only
	@echo -n "waiting for slave initializing "
	@while ! mysql -h 127.0.0.1 -uroot -P3308 -p1tIsB1g3rt -NBe "do 1;" 2>/dev/null; do \
	printf '.' ; sleep 1 ; done ; echo '.'
	@echo "mysql slave environment is ready!"

	@echo "begin to set falcon push environment: "
	go run fixtures/push_mock.go &

.PHONY: daily
daily: fmt test_env test cover lint clean
	@echo "\033[93mdaily build successed! \033[m"
