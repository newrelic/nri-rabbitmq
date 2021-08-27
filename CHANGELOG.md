# Change Log

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/)
and this project adheres to [Semantic Versioning](http://semver.org/).

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
