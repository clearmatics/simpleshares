.PHONY: build clean fabric-network fabric-network-clean ethereum-network

# Tool commands (overridable)
DOCKER_CMD         ?= docker
DOCKER_COMPOSE_CMD ?= docker-compose

CHAINCODE_DIR = ./chaincode/src/simpleshares

clean: fabric-network-clean
	@rm -rf $(CHAINCODE_DIR)/vendor && \
		rm -rf build node_modules

fabric-network:
	@cd $(CHAINCODE_DIR) && \
		govendor init && \
		govendor add +external && \
		cd ../../../network && \
		make dockerenv-custom-up


fabric-network-clean:
	@cd ./network/fixtures/custom-docker && \
		$(DOCKER_COMPOSE_CMD) -f docker-compose.yaml down

ethereum-network:
	@npm install && \
		npm run testrpc