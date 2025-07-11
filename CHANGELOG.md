# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/)
and this project adheres to [Semantic Versioning](http://semver.org/).

Unreleased section should follow [Release Toolkit](https://github.com/newrelic/release-toolkit#render-markdown-and-update-markdown)
## Unreleased

## v2.15.1 - 2025-07-02

### ⛓️ Dependencies
- Updated golang version to v1.24.4

## v2.15.0 - 2025-03-12

### 🚀 Enhancements
- Add FIPS compliant packages

### ⛓️ Dependencies
- Updated golang patch version to v1.23.6

## v2.14.2 - 2025-01-29

### ⛓️ Dependencies
- Updated golang patch version to v1.23.5

## v2.14.1 - 2024-12-04

### ⛓️ Dependencies
- Updated golang patch version to v1.23.4

## v2.14.0 - 2024-10-09

### dependency
- Upgrade go to 1.23.2

### 🚀 Enhancements
- Upgrade integrations SDK so the interval is variable and allows intervals up to 5 minutes

## v2.13.9 - 2024-09-11

### ⛓️ Dependencies
- Updated golang version to v1.23.1

## v2.13.8 - 2024-07-10

### ⛓️ Dependencies
- Updated golang version to v1.22.5

## v2.13.7 - 2024-06-26

### 🐞 Bug fixes
- we were logging an error message even if there was no actual error causing a weird `\u003cnil\u003e` to show up in the logs.

## v2.13.6 - 2024-05-15

### ⛓️ Dependencies
- Updated golang version to v1.22.3

## v2.13.5 - 2024-04-17

### ⛓️ Dependencies
- Updated golang version to v1.22.2

## v2.13.4 - 2024-03-13

### 🐞 Bug fixes
- Updated golang to version v1.21.7 to fix a vulnerability

## v2.13.3 - 2024-02-21

### ⛓️ Dependencies
- Updated github.com/newrelic/infra-integrations-sdk to v3.8.2+incompatible

## v2.13.2 - 2023-11-01

### ⛓️ Dependencies
- Updated golang version to 1.21

## v2.13.1 - 2023-08-02

### ⛓️ Dependencies
- Updated golang to v1.20.7

## v2.13.0 - 2023-07-26

### 🚀 Enhancements
- bumped golang version pinning 1.20.6

## 2.12.0  (2023-06-08)
### Added
- Update Go version to 1.20

## 2.11.0  (2023-04-20)
### Added
- Exposed the 'QueuesMaxLimit' flag to configure the max amount of queues collected.
- Added a 'DisableEntities' flag to avoid generating inventory entries on large environments

## 2.10.1  (2023-04-05)
### Changed
- Fix log example file path in unix tarballs
- Fix: handle huge values in node's disk_free metric properly
- Disable CGO

## 2.10.0 (2023-03-08)
### Changed
- Upgrade Go to 1.19 and bump dependencies

## 2.9.0 (2023-02-27)
### Added
- Add an arguments that allow to set the timeout to connecto to Rabbit's API

## 2.8.0 (2023-02-20)
### Changed
- Remove old Healthcheck and use the 'running' metric to calculate node status event.

## 2.7.0  (2022-09-27)
### Added
- Logging configuration examples files.


## 2.6.0  (2022-07-14)
### Breaking

- Removing the Cluster entity generation. This entity didn't contain any metrics related. And its deprecation was announced on this [EOL](https://discuss.newrelic.com/t/q1-bulk-eol-announcement-fy23/181744)

## 2.5.1  (2022-07-04)

### Changed
- Bump dependencies
### Added
Added support for more distributions:
- RHEL(EL) 9
- Ubuntu 22.04

## 2.5.0  (2022-05-03)
### Changed
- Move tool deps to go.mod in tools
- Update pipeline to Go 1.18
## Breaking
- Replace the attribute `clusterName` to `rabbitmqClusterName` to avoid collisions with the `clusterName` attribute reported when running in k8s. Naming was taken from [HAproxy integration fix](https://github.com/newrelic/nri-haproxy/blob/master/src/collection.go#L160) . User that have been use `clusterName` attribute will need to replace it with `rabbitmqClusterName`.

- Adds the `host:port` to all entity keys in order to use the entityRewrite when running in k8s. This fixes #73 . When the new version of the integration is deployed entities will be recreated with this new name. Old entities will live for an [extra day](https://github.com/newrelic/entity-definitions/blob/main/definitions/infra-rabbitmqqueue/definition.yml#L6).
example: 
`entityKey`(before):`ra-queue:/aliveness-test:clustername=rabbit@rabbitmq-0.rabbitmq-headless.rabbitmq.svc.cluster.local`
`entityKey`(fixed):`ra-queue:k8s:k8s-cluster-name:rabbitmq:pod:rabbitmq-0:rabbitmq:15672:/aliveness-test:clustername=rabbit@rabbitmq-0.rabbitmq-headless.rabbitmq.svc.cluster.local`

## 2.4.2  (2022-03-17)
### Added
- `rabbitmq-log.yml.example` is now in Linux packages to help setting up log parsing.

## 2.4.1 (2021-10-20)
### Added
Added support for more distributions:
- Debian 11
- Ubuntu 20.10
- Ubuntu 21.04
- SUSE 12.15
- SUSE 15.1
- SUSE 15.2
- SUSE 15.3
- Oracle Linux 7
- Oracle Linux 8

## 2.4.0 (2021-08-27)
### Added

Moved default config.sample to [V4](https://docs.newrelic.com/docs/create-integrations/infrastructure-integrations-sdk/specifications/host-integrations-newer-configuration-format/), added a dependency for infra-agent version 1.20.0

Please notice that old [V3](https://docs.newrelic.com/docs/create-integrations/infrastructure-integrations-sdk/specifications/host-integrations-standard-configuration-format/) configuration format is deprecated, but still supported.

## 2.3.1 (2021-06-10)
### Changed
- ARM support

## 2.3.0 (2021-05-10)
### Changed
- Update Go to v1.16.
- Migrate to Go Modules
- Update Infrastracture SDK to v3.6.7.
- Update other dependecies.

## 2.2.4 (2021-03-30)
### Added
- ARM and ARM64 packages for Linux
### Changed
- Exit rather than warn on API errors to avoid nil pointer error
- Moved release pipeline to Github Actions

## 2.2.3 (2020-07-28)
### Changed
- Increased queue limit to 2000

## 2.2.2 (2020-03-23)
### Changed
- Added argument `management_path_prefix` to support custom prefix for all HTTP request to the rabbitmq management plugin as detailed [here](https://www.rabbitmq.com/management.html#path-prefix).

## 2.2.1 (2020-01-28)
### Changed
- Send an inventory value when it would otherwise be empty

## 2.2.0 (2019-11-18)
### Changed
- Renamed the integration executable from nr-rabbitmq to nri-rabbitmq in order to be consistent with the package naming. **Important Note:** if you have any security module rules (eg. SELinux), alerts or automation that depends on the name of this binary, these will have to be updated.

## 2.1.2 - 2019-11-14
### Fixed
- Exclude windows definition from linux build

## 2.1.1 - 2019-10-16
### Fixed
- Windows installer GUIDs

## 2.1.0 - 2019-08-09
### Added
- Windows build files

## 2.0.5 - 2019-06-20
### Fixed
- Re-added clusterName as a queryable value

## 2.0.4 - 2019-06-12
### Fixed
- Exit code 69 error for rabbitmqctl

## 2.0.1 - 2019-05-20
### Fixed
- Segfault regression

## 2.0.0 - 2019-04-18
### Changed
- Changed entity keys so they are more likely unique
- Updated to v3 SDK
- Added reportingEntity attribute

## 1.0.4 - 2019-03-20
### Fixed
- Collect nodes when vhosts are filtered

## 1.0.3 - 2019-02-19
### Fixed
- Fixed bug where Queue whitelist would not work

## 1.0.2 - 2019-02-06
### Fixed
- Queue limiting was happening against full Queue list rather than filtered list

## 1.0.1 - 2019-02-05
### Fixed
- Updated protocol version

## 1.0.0 - 2018-11-16
### Changed
- Bumped version to 1.0.0

## 0.1.4 - 2018-11-02
### Changed
- Increased Queue limit to 500

## 0.1.3 - 2018-10-01
### Added
- Added limiting of Queue entities

## 0.1.2 - 2018-09-18
### Changed
- Changed sample file to be clearer for users to configure
- Correct misspellings

## 0.1.1 - 2018-08-24
### Fixed
- Exchange Binding Metric Data would show up as Queue Metric Data. While the bug showed exchange bindings as queue bindings, it would always show a count of zero too

## 0.1.0 - 2018-08-24
### Added
- Initial version: Includes Metrics, Inventory, and Events data
