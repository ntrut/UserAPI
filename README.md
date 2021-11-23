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
http://localhost:8080/infor/list?order=asc&sort=first_name
```
