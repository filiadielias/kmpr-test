# kmpr-test
kumparan technical assessment

### How To Run This Project

```bash
# go get
go get github.com/filiadielias/kmpr-test

# go to directory
cd $GOPATH/src/github.com/filiadielias/kmpr-test

# run Project
go run app.go

```


## Project Information

Stack : Golang, PostgreSQL, ElasticSearch, NSQ (message queue), Redis (cache).

Hosted on AWS, to test the services please use endpoints below.

**Add News**
----
* **URL**

  _``http://3.0.147.116:8000/news``_

* **Method**

  `GET`

* **Sample Call**

  ```
  curl --request GET \
  --url 'http://3.0.147.116:8000/news?page=1'
  ```
  
**Get News**
----
* **URL**

  _``http://3.0.147.116:8000/news?page=1``_

* **Method**

  `POST`
  
* **Sample Call**

  ```
  curl --request POST \
  --url http://3.0.147.116:8000/news \
  --header 'content-type: application/json' \
  --data '{\n	"author":"J.R.R. Tolkien",\n	"body":"The Fellowship of the Ring"\n}'
  ```
  
