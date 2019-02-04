# Tool commands (overridable)
DOCKER_CMD         ?= docker
DOCKER_COMPOSE_CMD ?= docker-compose

MAKEFILE_THIS          := $(lastword $(MAKEFILE_LIST))
THIS_PATH              := $(patsubst %/,%,$(dir $(abspath $(MAKEFILE_THIS))))
FIXTURE_CUSTOM_DOCKERENV_PATH    := $(THIS_PATH)/fixtures/custom-docker

# Docker-compose
BASE_DOCKER_COMPOSE_FILES := -f ./docker-compose.yaml

DOCKER_COMPOSE_UP_FLAGS            := --remove-orphans --force-recreate
DOCKER_COMPOSE_UP_TEST_FLAGS       := $(DOCKER_COMPOSE_UP_FLAGS) --abort-on-container-exit
DOCKER_COMPOSE_UP_BACKGROUND_FLAGS := $(DOCKER_COMPOSE_UP_FLAGS) -d
DOCKER_COMPOSE_UP_STANDALONE_FLAGS := $(DOCKER_COMPOSE_UP_FLAGS)
DOCKER_COMPOSE_PULL_FLAGS :=

export DOCKER_COMPOSE_CMD

.PHONY: dockerenv-custom-up
dockerenv-custom-up:
		cd $(FIXTURE_CUSTOM_DOCKERENV_PATH) && \
		$(DOCKER_COMPOSE_CMD) $(BASE_DOCKER_COMPOSE_FILES) up $(DOCKER_COMPOSE_UP_STANDALONE_FLAGS)

