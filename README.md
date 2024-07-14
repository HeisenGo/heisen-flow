# 📝 HeisenFlow Task Management System

Welcome to the Task Management System! This project is built with Golang and Fiber, designed to help you manage tasks efficiently. Let's dive into the details!

## 🚀 Project Summary

### **User Registration and Authentication**
- **Sign Up**: Users can register by providing their name, email, and password.
- **Log In**: Post-registration, users log in with their credentials, undergo authentication, and gain access to protected resources. Access control ensures user data security.

### **Notification Inbox**
- **Events**: User-related events like invitations to new boards, task status changes, and new notifications are timestamped and maintained.
- **Status**: Messages can be marked as read or unread, and users can change their status.

### **Board Management**
- **Create Boards**: Users can create new boards, which function as independent projects containing tasks and subtasks.
- **Visibility**: Boards can be private (visible only to invited members) or public (visible to all users).
- **Roles**: Each user on a board can have one of four roles: viewer, editor, maintainer, or owner, with varying levels of access and permissions.

### **Task and Subtask Management**
- **Role-Based Activities**: Tasks can have an unlimited number of subtasks in a hierarchical structure.
- **Columns**: Boards must have various columns, including a "done" column for completed tasks. Columns and tasks can be reordered for optimal workflow management.
- **Task Details**: Users can create new tasks with details such as name, description, start/end dates, assignees, and story points.
- **Comments and Notifications**: Task comments are role-based, and status changes trigger notifications to relevant users.
- **Dependencies**: Task dependencies are tree-structured without loops, and completed tasks become independent.
- **Progress Tracking**: Users can track subtask progress, and tasks can be imported/exported in various formats.
- **Custom Fields**: Users can define custom fields for tasks.

### **Technical and Implementation Requirements**
- **Git Workflow**: Use appropriate git workflow for teamwork.
- **Frameworks and Packages**: Utilize preferred frameworks and packages.
- **Docker**: Dockerize the project with at least two containers (application and database) connected via Docker network.
- **Modular Structure**: Implement a modular structure.
- **Database**: Store data in a relational database like MySQL.
- **Logging**: Log transactions in dedicated log files, e.g., `transaction.log`.
- **Security**: Hash sensitive data like passwords and use UUIDs for unique identifiers.
- **Custom Exceptions**: Define and manage custom exceptions.
- **Unit Tests**: Implement unittests for all project components.
- **Caching**: Cache data to improve system performance.
- **API Documentation**: Fully document APIs and modules, ensuring all backend routes have Swagger UI for request/response testing.
- **Secure Methods**: Use secure methods for user data and transaction storage to prevent attacks.
- **Notification System**: Implement a notification system for task status updates.
- **Creative Enhancements**: Additional features and creative enhancements are encouraged if the primary project scope is met.

## 🛠️ How to Run the Project

### Method 1: Using Docker
**Prerequisites**: Docker

1. Rename the `config.yaml.example` file to `config.yaml`.
2. Run the `make dockerize` command.
3. If you haven't changed the variables in the `config.yaml.example` file, you will see the main page of the application by clicking the link http://0.0.0.0:8080 in your browser.

### Method 2: Using Pre-installed PostgreSQL, Redis, and Go
**Prerequisites**: Go, PostgreSQL, Redis

1. Change the file `config.yaml.local-example` to `config.yaml` and fill its values according to the settings of your local databases.
2. Run the `make run` command.
3. Now you will have the project running according to the `config.yaml` values.

### Running Tests
- Use the `make test` command to run the tests.

---

Happy coding! 🐛 Be bug-free and enjoy managing your tasks! 🎉
