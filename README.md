# go-mongo
## rest-api with go, using fiber and mongodb
### rest-api for digital event/webinar organizing
### try it out on https://guarded-stream-71687.herokuapp.com
### add endpoint suffix beside on route that you want to access
###
#### Endpoints
####  
#### Post => "/api/users/register/" for UserRegistration
#### Post => "/api/users/login/email/" for LoginWithEmail
#### Delete => "/api/users/:id" for DeleteUser
#### Get => "/api/users/" for GetAllUser
#### Post => "/api/users/participant/:eventId/:userId" for JoinEvent
#### Get => "/api/users/:id" for GetUserDetail
####
#### Get => "/api/events/" for GetAllEvents
#### Get => "/api/events/:id" for GetEvent
#### Post => "/api/events/" for AddEvent
#### Put => "/api/events/:id" for UpdateEvent
#### Delete => "/api/events/:id" for DeleteEvent