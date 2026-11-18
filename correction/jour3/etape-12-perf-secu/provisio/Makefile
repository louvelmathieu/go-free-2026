.PHONY: test race cover mocks lint

test:            ## tests unitaires
	go test ./...

race:            ## avec le détecteur de races
	go test -race ./...

cover:           ## couverture par package + rapport HTML
	go test -coverprofile=cover.out ./...
	go tool cover -func=cover.out | tail -1
	go tool cover -html=cover.out -o cover.html

mocks:           ## régénère les mocks (mockery v3, config dans .mockery.yml)
	mockery

synctest:        ## tests à temps virtuel (stable depuis Go 1.25)
	go test ./internal/orchestrator -run Synctest -v

lint:
	go vet ./...
