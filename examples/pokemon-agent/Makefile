help:
	@echo "Makefile targets:"

setup:
	cd emulator-agent && pip install -r requirements.txt
	cd llm-orchestrator && pip install -r requirements.txt

test:
	cd emulator-agent && pytest tests/
	cd llm-orchestrator && pytest tests/

lint:
	cd emulator-agent && flake8 .
	cd llm-orchestrator && flake8 .

interfaces-install:
	astartectl utils install-interface --astarte-url $(ASTRATE_URL) --realm-key $(REALM_KEY) --realm $(REALM) astrate-interfaces/*.json

run-emulator-stub:
	cd emulator-agent && python -m emulator_agent.main --stub --realm $(REALM) --device-id $(DEVICE_ID) --cert $(CERT) --key $(KEY) --ca $(CA)

run-orchestrator:
	cd llm-orchestrator && python -m llm_orchestrator.main

clean:
	find . -type d -name "__pycache__" -exec rm -rf {} +
	find . -type d -name ".pytest_cache" -exec rm -rf {} +
