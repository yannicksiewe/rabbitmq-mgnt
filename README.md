# rabbitmq-mgnt
Managing frequent operations with rabbitmq such as updating a user, checking the connection, creating and modifying queues 


# Usage 

```shell
docker run --rm -e ADMINUSER="guest" -e ADMINPASSWORD="guest" -e LOGINACCOUNT="test" -e LOGINPASSWORD="test" -e MQENDPOINT="rabbitmq" rabbitmq-mgnt:latest 
```
