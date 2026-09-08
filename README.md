# Caloire Snap Notification Service

- store user notifciation profile
- when send noti, look up notification profile from redis,if cache miss, look up db
- if db goes down, programmatically pause queue consumer with circuit breaker pattern instead of routing to DLQ automatically