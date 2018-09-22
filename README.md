# contactbook

A contactbook backend application which illustrates CRUD APIs over REST.


### Build and install

```sh
$ go get -d github.com/prashanthpai/contactbook
$ cd $GOPATH/src/github.com/prashanthpai/contactbook
$ make
$ make install
```


### Features

* Each contact should have a unique email address.
* APIs support adding/editing/deleting contacts.
* Allows searching by name or email address.
* Search by name supports pagination (10 items by default per invocation)
* Uses an embedded DB

### Usage

```sh
Usage of contactbook:
  -addr string
    	Address to listen on for HTTP server. (default ":8080")
  -db-file string
    	Path to db file. (default "contactbook.db")
  -password string
    	Username of HTTP user. (default "livpo")
  -user string
    	Username of HTTP user. (default "ppai")
```

**Running server**
```sh
$ contactbook --user user --password password
```

**Create a contact**
```sh
$ cat create.json
{
    "name": "Some Name",
    "email": "some_email@example.com",
    "phone": "987654321"
}
```
```sh
$ curl --user user -i -X POST http://localhost:8080/contacts --data @create.json
Enter host password for user 'user':
HTTP/1.1 201 Created
Location: /contacts/some_email@example.com
Date: Sun, 23 Sep 2018 11:04:20 GMT
Content-Length: 0

```

**Delete a contact**
```sh
$ curl --user user -i -X DELETE http://localhost:8080/contacts/some_email@example.com
Enter host password for user 'user':
HTTP/1.1 204 No Content
Date: Sun, 23 Sep 2018 11:05:32 GMT

```

**Get a contact by email**
```sh
$ curl --user user -i -X GET http://localhost:8080/contacts?email=some_email@example.com
Enter host password for user 'user':
HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
Date: Sun, 23 Sep 2018 11:06:20 GMT
Content-Length: 178

[{"name":"Some Name","email":"some_email@example.com","phone":"987654321","created-at":"2018-09-23T16:36:06.804896664+05:30","updated-at":"2018-09-23T16:36:06.804896664+05:30"}]
```

**List contacts** (10 per page by default)
```sh
$ curl --user user -i -X GET http://localhost:8080/contacts
Enter host password for user 'user':
HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
Date: Sun, 23 Sep 2018 11:07:35 GMT
Content-Length: 356

[{"name":"Some Name","email":"some_email@example.com","phone":"987654321","created-at":"2018-09-23T16:36:06.804896664+05:30","updated-at":"2018-09-23T16:36:06.804896664+05:30"},{"name":"Some Name 2","email":"some_email2@example.com","phone":"9876543224","created-at":"2018-09-23T16:37:29.64788043+05:30","updated-at":"2018-09-23T16:37:29.64788043+05:30"}]

```

**Pagination**
```sh
$ curl --user user -i -X GET http://localhost:8080/contacts?page=1
$ curl --user user -i -X GET http://localhost:8080/contacts?name=Name&page=2
```
