## Prerequisites

- [Docker Desktop](https://docs.docker.com/desktop/)
- [Migrate CLI](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)
## Installation

Follow these steps to set up the project locally:

1. **Clone the repository**:
```bash
git clone https://github.com/SonikSeven/password-manager.git
```

2. **Navigate to the project folder**:
```bash
cd password-manager
```

3. **Start the Docker containers** (make sure [Docker Desktop](https://docs.docker.com/desktop/) is running):
```bash
docker-compose up -d
```

4. **Run database migrations** using the [Migrate CLI](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate) tool:
```bash
.\tools\migrate.exe -path db/migration -database "postgresql://postgres:postgres@localhost:5432/password_manager?sslmode=disable" -verbose up
```

5. **Access the app**:
Open your browser or [Postman](https://www.postman.com/) and navigate to:
```
http://localhost:8000
```
You can now interact with the API endpoints listed below.
## Endpoints

| Endpoint              | HTTP Method & Access                         | JSON Fields                                                  | URL Parameters   |
| --------------------- | -------------------------------------------- | ------------------------------------------------------------ | ---------------- |
| `/api/register`       | `POST`: anyone                               | `email`\*<br>`password`\*                                    |                  |
| `/api/login`          | `POST`: anyone                               | `email`\*<br>`password`\*                                    |                  |
| `/api/passwords`<br>  | `GET`, `POST`: authenticated user            | `username`\*<br>`password`\*<br>`url`\*<br>`notes`<br>`icon` | search<br>domain |
| `/api/passwords/<id>` | `GET`, `PATCH`, `DELETE`: authenticated user | `username`\*<br>`password`\*<br>`url`\*<br>`notes`<br>`icon` |                  |
### **Note**

- The included `app.env` is for **demonstration purposes only**. It contains placeholder credentials and settings to help you run the app locally. **Do not use it in production.**
