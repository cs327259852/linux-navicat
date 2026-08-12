# LinuxNavicat

[![Java CI with Maven](https://github.com/your-username/linux-navicat/actions/workflows/ci.yml/badge.svg)](https://github.com/your-username/linux-navicat/actions)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](LICENSE)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](CONTRIBUTING.md)
[![Language](https://img.shields.io/badge/Language-Java_8%2B-orange.svg)](https://www.java.com)

A lightweight, cross-platform GUI database client for **MySQL** built with Java Swing and JDBC. Designed to be a fast, open-source alternative to Navicat on Linux, macOS, and Windows.

[中文文档 (Chinese Readme)](README_ZH.md)

---

## ✨ Features

- **🚀 Lightweight & Fast**: Launches instantly with minimal RAM usage.
- **🌐 Cross-Platform**: Runs seamlessly on Linux, macOS, and Windows.
- **🔌 Connection Management**: Save multiple MySQL connections with custom host, port, username, and password.
- **📊 Table & Data Viewer**: 
  - Double-click connections and tables to browse structure and rows.
  - Page-by-page table pagination support.
  - One-click DDL (Data Definition Language) viewer.
- **📝 SQL Query Editor**: Execute DDL and DML statements with tabular query results.
- **🌍 Internationalization (i18n)**: Supports English, 简体中文, 繁体中文, and 日本语.
- **🎨 Modern UI**: Powered by FlatLaf for a sleek, native OS feel.

---

## 🛠️ Quick Start

### Prerequisites
- **JDK 8** or higher
- **Maven 3.6+**

### Build and Run

1. **Clone the repository:**
   ```bash
   git clone https://github.com/your-username/linux-navicat.git
   cd linux-navicat
   ```

2. **Run locally using Maven:**
   ```bash
   mvn clean exec:java
   ```

3. **Package as a Runnable Jar:**
   ```bash
   mvn clean package
   java -jar target/linux-navicat-0.2.0.jar
   ```

---

## 🗺️ Project Roadmap & Go Refactoring Plan

We are progressively refactoring `LinuxNavicat` to **Go + Wails + Vue 3 / TS**! Check out our [Progressive Daily Refactoring Roadmap](REFACTORING_GO_ROADMAP.md).

- [ ] **SQL Syntax Highlighting**: Integrate `RSyntaxTextArea` for SQL auto-completion and syntax highlighting.
- [ ] **Data Export / Import**: Support exporting query results to CSV, JSON, and Excel.
- [ ] **SSH Tunnel Support**: Support connecting to remote databases via SSH tunnels (JSch).
- [ ] **Multi-Database Support**: Extend JDBC abstraction to support PostgreSQL, SQLite, and MariaDB.
- [ ] **Query History**: Save executed SQL history for quick re-execution.
- [ ] **Dark Mode Theme**: Add a toggle for FlatLaf Dark Theme.

Check out our [Good First Issues](https://github.com/your-username/linux-navicat/issues?q=is%3Aissue+is%3Aopen+label%3A%22good+first+issue%22) to get started!

---

## 🤝 Contributing

Contributions make the open-source community an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**!

Please read our [Contributing Guide](CONTRIBUTING.md) to learn how you can report bugs, request features, or submit pull requests.

---

## 📄 License

Distributed under the Apache License 2.0. See [`LICENSE`](LICENSE) for more information.
