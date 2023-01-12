# UserAPI
An restful api that uses CRUDL to support user endpoints. Stores the users in mysqlite. Creates 10 users when the program starts. 

## How to start the API
First you need to create a folder and call it whatever you want. Then change your directory to that folder you just created.
Next you want to clone the github repo using this command

```bash
git clone https://github.com/ntrut/UserAPI
```
Make sure your still in the directory you created. Next go into UserAPI folder that was just cloned from github. Also, you need to have go installed on your machine. In your terminal run this command
```bash
go run .
```
The program runs on localhost with port 8080, just make sure nothing is being used on that port. Also, you dont need to do anything with the db file. The program will automatically overwrite it when the program first runs. The program will always start with the same 10 users. Now you are ready to use the API!

## API Docs

### GET Requests
All http GET Requests. The read request will throw a 404 NOT FOUND if the id doesnt exist. Id needs to be a number that is creater than 0 else it will throw a 400 BAD Reqeust. IF the parameters are incorrect for the list request, it will throw a 400 BAD REQUEST.
```bash
#Returns one user based on the id
http://localhost:8080/infor/read/{id}
```

```bash
#Returns all users from the database
http://localhost:8080/infor/list
```

```bash
#Returns all users from the database in either asc or desc order and can be sorted by first name or last name or email or id.
http://localhost:8080/infor/list?order={asc or desc}&sort={user parameter}
```

### POST Requests
This is how the body needs to look like when sending the POST request. This http request will throw 409 CONFLICT if the id already exists in the database.
```json
{
    "id": 11,
    "email": "helloworld@gmail.com",
    "first_name": "hello",
    "last_name": "world"
}
```
Example on how it works with curl
```bash
#Creates one user
curl -X POST http://localhost:8080/infor/create -H 'Content-Type: application/json' -d '{"id":11,"email":"helloworld@gmail.com","first_name":"hello","last_name":"world"}'
#response returns the user created with the timestamp as well
-> {"id":11,"email":"helloworld@gmail.com","first_name":"hello","last_name":"world","updated":"2021-11-23 14:00:11.463186244 -0800 PST m=+919.601937982"}

```
### PUT Requests
This is how the body needs to look like when sending the PUT request. We update the email first name and last name. This http request will throw 404 NOT FOUND if user with that id doesnt exist. Id needs to be a number that is greater than 0 else it will throw a 400 BAD Reqeust.
```json
{
    "email": "helloworld@gmail.com",
    "first_name": "hello",
    "last_name": "world"
}
```
```bash
#Updates one user with that id
http://localhost:8080/infor/update/{id}
```

### DELETE Reqeusts
Simple delete request. Will throw a 404 NOT FOUND if the id doesnt exist. Id needs to be a number that is creater than 0 else it will throw a 400 BAD Reqeust.
```bash
#Deletes one user with that id
http://localhost:8080/infor/delete/{id}
```
