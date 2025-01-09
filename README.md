# MPC (Multi-Party Computation) Project

A distributed key generation and signing system using threshold signatures with Ethereum integration.

## Overview

This project implements a Multi-Party Computation (MPC) system that enables distributed key generation and signing operations. It uses threshold signatures where a minimum number of parties must cooperate to generate signatures, enhancing security through decentralization.

Link to front-end (React Native) repo: [https://github.com/vietddude/mpcoin-fe](https://github.com/vietddude/mpcoin-fe)

## Features

- Distributed Key Generation (DKG)
- Threshold Signature Scheme (TSS)
- Ethereum Integration
- Real-time Transaction Monitoring
- RESTful API Interface
- Redis-based Session Management
- PostgreSQL Database Storage

## Prerequisites

- Go 1.22.2 or higher
- PostgreSQL
- Redis
- Access to an Ethereum node (for blockchain integration)

## Installation

1. Clone the repository:

```bash
git clone https://github.com/vietddude/mpcoin-be
cd mpc
```

2. Install dependencies:

```bash
go mod download
```

3. Set up environment variables (example):

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=mpc_db
REDIS_URL=localhost:6379
ETH_NODE_URL=wss://ethereum-sepolia-rpc.publicnode.com
```

## Project Structure

- `/cmd`
  - `/api` - API server implementation
  - `/test` - TSS testing implementation
  - `/worker` - Blockchain monitoring worker
- `/internal`
  - `/api` - API handlers and routing
  - `/config` - Configuration management
  - `/db` - Database operations
  - `/service` - Business logic layer
  - `/repository` - Data access layer
- `/pkg` - Reusable packages

## API Documentation

The API documentation is available via Swagger UI at:

```
http://localhost:5001/swagger/index.html
```

## Usage

1. Start the API server:

```bash
go run cmd/api/main.go
```

2. Start the blockchain worker:

```bash
go run cmd/worker/main.go
```

## Security

This project implements threshold signatures where `t` out of `n` parties must cooperate to generate valid signatures, providing security through decentralization.
