# Change Log

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/)
and this project adheres to [Semantic Versioning](http://semver.org/).


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
