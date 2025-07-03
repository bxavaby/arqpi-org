## Overview

arqpi-org serves fragments of text from Fernando Pessoa and his heteronyms through a lightweight, fast API. The entire dataset (4079 fragments) is loaded in memory for efficient access and search.

## Endpoints

- `/fragment?id=123` - Get a specific fragment by ID
- `/random` - Return a random fragment
- `/search?q=term` - Full text search across all fragments
- `/info` - API metadata (fragments count, authors, categories)
- `/quote` - Return a short, quotable fragment
- `/status` - API status and statistics

## API Keys for Supporters

Thank you for supporting arqpi.org through [Ko-fi](https://ko-fi.com)! As a token of appreciation, supporters receive an unlimited-use API key that bypasses rate limits.

### How It Works:

1. After your donation, Ko-fi will immediately display your unique API key
2. This key is manually activated within 24 hours (usually sooner!)
3. Once activated, add `?key=YOUR_API_KEY` to any API request
4. Enjoy unlimited access to the API!

> I use a manual activation process to ensure security and proper record-keeping. If you haven't received your activation confirmation within 24 hours, please reach out!
