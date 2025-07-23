# LodgeX Backend Exercise 

Challenge for new new backend hires. If you have any questions please email liam@lodgex.com.au

# Prerequisites 

Everything here was designed to work on Linux. If you're on windows, it might just work if you have the correct tools, or you can just use WSL.

In order to run, make sure you have the following installed:
* Golang (1.23)
* Docker
* Docker Compose
* Make
* VS Code (Optional) 
* Docker Desktop (Optional)
* Dbeaver (Optional)
* Postman (Optional) 
* Ability to google

And to ensure everything works, run:

```
make build
```

# Running 

The server requires a database to run properly. To build and start the database, run:

```
make startdb
```

This should build and start a postgres database on port `5400` or as defined in `docker-compose.yml`.

If you wan't to make changes to the database schema and/or data, you can also do so by changing `database/schema.sql` and running `make dbrebuild`.

If you choose to use VS Code, you should be able to launch the backend directly via VS Code "Run and Debug" tab (assuming you have Golang installed). Otherwise, you can run `make run` which launches the server via a golang docker container.

# Project Directory 

Hopefully most of the project structure is fairly self-explanatory:

* `/server` is the actual webserver 
* `/server/handlers` has some "handlers" or "controllers" for the API endpionts. 
* `/server/queries` contains functions for direct database calls 
* `/docker` contains some dockerfiles 

# Tasks 

* Task 1: Ensure everything works. Start the database, then start the webserver. If things don't work, make them work. Call the `GET /users` endpoint to make sure the pre-populated data comes back as expected.

* Task 2: Create an endpoint to create a user.

* Task 3: Given the default database data, the endpoint for fetching users sometimes returns duplicates. Fix that.

* Task 4: Create an endpoint to create a user message, and then modify `GET /users` to return a message count. 

* Task 5: Create a PATCH request for modifying users. Should allow the removal of user nicknames (setting nickname=NULL in the database). 

* Task 6: Write some simple unit tests to cover fetching user data and creating users.

* Task 7 (Optional): Implement some level of JWT Authentication.

* Task 8 (Optional): Make any amount of refactors you wish. 


# Information 

The database uses the below information for connection:

```
POSTGRES_DB=exercise
POSTGRES_USER=super_secure_username
POSTGRES_PASSWORD=super_secure_password
```





