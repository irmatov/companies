# Requirements

In this exercise you need to build a REST API microservice to handle Companies. Company is an entity defined by the following attributes:

* Name
* Code
* Country
* Website
* Phone

All four CRUD operations are required. For read operation, fetching one or many companies should be available. Each Company’s attribute should be available as filtering in the CRUD operations.

Creation and deletion operations must be allowed only for requests received from users located in Cyprus. The location must be retrieved based on the user’s IP address via the service
https://ipapi.co/.

## Optional

On each mutating operation, a JSON formatted event must be produced to a service bus (Kafka, RabbitMQ etc.). In a production environment, those events could be used to notify other microservices of a data change and conduct some business logic, i.e. sending emails.
