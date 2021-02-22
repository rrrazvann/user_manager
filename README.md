### README ###

Using docker-compose run the application which is exposing port 4273, so you can test it using Postman.


The app is using json as content-type for the endpoints
There are 4 endpoints:
	- POST /users -> insert user into db (fields: first_name, last_name, nickname, password, email, country)
	- GET /users  -> find users by specific conditions (fields: first_name, last_name, nickname, password, email, country)
	- PUT /users/{id} -> update user with specific id with the specified fields (fields: first_name, last_name, nickname, password, email, country)
	- DELETE /users/{id} -> self explanatory

Users has the following fields: first_name, last_name, nickname, password, email, country.
Nickname and email are unique into db.
The data is directly added into db, the password needs to be hashed before calling the endpoints.

Improvements:
	- retry buffer for subscribed services for event notifications
	- endpoint validation (regex validation + ping it and expect a 200)
	- use a router which has route parameters implemented
	- storage system health check
	- add posibility of the server to be used in a cluster