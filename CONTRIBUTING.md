# Contributing to Aaxion

First off, thank you for considering contributing to Aaxion! It's people like you that make Aaxion such a great tool for secure, zero-buffer file streaming.

## 🧠 Core Philosophy

Aaxion is built to be **incredibly lightweight** and **efficient**. When writing code for Aaxion, keep memory allocation and disk I/O in mind. Features should degrade gracefully on low-end hardware.

## 🛠️ Development Environment

### Prerequisites

- **Go 1.22+**
- **Node.js 18+** (for the Next.js website/documentation)
- **Make** (optional, but recommended)

### Local Setup

1. Clone the repository:

   ```bash
   git clone https://github.com/codershubinc/aaxion.git
   cd aaxion
   ```

2. Build the binary:

   ```bash
   go build -o aax cmd/aaxion/main.go
   ```

3. Initialize the database and create a local admin:

   ```bash
   ./aax create-admin dev:devpass
   ```

4. Start the server:

   ```bash
   ./aax serve
   ```

### Website & Documentation

The Aaxion documentation site is built using Next.js.

```bash
cd website
npm install
npm run dev
```

## 📝 Commit Guidelines

We enforce conventional commits to generate automated release notes.
Format: `<type>(<scope>): <subject>`

**Types:**

- `feat`: A new feature
- `fix`: A bug fix
- `docs`: Documentation only changes
- `style`: Changes that do not affect the meaning of the code (formatting, missing semi-colons, etc)
- `refactor`: A code change that neither fixes a bug nor adds a feature
- `perf`: A code change that improves performance
- `chore`: Changes to the build process or auxiliary tools

**Example:**
`feat(auth): add --force flag to create-admin command`

## 🚀 Pull Request Process

1. Ensure any new features are thoroughly tested on both Linux and Windows (using the build tags when dealing with filesystem syscalls).
2. Update the `CHANGELOG.md` with notes of your changes under the `[Unreleased]` block.
3. Submit a PR against the `main` branch.
4. A maintainer will review your code. We may request performance benchmarks for heavy features.

## 🐛 Reporting Bugs

If you find a bug, please create a GitHub Issue with:

- Your OS and Architecture
- Aaxion Version (`aax version`)
- Steps to reproduce the issue
- Relevant logs (if using systemd, `journalctl -u aaxion.service -n 50`)

Thank you for contributing!
