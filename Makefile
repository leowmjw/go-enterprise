run:
	@go run *.go


demo-setup:
	@echo "Start running OpenFGA + setup Store with Examples .."

demo:
	@echo "Running demo ..."	
	# Setup OpenFGA

start-openfga:	
	@openfga run

start-temporal:
	@temporal server start-dev

stop:
	@kill `pgrep openfga`	
	@kill `pgrep temporal`

health:
	@curl -X GET ${FGA_API_URL}/healthz	

