# MiniBank

This is a simple Go-based API for a mini banking application.

## Features

* Create accounts
* Get account details by ID
* Get all accounts
* Transfer money between accounts (basic implementation)
* Secure authentication with JWT

## Getting Started

1. **Prerequisites:**
   - Go installed on your system
   - PostgreSQL database running

2. **Clone the repository:**

```bash
   git clone https://github.com/Nathene/MiniBank.git
```

Set up the database:

Create a PostgreSQL database named postgres.
Ensure the database user postgres has the necessary permissions to create tables and insert data.
### Build and run the application:

```Bash
make build 
make run
```

The API server will start running on port 3000.

## API Endpoints

* **`/login` (GET):**  Displays a login form.
* **`/login` (POST):**  Authenticates a user and returns a JWT.
* **`/account` (GET):**  Retrieves all accounts.
* **`/account` (POST):** Creates a new account.
* **`/account/{id}` (GET):** Retrieves an account by ID.
* **`/account/{id}` (DELETE):** Deletes an account by ID.
* **`/transfer` (POST):**  Transfers money between accounts.

### Running Tests

```bash
make test
```


### Seeding the Database
```Bash
make seed
```


### Removing the Account Table

```Bash
make rmtable
```


Authentication
The API uses JWT (JSON Web Tokens) for authentication. To access protected endpoints (e.g., /account/{id}), you need to include a valid JWT in the x-jwt-token header of your request.


## Example

```Go
package main

import (
    "log"

    "github.com/Nathene/MiniBank/internal"
    "github.com/Nathene/MiniBank/pkg/api"
)

func main() {
    store, err := internal.NewPostGresStore()
    if err != nil {
        log.Fatal(err)
    }

    if err := store.Init(); err != nil {
        log.Fatal(err)
    }

    server := api.NewAPIServer(":3000", store)
    server.Run()
}
```