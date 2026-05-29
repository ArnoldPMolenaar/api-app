# 🔧 API-App

![Go](https://img.shields.io/badge/Go-1.25.9-blue)
![Fiber](https://img.shields.io/badge/Fiber-3.0-green)
![Gorm](https://img.shields.io/badge/Gorm-1.21.12-orange)

## 📜 Description

This API provides endpoints to manage applications and their associated domain names. Each application can have multiple domain names, and each domain name has its own settings.  
  
**Features:**
- **Operations for Apps:** Create, read, update, and delete applications.
- **CRUD Operations for Domains:** Create, read, update, and delete domain names associated with applications.
- **Settings Management:** Retrieve settings for domain names by domain ID or domain name.
- **Event Dispatching:** A event should be dispatched to all other microservices to notify them of the new application.

## 📋 Endpoints
### Private Routes

- **Apps**
    - `GET /v1/apps/` - Get all apps
    - `POST /v1/apps/` - Create a new app
    - `GET /v1/apps/:id` - Get an app by ID
    - `PATCH /v1/apps/:id` - Update an app by ID
    - `DELETE /v1/apps/:id` - Delete an app by ID
    - `POST /v1/apps/:id/restore` - Restore a deleted app by ID
    - `GET /v1/apps/settings` - Get settings by app name
    - `GET /v1/apps/:id/settings` - Get settings by app ID

- **Domains**
    - `POST /v1/domains/` - Create a new domain
    - `GET /v1/domains/:id` - Get a domain by ID
    - `PATCH /v1/domains/:id` - Update a domain by ID
    - `DELETE /v1/domains/:id` - Delete a domain by ID
    - `POST /v1/domains/:id/restore` - Restore a deleted domain by ID
    - `GET /v1/domains/settings` - Get settings by domain name
    - `GET /v1/domains/:id/settings` - Get settings by domain ID

### Public Routes

- **Settings**
    - `GET /v1/settings/apps` - Get settings by app name
    - `GET /v1/settings/apps/:id` - Get settings by app ID
    - `GET /v1/settings/domains` - Get settings by domain name
    - `GET /v1/settings/domains/:id` - Get settings by domain ID

## 🚀 Getting Started

The project can be easily started with Docker by using the `dev` or `prod` environment.

### Development

```sh
docker compose up -d dev
```

### Production

```sh
docker compose up -d prod
```

## 🤝 Contributing
We welcome contributions! Please fork the repository and submit a pull request.

## 📝 License

This project is licensed under the MIT License.

## 📞 Contact

For any questions or support, please contact [arnold.molenaar@webmi.nl](mailto:arnold.molenaar@webmi.nl).
<hr></hr> Made with ❤️ by Arnold Molenaar