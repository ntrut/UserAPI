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
Now you are ready to use the API!

## API Docs

### GET Requests
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
This is how the body needs to look like when sending the post request
```json
{
    "id": 9,
    "email": "helloworld@gmail.com",
    "first_name": "hello",
    "last_name": "world"
}
```
Example on how it works with curl
```bash
#Creates one user
curl -X POST http://localhost:8080/infor/create -H 'Content-Type: application/json' -d '{"id":11,"email":"helloworld@gmail.com","first_name":"hello","last_name":"world"}'
#response
-> {"id":11,"email":"helloworld@gmail.com","first_name":"hello","last_name":"world","updated":"2021-11-23 14:00:11.463186244 -0800 PST m=+919.601937982"}

```
