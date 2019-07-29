# kmpr-test
kumparan technical assessment

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

  ``_http://3.0.147.116:8000/news?page=1_``

* **Method**

  `POST`
  
* **Sample Call**

  ```
  curl --request POST \
  --url http://3.0.147.116:8000/news \
  --header 'content-type: application/json' \
  --data '{\n	"author":"J.R.R. Tolkien",\n	"body":"The Fellowship of the Ring"\n}'
  ```
  
