# Sea of Wisdom :: LIBRARY

## Spec

_api_address:_ `http://65.108.147.133:8005`
_Roles:_ 1. guest: 0 2. reader: 1 3. author: 2 4. validator: 4 5. admin: 5

_(Scientific) Paper = Work_

### 1. Registration / Auth / Basic information

_web3_address_ -- address of the participant(user) from his wallet extension(Metamask, UniPass ...).

#### 1.1. Registration

_POST_ `api_address/new_participant`

**Sample request:**

```
{
	"nickname": "mellaught",
	"web3_address": {web3_address}
}
```

**Sample response:**
_Status code 200_

```
{
    "nickname": "mellaught",
	"role": 1,
    "jwt_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2ODYwMzQ3NjgsImlhdCI6MTY4NTk0ODM2OCwiaXNzIjoic293X2xpYnJhcnkiLCJsYW5ndWFnZSI6IiIsInJvbGUiOjQsInN0YXR1cyI6IiIsInN1YiI6ImRlNWU1MThmLTM4MTUtNGYxNi04ZTExLTQ0MWE4NjU2NjU0ZiIsIndlYjNfYWRkcmVzcyI6IjB4MTAwZGQ2YzI3NDU0Y2IxREFkZDEzOTEyMTRBMzQ0QzYyMDhBOEM4MCJ9.wXbDfDv_rY4ZQKjun9qAbsTASIYWR89qPV4-9Xc6_2A",
}
```

_In case of an error._

```
{
	"error": "The error message"
}
```

#### 1.2 Auth

- _GET_ `api_address/auth/{web3_address}`

**Sample responce:**

```
{
    "jwt_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2ODYwMzQ3NjgsImlhdCI6MTY4NTk0ODM2OCwiaXNzIjoic293X2xpYnJhcnkiLCJsYW5ndWFnZSI6IiIsInJvbGUiOjQsInN0YXR1cyI6IiIsInN1YiI6ImRlNWU1MThmLTM4MTUtNGYxNi04ZTExLTQ0MWE4NjU2NjU0ZiIsIndlYjNfYWRkcmVzcyI6IjB4MTAwZGQ2YzI3NDU0Y2IxREFkZDEzOTEyMTRBMzQ0QzYyMDhBOEM4MCJ9.wXbDfDv_rY4ZQKjun9qAbsTASIYWR89qPV4-9Xc6_2A",
    "role": 4
}
```

#### 1.3. Basic information

_GET_ `api_address/get_basic_info`

- REQUIRES HEADER: - `Authorization` : `Bearer {jwt_token}`

**Sample responce:**

```
{
    "nickname": "mellaught",
	"role": 1,
	"language": "russian"
}
```

## 2. Reader

This section consists of APIs for filling up the **Reader** page.
Let's us to kick off by describing the basic work structure:

**WORK_STRUCTURE:**

```
{
	"name": "My wdwdadw12312a work",
	"description": " descr232132iption",
	"annotation": "!!!wrawdad!!!s",
	"content": {
		"work_data": "12312312312312 ПЛАВАНИКЕ ЭТО КРУЦТО XX"
	},
	"science": "подводное плавание",
	"tags": [
		"спорт",
		"плавание"
	]
}
```

### 2.1. Get all available works(papers)

_GET_ `api_address/works`

- REQUIRES HEADER: - `Authorization` : `Bearer {jwt_token}`

_Status code 200_
**Sample response:**

```
[
	WORK_STRUCTURE,
	WORK_STRUCTURE,
	...
	WORK_STRUCTURE
]
```

_In case of an error._

```
{
	"error": "The error message"
}
```

### 2.2. Get the works by the specific author

_GET_ `api_address/works/{web3_address}`

- REQUIRES HEADER: - `Authorization` : `Bearer {jwt_token}`

_Status code 200_
**Sample response:**

```
[
	WORK_STRUCTURE,
	WORK_STRUCTURE,
	...
	WORK_STRUCTURE
]
```

_In case of an error._

```
{
	"error": "The error message"
}
```

### 2.3. Get the specific work by its ID

UNDER CONSTRUCTION

### 2.4. Bookmarks

#### 2.4.1. Add the paper as bookmark

_POST_ `api_address/add_bookmark/{work_id}`

- REQUIRES HEADER: - `Authorization` : `Bearer {jwt_token}`

**Sample response:**
_Status code 200_

```
{
	"OK"
}
```

_In case of an error._

```
{
	"error": "The error message"
}
```

#### 2.2. Get the participant's bookmarks

_GET_ `api_address/bookmarks`

- REQUIRES HEADER: - `Authorization` : `Bearer {jwt_token}`

**Sample response:**

```
[
	WORK_STRUCTURE,
	WORK_STRUCTURE,
	...
	WORK_STRUCTURE
]
```

#### 2.3. Remove the paper from bookmarks

- _POST_ `api_address/remove_bookmark/{work_id}`
  - REQUIRES HEADER:
    - `Authorization` : `Bearer {jwt_token}`

_The work is under admin review:_
**Sample response:**

```
{
	"OK"
}
```

_In case of an error._
**Sample response:**

```
{
	"error": "The error message"
}
```

#### Edit the reader profile

Under construction.

### 2.5. Purchase of Works

The sections describes the methods related to buying scientific papers(works) by SOW tokens.
The participant cannot purchase a work in case of insufficient tokens in their account.

#### 2.5.1 Buy a work

_GET_ `api_address/purchase_work/{work_id}`

- REQUIRES HEADER:
  - `Authorization` : `Bearer {jwt_token}`

_The work is under admin review:_
**Sample response:**

```
{
	"OK"
}
```

#### 2.5.2. Get works of a particular participant

_GET_ `api_address/purchased_works` - REQUIRES HEADER: - `Authorization` : `Bearer {jwt_token}`
**Under construction.**
