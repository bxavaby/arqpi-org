# arqpi-org

A REST API providing programmatic access to Fernando Pessoa's works.

## Overview

arqpi-org serves fragments of text from Fernando Pessoa and his heteronyms through a lightweight, fast API. The entire dataset (4079 fragments) is loaded in memory for efficient access and search.

## Endpoints

- `/fragment?id=123` - Get a specific fragment by ID
- `/random` - Return a random fragment
- `/search?q=term` - Full text search across all fragments
- `/info` - API metadata (fragments count, authors, categories)
- `/quote` - Return a short, quotable fragment
- `/status` - API status and statistics

## Usage

For free, unlimited access to the API, consider supporting the project on [Ko-fi](https://ko-fi.com).
