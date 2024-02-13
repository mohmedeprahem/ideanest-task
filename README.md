# Golang API Application

## Description

This repository houses a Golang-based API application designed for managing organizations. The application includes features such as token management, CRUD operations for organizations, user invitations, and integration with MongoDB using Docker.

## DataBase Design

Database Design

<div style="text-align:center;">
  <img src="https://cdn.discordapp.com/attachments/411190888212070400/1206783173011316766/image.png?ex=65dd4369&is=65cace69&hm=93c35f56cd8c4d3a45bbfbec47837e4602b2fe29feee63e1dec0de442968590b&" alt="Database design" />
</div>

## Installation

### Clone the Repository:

```bash
# Clone repo
$ git clone https://github.com/mohmedeprahem/ideanest-task.git
```

## Configuration

Before running the application, you need to create two configuration files: `app-config.yaml` and `database-config.yaml`. Follow the instructions below to create these files:

### app-config.yaml

```yaml
jwt:
  atSecret: Your access secret key
  rtSecret: Your refresh secret key
redis:
  address: redis:6379
  password: your redis password if you want
  db: 0
```

### database-config.yaml

```yaml
mongodb:
  uri: database uri
```

## License

[MIT licensed](LICENSE).
