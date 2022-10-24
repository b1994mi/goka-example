# Goka Example Usage

## Summary

An example of how a working back-end application that publishes (emit) & subscribes (consume/process) messages to and from kafka. It also uses protocol buffer (protobuf) to encode & decode messages to reduce messages size.

## Endpoints

| Endpoint      | HTTP   | Formdata Fields            | Body/Params                                                             |
| ------------- | ------ | ---------------------- | ---------------------------------------------------------------- |
| `/`    | GET    | Will give a simple response `{"message": "pong"}`  |                                                                  |
| `/wallet` | GET    | Get wallet data by Id   | Query params : `wallet_id, with_trx`. Example: `/wallet?wallet_id=6&with_trx=true`                                                                 |
| `/wallet`    | POST   | Create new wallet transaction data | Body (JSON): `wallet_id` (Number, required), `amount` (Number, required) |

## How to Start

1. Make sure you have installed docker and run the `docker-compose.yml` file to run `zookeeper` and `kafka`.
2. Enter command `go run .` from this directory to run the app (make sure to install go). 
