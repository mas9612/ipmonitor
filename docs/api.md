# API

## GET /hosts
Get all hosts information.

### Success
* Status: 200 OK
* Response:
```json
{
  "count": 2,
  "hosts": [
    {
      "id": 1,
      "created_at": "2018-08-29T00:00:00Z",
      "address": "10.1.240.151",
      "hostname": "k8s-01",
      "description": "k8s node #1"
    },
    {
      "id": 2,
      "created_at": "2018-08-29T00:00:00Z",
      "address": "10.1.240.152",
      "hostname": "k8s-02",
      "description": "k8s node #2"
    }
  ]
}
```

## POST /hosts
Add new host information.

### Request Body
| field       | description | Required |
|:------------|:------------|:---------|
| address     | IP address  | true     |
| hostname    | Hostname    | true     |
| description | Description | false    |

```json
{
  "address": "10.1.240.151",
  "hostname": "k8s-01",
  "description": "k8s node #1"
}
```

### Success
* Status: 201 Created
* Response:
```json
{
  "id": 1,
  "created_at": "2018-08-29T00:00:00Z",
  "address": "10.1.240.151",
  "hostname": "k8s-01",
  "description": "k8s node #1"
},
```

### Error
* Status: 400 Bad Request (If `address` or `hostname` was not included in request body)
* Response:
```json
{
  "status": 400,
  "message": "Key \"address\" and \"hostname\" are requied"
}
```

## GET /hosts/{id}
Get specific host information.

### Success
* Status: 200 OK
* Response:
```json
{
  "id": 1,
  "created_at": "2018-08-29T00:00:00Z",
  "address": "10.1.240.151",
  "hostname": "k8s-01",
  "description": "k8s node #1"
}
```

### Error
1. If `id` is not integer
    - Status: 400 Bad Request
    - Response:
    ```json
    {
      "status": 400,
      "message": "Key \"id\" must be integer"
    }
    ```
1. If record not found
    - Status: 404 Not Found
    - Response:
    ```json
    {
      "status": 404,
      "message": "ID \"999\" not found"
    }
    ```

## PUT /hosts/{id}
Update (or create if specified `id` not found) host record.

### Request Body
| field       | description | Required |
|:------------|:------------|:---------|
| address     | IP address  | true     |
| hostname    | Hostname    | true     |
| description | Description | false    |

```json
{
  "address": "10.1.240.151",
  "hostname": "k8s-01",
  "description": "k8s node #1"
}
```

### Success
* Status: 200 OK
* Response:
```json
{
  "id": 1,
  "created_at": "2018-08-29T00:00:00Z",
  "address": "10.1.240.151",
  "hostname": "k8s-01",
  "description": "k8s node #1"
}
```

### Error
* Status: 400 Bad Request (If `address` or `hostname` was not included in request body)
* Response:
```json
{
  "status": 400,
  "message": "Key \"address\" and \"hostname\" are requied"
}
```

## DELETE /hosts/{id}
Delete host record.

### Success
* Status: 204 No Content
* Response:
```json
{
  "id": 1,
  "created_at": "2018-08-29T00:00:00Z",
  "address": "10.1.240.151",
  "hostname": "k8s-01",
  "description": "k8s node #1"
}
```

### Error
1. If `id` is not integer
    - Status: 400 Bad Request
    - Response:
    ```json
    {
      "status": 400,
      "message": "Key \"id\" must be integer"
    }
    ```
1. If record not found
    - Status: 404 Not Found
    - Response:
    ```json
    {
      "status": 404,
      "message": "ID \"999\" not found"
    }
    ```
