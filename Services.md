# Services

List of services and a resume of it responsabilities

## Wallet

This is the core service, responsible for the main logic of managing the wallets,
and registering the transactions.

The service should provide APIs for Wallet CRUD operation and also to perform
transactions. Since it knows the core database, it should also be responsible for
returning analytics data to populate dashboards.

## Notifier

This service is responsible for real time updates and notifications. It should
handle websocket connections for live updates and should be able to send e-mail
notifications.