# LinuxNavicat

[![Java CI with Maven](https://github.com/your-username/linux-navicat/actions/workflows/ci.yml/badge.svg)](https://github.com/your-username/linux-navicat/actions)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](LICENSE)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](CONTRIBUTING.md)
[![Language](https://img.shields.io/badge/Language-Java_8%2B-orange.svg)](https://www.java.com)

一款基于 Java Swing 和 JDBC 开发的轻量级、跨平台 **MySQL** 数据库图形客户端。旨在为 Linux、macOS 和 Windows 用户提供一个极速、开源免费的 Navicat 替代品。

[English Documentation](README.md)

---

## ✨ 核心特性

- **🚀 极其轻量**: 启动秒级响应，内存占用小。
- **🌐 跨平台支持**: 无缝运行于 Linux、macOS 和 Windows。
- **🔌 连接管理**: 支持保存多个 MySQL 连接配置（主机、端口、用户名、密码）。
- **📊 数据库与数据表浏览**:
  - 双击连接与数据表快速查看表结构与行数据。
  - 支持数据表翻页（分页查看）。
  - 一键查看建表 DDL 语句。
- **📝 SQL 查询编辑器**: 支持编写并执行 DDL 与 DML 语句，结果以表格形式直观展示。
- **🌍 多语言国际化 (i18n)**: 支持 简体中文、繁体中文、English 及 日本语。
- **🎨 现代化 UI**: 引入 FlatLaf 现代化外观主题，简洁顺滑。

---

## 🛠️ 快速开始

### 环境准备
- **JDK 8** 或更高版本
- **Maven 3.6+**

### 编译与运行

1. **克隆代码仓库:**
   ```bash
   git clone https://github.com/your-username/linux-navicat.git
   cd linux-navicat
   ```

2. **使用 Maven 快速运行:**
   ```bash
   mvn clean exec:java
   ```

3. **打包为可执行 Jar 包:**
   ```bash
   mvn clean package
   java -jar target/linux-navicat-0.2.0.jar
   ```

---

## 🗺️ 项目路线图 (Roadmap)

我们非常欢迎各种程度的开发者参与贡献！以下是项目当前规划的功能特性：

- [ ] **SQL 语法高亮与补全**: 集成 `RSyntaxTextArea` 实现 SQL 关键字高亮及智能补全。
- [ ] **数据导入与导出**: 支持将查询结果一键导出为 CSV、JSON 或 Excel 格式。
- [ ] **SSH 隧道连接**: 支持通过 SSH 隧道 (JSch) 安全连接远程数据库。
- [ ] **多数据库支持**: 抽象 JDBC 层，扩展支持 PostgreSQL、SQLite 及 MariaDB。
- [ ] **SQL 执行历史**: 保存已执行的 SQL 历史记录，方便快速重查。
- [ ] **暗黑模式 (Dark Mode)**: 支持一键切换 FlatLaf 暗黑主题。

欢迎查看我们的 [Good First Issues](https://github.com/your-username/linux-navicat/issues?q=is%3Aissue+is%3Aopen+label%3A%22good+first+issue%22) 快速上手参与！

---

## 🤝 参与贡献

社区的成长离不开每一位贡献者的付出！无论是提交 Bug 反馈、编写代码、改进文档还是提出建议，我们都非常欢迎。

请阅读 [贡献指南 (CONTRIBUTING.md)](CONTRIBUTING.md) 了解如何提交 Issue 和 Pull Request。

---

## 📄 开源许可

本项目基于 Apache License 2.0 协议开源，详情请参阅 [`LICENSE`](LICENSE) 文件。
