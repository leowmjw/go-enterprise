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

start-openfga:	
	@openfga run

start-temporal:
	@temporal server start-dev

test-direct-access:
	@echo "Demo Direct Access .."
	@openfga-cli model transform --file direct-access.fga>direct-access.json

stop:
	@kill `pgrep openfga`	
	@kill `pgrep temporal`

health:
	@curl -X GET ${FGA_API_URL}/healthz	

