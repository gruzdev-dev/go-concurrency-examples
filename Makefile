.PHONY: verify

verify:
	@chmod +x scripts/verify.sh
	@echo "=== Test 1: unsafe ==="
	@./scripts/verify.sh configs/race.json
	@echo ""
	@echo "=== Test 2: mutex_copy ==="
	@./scripts/verify.sh configs/static_trap.json
	@echo ""
	@echo "=== Test 3: interface_tearing ==="
	@./scripts/verify.sh configs/tearing.json
