# data-ingestion-go

Simple data generator for parts and suppliers, written in Go.

## Overview
This project generates synthetic data for parts and suppliers, useful for testing, demos, and populating development databases.

## Features
- Generate random parts and suppliers with realistic fields
- Output data to CSV files
- Easily configurable and extendable

## Getting Started

1. Clone the repository:
	```sh
	git clone https://github.com/bitterfq/data-ingestion-go.git
	cd data-ingestion-go
	```
2. Run the generator:
	```sh
	go run .
	```
3. Run tests:
	```sh
	go test ./...
	```

## Project Structure
- `parts/` — Logic for generating part data
- `suppliers/` — Logic for generating supplier data
- `internal/db/` — Database models and queries (auto-generated)
- `schema.sql` — Database schema
- `queries.sql` — SQL queries for data operations
