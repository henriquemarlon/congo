.PHONY: infra
infra:
	$(START_LOG)
	@docker compose \
	-f ./deployments/compose.infra.yaml \
	up --build -d
	$(END_LOG)