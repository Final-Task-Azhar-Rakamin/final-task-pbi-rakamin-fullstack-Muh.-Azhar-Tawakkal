# API-Final Task

These are all the APIs that will be used in the Final Task project. The API is made to run these functions :
- Handle register and login for user
- Handle update user data
- Handle user posts

## **1. List API**


### **1.1 User API**
#### Register API
This API used to create users account
```
URL: users/register
Method: POST
Request Body (Form-data):
{
    username,
    email,
    password,
    profile_picture,
}
Success Response :
{
  "message": "Register Success",
  "data": {{user data}}
}
```

#### Login API
This API used to validate user login attempt and put token in the cookie
```
URL: users/login
Method: GET
Request Body (JSON):
{
    email,
    password
}
Success Response:
{
    "data": {
        "email": {{email}} ,
        "username": {{username}}
    },
    "message": "login success",
    "token": {{Token authorization}}
}
```

#### Update User API (Login Required)
This API used to update user profile in the database
```
URL: users/:id
Method: PUT
Request Body (Form-data):
{
    username,
    password,
    profile_picture
}
Success Response:
{
    "message": "Update user success"
}
```

#### Delete Users API (Login Required)
This API used to delete user account
```
URL: users/:id
Method: DELETE
Success Response:
{
    "message": "Delete data success"
}
```

#### Logout User API (Login Required)
This API used to logout and clear the cookie
```
URL: users/logout
Method: GET
Response:
{
    "message": "User logged out! Good bye!"
}
```
### **1.2 Photos API**

#### List Posts API (Login Required)
This API used to list the photo that the user has post
```
URL: /photos
Method: GET
Response:
{
    "photos": [{{list of user posts}}],
}
```

#### Get Post Detail API (Login Required)
This API used to see user's photo detail
```
URL: /photos/:id
Method: GET
Response:
{
    "photos": {{post detail}},
}
```

#### Add Post API (Login Required)
This API used to add a photo into the database
```
URL: /photos/:id
Method: POST
Request Body (Form-data):
{
    title,
    caption,
    photo_file
}
Response:
{
    "message": "Input data success",
    "photo Id": "{{Photo Id}}"
}
```

#### Update Post API (Login Required)
This API used to update a photo into the database
```
URL: /photos/:id
Method: PUT
Request Body (JSON):
{
    title,
    caption,
}
Response:
{
    "message": "Update data success"
}
```

#### Delete Post API (Login Required)
This API used to delete a photo from the database
```
URL: /photos/:id
Method: DELETE
Response:
{
    "message": "Delete data success"
}
```
