.PHONY: mocks test

mocks:
	mockgen -source=cache/cache.go -destination=cache/cache_mock.go -package=cache
	mockgen -source=cache/store/store.go -destination=cache/store/store_mock.go -package=store
	mockgen -source=logger/logger.go -destination=logger/logger_mock.go -package=logger

test:
	cd core; GOGC=10 go test -v -p=4 ./...
