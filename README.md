# ✨ Pack Optimizer Application ✨

This application is designed to calculate the optimal distribution of item packs to fulfill a customer's order, adhering to a set of specific business rules. It includes a Golang backend with a RESTful API and a simple web-based UI for user interaction.

-----

## 💡 Project Explanation

The core functionality of this application is to solve a packing optimization problem. Given an order quantity and a set of available pack sizes, the service must determine the pack combination that:

1.  Fulfills the order with only **whole packs**.

2.  Uses the **least number of total items** (i.e., minimal overage).

3.  Uses the **fewest number of individual packs**.

For example, to fulfill an order of `251` items with pack sizes of `250` and `500`, the optimal solution is one `500`-item pack, as it minimizes the item overage compared to using two `250`-item packs.

The application is built with a flexible architecture, allowing new pack sizes to be added to the PostgreSQL database without requiring any code changes.

## 🏗️ Infrastructure and Architecture

### 🗂️ Clean Architecture

The project is structured according to the principles of Clean Architecture to ensure separation of concerns, testability, and maintainability.

* **`internal/domain`**: Contains the core business logic, including the `Pack` entity and the `PackRepository` interface.
* **`internal/usecase`**: Implements the business logic. The `PackUseCase` takes an order quantity and returns the optimal pack distribution by interacting with the `PackRepository` interface.
* **`internal/repository`**: The data access layer. The `sql_repo` package provides a concrete implementation of the `PackRepository` using GORM and PostgreSQL.
* **`internal/handler`**: The API layer, responsible for handling HTTP requests, calling the use case, and formatting responses.

### 💻 Technology Stack

* **Backend**: Golang with the **Fiber** web framework.
* **Frontend**: A simple HTML/CSS/JS frontend styled with **Tailwind CSS**.
* **Database**: **PostgreSQL** for pack size persistence.
* **Containerization**: **Docker** and **Docker Compose** for easy setup and deployment.
* **Migrations**: Database schema management is handled by **golang-migrate/migrate**.

## 🚀 Getting Started

### ✅ Prerequisites

* **Docker**
* **Docker Compose**
* **Go** (for local development and testing)
* **Python3 and pip or pipx** (for pre-commit)

### ▶️ How to Run

1.  **Clone the repository:**

    ```bash
    git clone https://github.com/kiarashfz/pack_optimizer
    cd pack_optimizer
    ```

2.  **Rename `.env.sample` file in the root of project to `.env`.**

    The original `.env` file is excluded from version control via `.gitignore` for security reasons. That's why it doesn't come with the repository when you clone it.


3.  **Start the services:**
    This command will build the Docker images, run database migrations, and start the application and database containers.

    ```bash
    make up
    ```

4.  **Access the application:**
    Open your browser and navigate to [http://localhost:8080](http://localhost:8080).

### 🛠️ Development Setup


#### 1️⃣ Install pre-commit

**Ubuntu:**

```bash
sudo apt install pipx
pipx ensurepath
pipx install pre-commit
```

**MacOS:**

```bash
brew install pre-commit
pipx ensurepath
pipx install pre-commit
```

#### 2️⃣ Install Golang tools
If you don't have the necessary Go tools installed, you can run:

```bash
make dev-setup
```
This command will install the required Go tools such as `goimports` and `golangci-lint` if they are not already installed.
#### 3️⃣ Install pre-commit hooks
To ensure code quality and consistency, the project uses pre-commit hooks. You can install them by running:

```bash
pre-commit install
```
This command sets up the pre-commit hooks defined in the `.pre-commit-config.yaml` file, which will automatically run checks on your code before each commit.

#### 4️⃣ lint your code
To lint your code, you can use the following command:

```bash
make lint
```



### 🔧 Makefile Functionality

The included `Makefile` provides convenient commands for managing the application's lifecycle.

* `make up`: Builds and starts all services in the background.
* `make down`: Stops and removes all containers, networks, and volumes.
* `make restart`: Restarts all services.
* `make logs`: Streams the logs for all services.
* `make test`: Runs all unit and integration tests with verbose output and race detection.
* `make dev-setup`: Installs necessary Go tools if missing.
* `make lint`: Runs `goimports` and `golangci-lint` on the codebase.

### 🟧 Postman Collection

A ready-to-use Postman collection is available in the `docs/api` directory.

You can import it directly into your Postman to test the API endpoints.



### 🧪 How to Test

All tests are located in the root of the project under the test directory. This includes both:

* **Unit Tests**: Located in the project root directory (e.g., `test/unit/usecase/pack_usecase_test.go`), these tests verify the functionality of individual components in isolation using mocks.


* **Integration Tests OR API Tests**: Found in the `test/integration` directory, these tests verify that the different layers of the application (handler, use case) work together correctly.

To run all tests, use the following command:

```bash
make test
```
### ☁️ Deployment

The application is successfully deployed and accessible at:
👉 [https://kiarash-pack-optimizer.up.railway.app/](https://kiarash-pack-optimizer.up.railway.app/)
