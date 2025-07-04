<div align="center">

# <img src="assets/fpis.png" alt="Fernando Pessoa API" height="40"> ARQPI.ORG

[![Live API](https://img.shields.io/badge/LIVE_API-111d3b?style=for-the-badge&logoColor=white)](https://arqpi-org.onrender.com)
[![Frontend](https://img.shields.io/badge/EXPLORE-c7f0ff?style=for-the-badge&logoColor=111d3b)](https://arqpi.org)
[![Ko-fi](https://img.shields.io/badge/Support-ffdfdf?style=for-the-badge&logo=ko-fi&logoColor=111d3b)](https://ko-fi.com/bxav)

**Fernando Pessoa fragments API - Text corpus of Portugal's literary genius**

────

**4079 fragments** • **Memory-optimized** • **Full-text search** • **Rate limited**

────

![REST API](https://img.shields.io/badge/REST_API-c7f0ff?style=flat-square&logoColor=111d3b)
![Go](https://img.shields.io/badge/Golang-94ffab?style=flat-square&logo=go&logoColor=111d3b)
![Full-text Search](https://img.shields.io/badge/Full--text_Search-bbff94?style=flat-square&logoColor=111d3b)
![Heteronyms](https://img.shields.io/badge/Heteronyms-ffdfdf?style=flat-square&logoColor=111d3b)
![Fernando Pessoa](https://img.shields.io/badge/Fernando_Pessoa-111d3b?style=flat-square&logoColor=white)

</div>

## Overview

arqpi.org serves fragments of text from Fernando Pessoa and his heteronyms through a lightweight, fast API. The entire dataset (4079 fragments) is loaded in memory for efficient access and search.

> ⚠️ If your request returns a 502, retry after ~1 min as the free service may need to wake up!

## Endpoints

<div align="center">

| Endpoint | Description | Example |
|----------|-------------|---------|
| `/fragment?id=123` | Get fragment by ID | [Try it](https://arqpi-org.onrender.com/fragment?id=123) |
| `/random` | Return random fragment | [Try it](https://arqpi-org.onrender.com/random) |
| `/search?q=term` | Full-text search | [Try it](https://arqpi-org.onrender.com/search?q=dream) |
| `/info` | API metadata | [Try it](https://arqpi-org.onrender.com/info) |
| `/quote` | Return quotable fragment | [Try it](https://arqpi-org.onrender.com/quote) |
| `/status` | API status and stats | [Try it](https://arqpi-org.onrender.com/status) |

</div>

<div align="center">

────

### Rate Limits

**Free tier**: 5 requests per minute
**Supporters**: Unlimited access

────

</div>

## API Keys for Supporters

Thank you for supporting arqpi.org through [Ko-fi](https://ko-fi.com/bxav)! As a token of appreciation, supporters receive an unlimited-use API key that bypasses rate limits.

### How It Works:

1. After your donation, Ko-fi will immediately display your unique API key
2. This key is manually activated within 24 hours (or sooner!)
3. Once activated, add `?key=YOUR_API_KEY` to any API request
4. Enjoy unlimited access to the API!

> I use a manual activation process to ensure security and proper record-keeping. If you haven't received your activation confirmation within 24 hours, please reach out!

<div align="center">

────

### Frontend

A minimalist interface is available at [arqpi.org](https://arqpi.org) for easy exploration of the API.

<img src="https://raw.githubusercontent.com/bxavaby/arqpi-org/main/docs/assets/screenshot.png" alt="Fernando Pessoa API Interface" width="650">

────

**MIT License © 2025 bxavaby**

</div>
