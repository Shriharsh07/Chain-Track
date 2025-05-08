# 🔗 Blockchain Ledger System

A simplified blockchain-based ledger application built with **Go**, **GORM**, **MySQL**, and **Gorilla Mux**. This project demonstrates how blockchain principles like immutability, hashing, and transaction verification can be implemented in a backend system.

---

## 🧠 Overview

This system allows users to submit transactions which are later grouped into blocks. Each block includes a hash of the previous block to ensure tamper resistance. This mimics how blockchains maintain a secure, verifiable chain of records.

---

## 🛠️ Tech Stack

- **Go (Golang)** — backend logic and concurrency
- **GORM** — ORM for MySQL
- **MySQL** — database for storing transactions and blocks
- **Gorilla Mux** — HTTP routing
- **Docker** *(optional)* — containerization for deployment

---

## ✨ Features

- Submit new transactions
- Automatically create new blocks after a set of transactions
- Link blocks using SHA-256 hashes
- Prevent tampering by enforcing chain integrity
- API access to view transactions and blocks

---

## 🧱 Architecture Diagram
<img src="resources/flowchart.png" alt="Blockchain Flow" width="300" height="450"/>

## 🚀 Getting Started

### Prerequisites

- Go 1.18+
- MySQL server
- Git

### Setup

```bash
# Clone the project
git clone https://github.com/Shriharsh07/Chain-Track.git
cd Chain-Track

# Install dependencies
go mod tidy
```

### 🔐 Configure Environment Variables
Create a .env file in the project root with the following content:
```
DB_USER=your_db_username
DB_PASS=your_db_password
DB_NAME=database_name //Example : chaintrack
DB_HOST=localhost
DB_PORT=3306
```

### Run the project
```
go run cmd/main.go
```


