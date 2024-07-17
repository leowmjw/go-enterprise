run:
	@go run *.go


demo-setup:
	@echo "Start running OpenFGA + setup Store with Examples .."

demo:
	@echo "Running demo ..."	
	# Setup OpenFGA

start-server:
	@echo "Start server hosting app to check ..."
	@cd cmd/authz && go run *.go
	# @cd cmd/kilcron && go run main.go

start-openfga:	
	@openfga run

start-temporal:
	@echo "Start Temporal Server"
	@temporal server start-dev --http-port 9090 --ui-port 3001 --metrics-port 9091 --log-level info

test-direct-access:
	@echo "Demo Direct Access .."
	@openfga-cli model transform --file direct-access.fga>direct-access.json

stop:
	@kill `pgrep openfga`	
	@kill `pgrep temporal`

health:
	@curl -X GET ${FGA_API_URL}/healthz	

