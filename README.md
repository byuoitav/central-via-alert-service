# central-via-alert-service

This service is configured to find all of the VIAs on campus and deploy an emergency message to all those devices. 

endpoints are configured as follows: 
GET -
https://(ServerName):(Port#)/api/v1/emessage/timer/:timing/building/:BuildingName - Returns all the buildings from the database

POST - 
https://(ServerName):(Port#)/api/v1/emessage/timer/:timing/via/:nameofvia - post to a single via
https://(ServerName):(Port#)/api/v1/emessage/timer/:timing/test - post to a test group
https://(ServerName):(Port#)/api/v1/emessage/timer/:timing/building/:BuildingName - post to a building
https://(ServerName):(Port#)/api/v1/emessage/timer/:timing/all - post to campus

