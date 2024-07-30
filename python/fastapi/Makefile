#### DEFAULTS ##########################################################################################################
#### Boilerplate Makefile setup
MAKEFLAGS += --warn-undefined-variables
SHELL := bash
.SHELLFLAGS := -ec -o pipefail
.DEFAULT_GOAL := help
.DELETE_ON_ERROR:
.SUFFIXES:


#### INCLUDES ##########################################################################################################
#### Includes of other Makefiles, or files that declare environment variables (e.g. .env)
# include .env _should_ allow vars to be defined here
# However be careful of vars referencing other vars. or vals with $ in it (e.g. passwords), which would need escaping

env_file_exists := $(wildcard .env)
ifneq ($(strip ${env_file_exists}),)
include .env
endif

#### VARS ##############################################################################################################
#### Variables used in this Makefile.
#### Uppercased are environment vars, or make-specific vars. All others should be lower-snake-case

# default ENV_CODE to lcl if not set
ENV_CODE ?= lcl
venv_name := appstack-fastapi

#### RULES #############################################################################################################
#### Rules defining file targets that need to be made, or PHONY targets, which don't actually produce a file
#### Reminder: The name of non-PHONY targets needs to be the name of a file on disk, or it will always be invoked
#### NOTE: Most rules here deal with project-setup
####       Rules orchestrating project workloads are in the included Makefile

.PHONY: help
help: ## print this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' Makefile | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: printvars
printvars: ## Print the Environment variables present in calls to make, plus variables defined in the executed Makefile
	@$(info ==== Makefile Variables ====)
	@$(info )
	@$(foreach V,$(sort $(.VARIABLES)), \
		$(if $(filter-out environment% default automatic, \
		$(origin $V)),$(info $V=$($V) ($(value $V)))))
	@printf "\n==== Environment Variables ====\n\n"
	@printenv

.PHONY: check
check: ## Checks the makefile is working
	@printf 'make is working'

