# PROMOTIONAL CAMPAIGN SYSTEM
- Author: Luong Quoc Hai

## Overview:
> This is an promotion campaign system implementation that supports basic logic requirements. To enhance the system, we need to implement more advanced logic and handle emergencies (like millions of requests/second, checking schedules for automated campaigns, resuming or extending campaigns, etc.)

## Requirement:
> Design and implement a promotional campaign system for the Trinity app that offers a 30% discount on Silver subscription plans to the first 100 users who register through campaign links (note that vouchers generated for the campaign are only valid for a limited period).


## Structure:
```plaintext
├── Makefile               # Automation commands for building, testing, running, etc.
├── README.md              # Project overview and setup instructions
├── cmd
│   └── server
│       └── main.go        # Entry point of the application
├── config
│   ├── config.go          # Handles app configurations
│   ├── file
│   │   └── config.yaml    # YAML file with config settings (e.g., DB credentials)
│   └── jwt.go             # Manages JWT configuration and generation
├── controllers
│   ├── campaignController.go  # Handles campaign-related API endpoints
│   ├── purchaseController.go  # Manages purchase-related endpoints
│   ├── userController.go      # Handles user-related actions
│   └── voucherController.go   # Manages voucher-related endpoints
├── database
│   ├── mysql
│   │   └── mysql.go       # MySQL database connection and queries
│   └── redis
│       └── redis.go       # Redis connection and caching
├── docker-compose.yml     # Docker services configuration (MySQL, Redis, etc.)
├── docs
│   ├── docs.go            # Setup for Swagger documentation
│   ├── swagger.json       # API specifications in JSON format
│   └── swagger.yaml       # API specifications in YAML format
├── go.mod                 # Go module dependencies
├── go.sum                 # Dependency checksum file
├── middleware
│   └── middleware.go      # Custom middleware for requests (e.g., JWT validation)
├── migrations
│   └── migrate.go         # Database migration scripts
├── models
│   ├── campaign.go        # Campaign entity definition
│   ├── purchase.go        # Purchase entity definition
│   ├── user.go            # User entity definition
│   └── voucher.go         # Voucher entity definition
├── routes
│   ├── routes.go          # Primary file for defining API routes
│   └── routes_backup.go   # Backup routes file (for testing or alternative setups)
├── services
│   ├── campaignService.go # Business logic for campaigns
│   ├── purchaseService.go # Business logic for purchases
│   ├── userService.go     # Business logic for users
│   └── voucherService.go  # Business logic for vouchers
└── utils
    └── utils.go           # Utility functions shared across the app
```

**Folder details**:
- cmd/server: Main application entry point (to initialize routes and services).
- config: Configuration files (for Redis, JWT, DB).
- controllers: API route handlers.
- database: Initialize the database (Redis, MySQL)
- database/redis: Redis-related logic (for caching, session management).
- docs: support for swaggering documentation
- middleware: Middleware logic for JWT - authentication and error handling. Using Redis to store user_id.
- migrations: helper for migrations to database
- models: Database models.
- routes: API route definitions.
- services: Business logic for each service.
- utils: utility functions shared across the application and response formatting.

## Architect diagram:
```plaintext
+---------------------------+
|        Frontend           |
|---------------------------|
| Campaign Landing Page     |
| Subscription Page         |
+------------+--------------+
             |
             v
+---------------------------+
|       API Gateway         |
|---------------------------|
| Routes requests to        |
| appropriate services      |
+------------+--------------+
             |
             v
+---------------------------+
|      Backend Services     |
|---------------------------|
| + Campaign Service       +|  <-- Tracks eligible users and campaign logic
| + Voucher Service        +|  <-- Generates and validates vouchers
| + User Service           +|  <-- Handles user registration
| + Purchase Service       +|  <-- Processes payments and applies discounts
+------------+--------------+
             |
             v
+---------------------------+
|       Database Layer      |
|---------------------------|
| + MySQL                  +|  <-- Persistent storage for users, campaigns, purchases
| + Redis                  +|  <-- Real-time tracking (voucherCode, first 100 users)
+------------+--------------+
             |
             v
+---------------------------+
|     External Systems      |
|---------------------------|
| + Payment Gateway        +|  <-- Handles subscription payment transactions
+---------------------------+

```

## Flow:
![alt text](assets/flow.png)

## Database:
![alt text](assets/Promotion.png)
Read file `database/schema.sql`

## API Documentation:
Please start project and access the link: `http://localhost:8080/swagger/index.html`

## Technical decisions:
- Programming Language: Go (Golang)
- Framework: Gin
- Database: MySQL
- Caching: Redis
- Authentication: JWT
- ORM: GORM

## Assumptions made:
1. The promotional campaign applies only to Silver subscription plans.
2. The discount is fixed at 30% and cannot be changed during the campaign.
3. Only the first 100 users who register through campaign links are eligible for the discount.
4. Vouchers are generated automatically when a user registers through a campaign link.
5. Vouchers have a predefined validity period after issuance.
6. A user can redeem only one voucher per campaign.
7. Campaign status (active, paused, expired) dictates whether vouchers can be issued or redeemed.
8. The system relies on MySQL for persistent storage and Redis for caching frequently accessed data.
9. JWT is used for user authentication and authorization.
10. Gin handles routing and API endpoints for the system.
11. Data integrity and concurrency are handled using GORM transactions.
12. The system includes mechanisms to track the total number of vouchers issued.
13. Campaign creators are pre-registered users with specific roles in the system.

## Future improvements

1. **Campaign Variability**:
   - Implement dynamic discount percentages for different user segments or campaigns.
   - Allow multiple campaign types (e.g., time-based, referral-based) with different reward structures.

2. **Voucher Expiration Management**:
   - Introduce automated expiry notifications for users about soon-to-expire vouchers.
   - Implement a grace period or extension feature for vouchers if required.

3. **Rate Limiting**:
   - Add rate limiting to prevent abuse of campaign links and ensure fair distribution of vouchers.

4. **Logging and Monitoring**:
   - Integrate centralized logging for better traceability of campaign events.
   - Implement real-time monitoring and alerting for campaign performance and system issues.

5. **Campaign Analytics**:
   - Develop a dashboard for campaign creators to view real-time stats on registrations, voucher usage, and conversion rates.
   - Implement advanced reporting tools to track campaign effectiveness.

6. **Database Optimization**:
   - Use read replicas in MySQL for scaling database reads.
   - Implement additional indexing strategies to speed up queries, especially for large campaigns.

7. **Security Enhancements**:
   - Implement **2FA** (Two-factor Authentication) for campaign creators or admins to add an extra layer of security.
   - Enhance **JWT security** by introducing token refresh mechanisms for long-lived sessions.

8. **Multi-Region Support**:
   - Extend the system to support global campaigns with geo-specific rules, currencies, and language preferences.

9. **User Experience Enhancements**:
   - Allow users to view their eligible campaigns and associated discounts within their account profile.
   - Implement email/SMS notifications for users when a campaign voucher is successfully applied or expired.

10. **API Rate Limiting**:
    - Integrate rate limiting for API requests to handle high traffic volumes, ensuring system stability under load.

11. **Automated Testing**:
    - Add end-to-end tests to simulate the entire user registration, voucher generation, and redemption process.

12. **Microservices Transition**:
    - Gradually migrate to a **microservices architecture** for scaling the campaign service independently from other parts of the system.

These improvements will ensure scalability, security, and enhanced user engagement in future versions of the system.

## Local setup guide

### **Local Setup Guide**

This guide will walk you through the steps to set up the promotional campaign system on your local machine.

### **Prerequisites**

- **Go** version 1.18+ installed
- **MySQL** or **MariaDB** version 5.7+ installed and running
- **Redis** installed and running
- **Golang dependencies** installed using `go mod`
  
### **Step 1: Clone the repository**

Clone the repository to your local machine using:

```bash
git clone github.com/luongquochai/promotional-campaign-system
cd promotional-campaign-system
```

### **Step 2: Setup the environment**

#### **Config File Setup**

Ensure you have the correct configuration values. The configuration is loaded from the `config.yaml` file. If the file doesn't exist, create it under the `config/file` directory.

Here’s the `config.yaml` template:

```yaml
dsn: "root:Chipgau164@@tcp(localhost:3306)/promotional_campaign?charset=utf8mb4&parseTime=True&loc=Local"
redis_addr: "localhost:6379"
redis_db: 0
port: ":8080"
```

#### **Database Setup**

1. **Create the Database**:
   If the database doesn't exist, run the following SQL command in MySQL to create it:

   ```sql
   CREATE DATABASE IF NOT EXISTS promotional_campaign;
   ```

2. **Run the `schema.sql` Script**:
   Execute the `schema.sql` script to set up the tables and indexes in your MySQL database.

   ```sql
   USE promotional_campaign;
   -- Execute the full schema.sql script
   ```

3. **Ensure Tables Are Created**:
   Verify that the following tables have been created:

   - `users`
   - `campaigns`
   - `purchases`
   - `vouchers`

   And ensure that the indexes are created as defined in the `schema.sql` file.

#### **Run Database Migration**

The `migration.go` script will handle the migration for you by automatically creating the necessary tables. You can run it with:

```bash
go run migrations/migrate.go
```

This will create the tables for `User`, `Campaign`, `Voucher`, and `Purchase` models.

### **Step 3: Start Redis**

Make sure Redis is running locally on `localhost:6379`. If Redis is not installed, you can install and start it as follows:

#### **For Linux/MacOS**:
```bash
# Install Redis (if not installed)
brew install redis

# Start Redis server
redis-server
```

#### **For Windows**:
- Follow the [official Redis installation guide](https://redis.io/download) for Windows.
- Run `redis-server` to start Redis.

### **Step 4: Run the application**

1. **Install Dependencies**:
   
   If you haven’t already, run the following to install Go dependencies:

   ```bash
   go mod tidy
   ```

2. **Run the Application**:
   
   Start the application with:

   ```bash
   go run cmd/server/main.go
   ```

   The application should now be running at `localhost:8080` (or a different port if specified in your `config.yaml`).

### **Step 5: Testing the Setup**

Access assets to get file API test in POST-Man:
`assets/promotion.postman_collection.json`
Import this file into your Postman and setup evironment variable: `"basicAuth: <token>"`

- **Test Database Connectivity**: Ensure the application can connect to MySQL and Redis properly.
- **Test API Endpoints**: Test the various endpoints to ensure that the promotional campaign logic works as expected.

### **Step 6: Clean Up**

To stop the application, simply press `Ctrl + C` in the terminal.

---

### **Troubleshooting**

- **MySQL Connection Issues**: Ensure that MySQL is running and that the credentials in `config.yaml` are correct.
- **Redis Connection Issues**: Make sure Redis is running on the correct port (`localhost:6379`).

---

### **Future Improvements**

- Automate Redis and MySQL startup using Docker for easier local setup.
- Add more configurable parameters to the `config.yaml` for flexibility in production setups.

## Results:
Below some examples in my projects, you can test more in file 
`assets/promotion.postman_collection.json`
```bash
make run
```
![alt text](assets/main.png)

### /user/register
- payload:
```json
{
    "username": "lqhai4",
    "password": "testing",
    "email": "test3@gmail.com",
    "role": "admin"
}
```
- response:
```json
{
    "message": "User created successfully"
}
```

### /user/login
- payload:
```json
{
    "email": "test3@gmail.com",
    "password":"testing"
}
```
- response:
```json
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzI0MjYxNTQsInVzZXJfaWQiOjl9.NQoS90Fp9Sv8q6RaDVURjsOkWy4s9ULRZbMIpSZjn50"
}
```

### campaigns: POST
**Logic**: Check duplicated campaign if campaign has the same discount price and event date in the existing campaign date -> true
If above condition is true but the campaign has the status != active => the new campaign will be created.

- payload:
```json
{
    "name": "Holiday Promo8",
    "discount_percentage": 44,
    "start_date": "2024-11-01T00:00:00Z",
    "end_date": "2024-12-31T23:59:59Z",
    "max_users": 4
}
```
- response:
```json
{
    "error": "duplicate campaign"
}{
    "data": {
        "id": 15,
        "name": "Holiday Promo8",
        "discount_percentage": 44,
        "start_date": "2024-11-01T07:00:00+07:00",
        "end_date": "2025-01-01T06:59:59+07:00",
        "max_users": 4,
        "status": "active",
        "created_at": "2024-11-20T22:05:35.633+07:00",
        "updated_at": "2024-11-20T22:05:35.633+07:00",
        "creator_id": 6,
        "creator": {
            "id": 0,
            "username": "",
            "email": "",
            "password": "",
            "role": "",
            "created_at": "0001-01-01T00:00:00Z",
            "updated_at": "0001-01-01T00:00:00Z"
        }
    },
    "error": {}
}
```

### /voucher/generate
- payload:
```json
{
    "campaign_id": 15
}
```

```json
{
    "voucher": {
        "id": 15,
        "code": "S2Q9GD65",
        "user_id": 9,
        "campaign_id": 15,
        "discount": 44,
        "valid_from": "2024-11-21T12:35:00.606096346+07:00",
        "valid_to": "2024-12-21T12:35:00.606096346+07:00",
        "created_at": "2024-11-21T12:35:00.606097546+07:00",
        "updated_at": "2024-11-21T12:35:00.606097546+07:00"
    }
}
```
User has only one voucherCode for 1 campaign:
```json
{
    "error": "voucher code {\"id\":15,\"code\":\"S2Q9GD65\",\"user_id\":9,\"campaign_id\":15,\"discount\":44,\"valid_from\":\"2024-11-21T12:35:00.606096346+07:00\",\"valid_to\":\"2024-12-21T12:35:00.606096346+07:00\",\"created_at\":\"2024-11-21T12:35:00.606097546+07:00\",\"updated_at\":\"2024-11-21T12:35:00.606097546+07:00\"} for campaign id 15 of user id 9 already generated"
}
```

### /purchase/create
Simulate the price is `100` and after applying voucher `(-44%)` discount and the final price will be `56`
- payload:
```json
{
    "campaign_id": 15
}
```
-- response:
```json
{
    "data": {
        "id": 0,
        "user_id": 9,
        "transaction_id": "TXN989351",
        "subscription_id": 15,
        "discount_applied": 44,
        "final_price": 56,
        "status": "completed",
        "created_at": "0001-01-01T00:00:00Z",
        "updated_at": "0001-01-01T00:00:00Z"
    },
    "message": "Purchase successful"
}
```