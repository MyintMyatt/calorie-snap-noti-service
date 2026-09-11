# Caloire Snap Notification Service

- store user notifciation profile
- when send noti, look up notification profile from redis,if cache miss, look up db
- if db goes down, programmatically pause queue consumer with circuit breaker pattern instead of routing to DLQ automatically


## Notification request
```json
{
  "id": "notif_987654321",
  "recipient": "recipient@exmaple.com",
  "priority": 4,
  "client_time": "2026-09-11T15:12:00Z",
  "template": "8342456",
  "locale": "mm",
  "subject" : "Welcome to Calorie Snap",
  "data": {
    "username": "Aung Aung",
    "bonus_amount": 5000,
    "is_new_user": true
  }
}
```