Here's an updated version of the `README.md` file, including the additional information about using the key-value store API for your Lua package manager "nep":

```markdown
# Lua Package Manager - Key-Value Store API

This repository contains a key-value store API implemented in Go, using SQLite as the database backend. This API is part of a Lua package manager system that allows you to store and retrieve key-value pairs.

## Prerequisites

- Go installed on your machine.
- SQLite installed.
- Git for cloning repositories.
- Postman or any other API testing tool.

## Getting Started

1. **Clone the Repository**

   Clone this repository to your local machine:

   ```sh
   git clone https://github.com/your-username/your-repo.git
   cd your-repo
   ```

2. **Install Dependencies**

   Install the necessary Go packages:

   ```sh
   go get github.com/mattn/go-sqlite3
   ```

3. **Run the Server**

   Start the Go server:

   ```sh
   go run main.go
   ```

   The server should now be running at `http://localhost:8080`.

4. **Testing with Postman**

   - **Adding a key-value pair:**
     - Open Postman.
     - Set the request type to `POST`.
     - Enter the URL: `http://localhost:8080/post`.
     - Go to the `Params` tab.
     - Add two parameters: `key` and `value`.
       - For example:
         - `key`: `mypackage`
         - `value`: `username/repo`
     - Click on `Send`.

   - **Retrieving a value by key:**
     - Open Postman.
     - Set the request type to `GET`.
     - Enter the URL: `http://localhost:8080/get?key=mypackage`.
     - Click on `Send`.

## API Endpoints

### POST /post

Adds a key-value pair to the SQLite database.

#### Parameters

- `key`: The name you want to refer to when running `nep` in the terminal.
- `value`: The `<git username>/<repo name>` as found in the GitHub URL.

### GET /get

Retrieves the value associated with a specific key from the SQLite database.

#### Parameters

- `key`: The key for which you want to retrieve the value.

## Example Usage

### Adding a Key-Value Pair

```sh
curl -X POST "http://localhost:8080/post?key=mypackage&value=username/repo"
```

### Retrieving a Value by Key

```sh
curl "http://localhost:8080/get?key=mypackage"
```

## Notes

- The `key` should be the name you want to refer to when running `nep` in the terminal.
- The `value` should be in the format `<git username>/<repo name>`, as found in the GitHub URL.
- This example uses SQLite as the database. If you plan to use a different database, you will need to adjust the SQL driver and connection accordingly.
- Ensure your SQLite database file (`kvstore.db`) is accessible and writable by your Go application.

```
