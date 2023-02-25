gen-cal:
	protoc calculator/calculatorpb/calculator.proto --go_out=:.
	protoc calculator/calculatorpb/calculator.proto --go-grpc_out=:.
run-server:
	go run ./calculator/server/
run-client:
	go run ./calculator/client/