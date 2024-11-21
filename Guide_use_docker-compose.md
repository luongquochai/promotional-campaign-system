## Local setup guide
### **Build and Start the Services**

Now you can build and start all the services with Docker Compose.
Run the following command from your project directory:

```bash
docker-compose up --build
```

This will:
1. Build your Go application.
2. Set up MySQL, Redis, and your application in Docker containers.
3. Start the services.

Once everything is up and running, you can access your application at `http://localhost:8080` (or the port specified in `config.yaml`).

- **Access MySQL**: You can access the MySQL service by running:

  ```bash
  docker exec -it mysql_promotional_campaign mysql -u root -p
  ```

  Then, use the password `Chipgau164@` to log into MySQL.

- **Access Redis**: You can connect to Redis by running:

  ```bash
  docker exec -it redis_promotional_campaign redis-cli
  ```

### **Clean Up**

To stop and remove the containers, network, and volumes created by Docker Compose:

```bash
docker-compose down --volumes
```

> We need to run Migrate database before testing:
### **Database Setup**
- Create the Database: If the database doesn't exist, run the following SQL command in MySQL to create it:

```sql
Copy code
CREATE DATABASE IF NOT EXISTS promotional_campaign;
```

### **Automatic Table Creation with GORM (On Application Start)**

The easiest approach is to have GORM automatically create or migrate your tables when your application starts. This can be done using the `AutoMigrate` method.

**Example**: In your `migration.go` file, you are already using `AutoMigrate` to create tables.

Here’s what happens:
- **AutoMigrate** will automatically create the tables defined in your models (`User`, `Campaign`, `Voucher`, `Purchase`) if they don’t exist.
- It will also update the schema to match your model definitions if there are any changes (e.g., new columns or modified types).

#### **How to run this migration:**

- **Manual Migration Command**

1. Start your containers:

   ```bash
   docker-compose up -d
   ```

2. Access the application container:

   ```bash
   docker exec -it promotional_campaign_app sh
   ```

3. Run the migration manually:

   ```bash
   go run cmd/migration/main.go
   ```

This will run the migration code and create the necessary tables in MySQL.

### **2. Using `schema.sql` for Manual Migration**

If you prefer to handle migrations manually (for example, to keep better control over the schema), you can execute the `schema.sql` file to create your database schema.

You can copy your `schema.sql` file to the Docker container and run it within MySQL.

1. **Copy the `schema.sql` to the MySQL Container:**

   In your `docker-compose.yml`, you can mount the `schema.sql` file into the MySQL container:

   ```yaml
   volumes:
     - ./schema.sql:/docker-entrypoint-initdb.d/schema.sql
   ```

2. **Rebuild and restart your containers:**

   ```bash
   docker-compose up --build
   ```

3. **MySQL Initialization:**

   The MySQL container will automatically execute any `.sql` files placed in the `/docker-entrypoint-initdb.d/` directory when it starts. So, this will automatically create your database and tables during the MySQL container's initialization.

4. **Check Database Tables:**

   You can check if the tables are created by logging into the MySQL container and inspecting the schema:

   ```bash
   docker exec -it mysql_promotional_campaign mysql -u root -p
   ```

   Then, run the following query to see the tables:

   ```sql
   SHOW TABLES;
   ```

You can combine both methods:
- Use `make migrate` for automatic table creation during development.
- Use the `schema.sql` for manual migrations or when deploying to production.

---
