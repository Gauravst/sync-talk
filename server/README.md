# Go API Template

![Version](https://img.shields.io/badge/version-1.0.5-blue)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

This is a production-ready Go backend API template designed to kickstart your next project.

## Features

- **Multiple Database Support**: Easily switch between different databases (e.g., PostgreSQL, MySQL, SQLite).
- **Docker Integration**: Run the application and its dependencies in isolated containers.
- **Makefile**: Simplify common tasks like building, testing, and running the application.
- **Database Migrations**: Manage database schema changes with ease.
- **Configuration Management**: Use `local.yaml` for environment-specific configurations.
- **Modular Structure**: Organized into `models`, `database`, `handlers`, `middleware`, `repositories`, and `services`.
- **Testing System**: Includes a robust testing framework for unit and integration tests.

## Setup

### 1. Create a New Repository

Click the **"Use this template"** button at the top right of this repository to create your own copy.

### 2. Configure the Application

Update the following files with your environment-specific settings:

- **`config/local.yaml`**: Add your application-specific configurations (e.g., server port, logging level).
- **`.env`**: Add your environment variables (e.g., database credentials, API keys).

### 3. Download Dependencies

`go mod download`

### 4. Run Your App

`make run`

### 5. Database Migrations

- To apply migrations: `make migrate-up`
- To rollback migrations: `make migrate-down`

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any improvements or bug fixes.

## License

This project is licensed under the **MIT License**. See the [LICENSE](LICENSE) file for details.
