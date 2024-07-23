present:
	@present -notes

run:
	@go run *.go

update-direct-access-model:
	@echo "Update OpenFGA Model to usable JSON format!"
	@cd openfga/models && openfga-cli model transform --file direct-access.fga>direct-access.json

start-server:
	@echo "Start server hosting app to check ..."
	@cd cmd/authz && go run *.go
	# @cd cmd/kilcron && go run main.go

start-kilcron:
	@echo "Start kilcron .."
	@cd cmd/kilcron && go run debug.go worker.go main.go

start-openfga:
	@echo "Start OpenFGA Server"
	@openfga run

start-temporal:
	@echo "Start Temporal Server"
	@temporal server start-dev --http-port 9090 --ui-port 3001 --metrics-port 9091 --log-level info

stop:
	@kill `pgrep openfga`	
	@kill `pgrep temporal`

health:
	@curl -X GET ${FGA_API_URL}/healthz	

