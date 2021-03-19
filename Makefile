install:
	# pout
	go install ./cmd/pout
examples:
	# python
	python3 -m grpc_tools.protoc -I./examples/protos --python_out=./examples/python/basic ./examples/protos/swapi/*
	# nodejs
	cd examples/nodejs/basic
	npm run gen:proto
