# web_scraper

### Description
The project is a web application that will extract large amounts of data from the Google search results page. It is built using the Gin web framework, which is a lightweight framework for building web applications in Go.

### Configuration
The server port and database configuration can be modified in the `app.env` file. To change the server port, update the `HTTP_SERVER_ADDRESS` variable. To modify the database configuration, update the `DB_SOURCE` variable.

Note that changing the database configuration may require additional changes to the Makefile.

### Setup infrastructure
- Start postgres container:

    ```bash
    make postgres
    ```
    
- Create web_scraper database:

    ```bash
    make createdb
    ```

- Run db migration up all versions:

    ```bash
    make migrateup
    ```

- Run db migration down all versions:

    ```bash
    make migratedown
    ```

### How to run

- Run server:

    ```bash
    make server
    ```

- Run server using docker:

    ```bash
    make serverdocker
    ```

- Run test:

    ```bash
    make test
    ```

- Open your web browser and navigate to http://localhost:8080 to view the application.
