## Spec

_api_address_: `http://65.108.147.133:8085`

### 1. Auth

_web3_address_ -- address of the participant(user) from his extension(Metamask, UniPass ...).

1. _GET_ `api_address/auth/{web3_address}`

   **Sample responce:**

```
{
	jwt_token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2ODQ1MzI4NjQsImlhdCI6MTY4NDQ0NjQ2NCwiaXNzIjoic293X2xpYnJhcnkiLCJyb2xlIjoiIiwic3RhdHVzIjoiIiwic3ViIjoiN2E0OWNlNWYtNzVmZC00NmYzLTk3ZTItMGViZjMxODg0NDk0In0.HZzOFYySxxkxYpGIFMbxw4eJuPe3eVeFBDtZV5GJ3nY"
}
```

### 2. Publishing a paper

_POST_ `api_address/publish_work`

- REQUIRES HEADER:

  - `Authorization` : `Bearer {jwt_token}`

  **Sample request:**

  ```
  {
  	"author_address": {web3_address},
  	"work":{
  		"name": "My first work",
  		"description": "My first description",
  		"annotation": "My first annotation",
  		"content": {
  			"work_data": "THE BEST WORK EVER!"
  		}
  	}
  }
  ```

  **Sample response:**

  ```
  {
  	"status": "WORK_UNDER_PRE_REVIEW"
  }
  ```

### 3. Get papers by author {web3_address}

- _GET_ `api_address/works/{web3_address}`

  - REQUIRES HEADER:
    - `Authorization` : `Bearer {jwt_token}`

  _The work is under admin review:_
  **Sample response:**

  ```
  {
  	"author_address": {web3_address},
  	"status": "WORK_UNDER_PRE_REVIEW"
  	"work":{}
  }
  ```

  _The work is under validators review._
  **Sample response:**

  ```
  {
  	"author_address": {web3_address},
  	"status": "WORK_UNDER_REVIEW"
  	"work": {
  		"name": "My first work",
  		"description": "My first description",
  		"annotation": "My first annotation",
  		"content": {}
  	}
  }
  ```

  _The work was successfully reviewed._
  **Sample response:**

  ```
  {
  	"author_address": {web3_address},
  	"status": "OPEN"
  	"work": {
  		"name": "My first work",
  		"description": "My first description",
  		"annotation": "My first annotation",
  		"content": {
  			"work_data": "THE BEST WORK EVER!"
  		}
  	}
  }
  ```

  _The work was declined by validators._
  **Sample response:**

  ```
  {
  	"author_address": {web3_address},
  	"status": "DECLINED"
  	"work": {}
  }
  ```
