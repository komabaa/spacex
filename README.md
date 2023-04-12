# Siemens Mobility coding challenge

This application provides simple REST API endpoints, takes it's data from [SpaceX API](https://github.com/r-spacex/SpaceX-API/tree/master/docs/v4) and based on users request calculates the requested data and returns it in JSON.


**URL**
Root api URL
  ```http://localhost:8080/api/v1/```

**Method:**
`GET` 

**Run**
```go run main.go```

| folder | description |
| ------- | ----------- |
| /load | Returns sum of all payloads weight that were in space on rocket named Falcon. |
| /crew | Returns number of crew members that were in space. |
