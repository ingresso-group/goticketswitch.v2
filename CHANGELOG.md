# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.0.6] - 2019-04-28
### Added
- Add possible discounts to avail response

## [1.0.5] - 2019-04-28
### Added
- Add price band description to availability

## [1.0.4] - 2019-02-15
### Added
- Added Cancel endpoint

## [1.0.3] - 2019-01-29
### Changed
- Moved AgentReference field from Customer to MakePurchase

## [1.0.2] - 2019-01-16
### Added
- Check for Callback Gone and Authentication errors

## [1.0.1] - 2018-08-24
### Fixed
- Close the response body after we're done with it

## [1.0.0] - 2018-07-11
### Changed
- Now each request method takes a context, see example/main.go for an example

## [0.0.2] - 2018-07-02
### Added
- GetDiscounts
- GetSendMethods
- MakeReservation
- ReleaseReservation
- MakePurchase
- GetStatus

## [0.0.1] - 2018-05-24
### Added
- This CHANGELOG file!
- basic client
- test.v1
- events.v1
- performances.v1
- availability.v1
- sources.v1

### Changed
- exposed DateRange utility

[Unreleased]: https://github.com/ingresso-group/goticketswitch.v2/compare/1.0.4...HEAD
