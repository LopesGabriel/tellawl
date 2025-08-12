# Tellawl

This is a financial system that helps managing and analyse you finance health.

The user is able to create wallets in order to have a better control and share
with other users.

## Features

Main features listed below.

### Create a wallet

The user is able to create a wallet and optionally can share it with other users
so both could add income and outcome data.

### Register transaction

The user will be able to register a new transaction that can be of two types,
`Income` or `Outcome`.

This operation should trigger an event in the system, enabling us to notify all 
users who share the wallet in which the operation was performed.

### Dashboard showing financial health

The application web page should provide the user with a dashboard with historical
data about the wallet. It should contain the patrimonial evolution.

It will be possible to group the data by week, month, and year. The user should
have access to a select box in order to choose the desired option.

## Application

### Testing

We'll have support for unit, integration, and e2e tests. The unit and integration
tests can be triggered by running `yarn test`, and the e2e by running the `yarn test:e2e`.

### Observability

The application must be able to push traces, metrics, and optionally logs to a
OTEL collector.

## Infrastructure

The application should be containerized and should follow the cloudevents
specification for events communication.

The database and other resources should be provisioned using terraform and the
credentials must be provided to the applications using environment variables.