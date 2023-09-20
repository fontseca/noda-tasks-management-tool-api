<div align="center">
  <img src="./assets/noda_logo.svg" alt="drawing" style="width:400px;"/>

---

  <a href="https://golang.org/doc/go1.20">
    <img src="https://img.shields.io/badge/Go-1.20-blue.svg" />
  </a>
  <a href="https://opensource.org/licenses/MIT">
    <img src="https://img.shields.io/badge/license-MIT-brightgreen.svg" />
  </a>
  <img src="https://img.shields.io/github/last-commit/fontseca/noda?color=61dfc6&label=last%20commit" />
</div>

**NODA** is a task management RESTful API designed to simplify the process of managing tasks, lists, and user interactions. It is uses PostgreSQL for data storage, JWT authentication for security, and Docker for easy deployment.

## Table of Contents

- [Table of Contents](#table-of-contents)
- [API Endpoints](#api-endpoints)
  - [Authentication](#authentication)
  - [Users](#users)
  - [Groups](#groups)
  - [Lists](#lists)
  - [Tasks](#tasks)
  - [Steps](#steps)
  - [Tags](#tags)
  - [Attachments](#attachments)

## API Endpoints

### Authentication

| HTTP Method | Endpoint              | Description                               | Role  | Notes |
| :---------: | --------------------- | ----------------------------------------- | :---: | ----- |
|  **POST**   | `/signup`             | Create a new user                         |  any  | —     |
|  **POST**   | `/signin`             | Log in an existent user                   | user  | —     |
|  **POST**   | `/me/logout`          | Log out the current user                  | user  | —     |
|  **POST**   | `/me/change_password` | Change the password of the logged in user | user  | —     |

### Users

| HTTP Verb  | Endpoint                 | Description                                          | Role  | Notes |
| :--------: | ------------------------ | ---------------------------------------------------- | :---: | ----- |
|  **GET**   | `/users`                 | Retrieve all users                                   | admin | —     |
|  **GET**   | `/users/search`          | Search for users                                     | admin | —     |
|  **GET**   | `/users/{user_id}`       | Retrieve a user                                      | admin | —     |
| **DELETE** | `/users/{user_id}`       | Permanently remove a user and all its related data   | admin | —     |
|  **PUT**   | `/users/{user_id}/block` | Block one user                                       | admin | —     |
| **DELETE** | `/users/{user_id}/block` | Unblock one user                                     | admin | —     |
|  **GET**   | `/users/blocked`         | Retrieve all blocked users                           | admin | —     |
|  **GET**   | `/me`                    | Get the logged in user                               | user  | —     |
|  **PUT**   | `/me`                    | Partially update the account of the logged in user   | user  | —     |
| **DELETE** | `/me`                    | Permanently remove the account of the logged in user | user  | —     |
|  **GET**   | `/me/settings`           | Retrieve all the settings of the logged in user      | user  | —     |

### Groups

| HTTP Method | Endpoint                 | Description                                    | Role  | Notes |
| :---------: | ------------------------ | ---------------------------------------------- | :---: | ----- |
|   **GET**   | `/me/groups`             | Retrieve all the groups                        | user  | —     |
|  **POST**   | `/me/groups`             | Create a new group                             | user  | —     |
|   **GET**   | `/me/groups/{groups_id}` | Retrieve a group                               | user  | —     |
|  **PATCH**  | `/me/groups/{groups_id}` | Partially update a list                        | user  | —     |
| **DELETE**  | `/me/groups/{groups_id}` | Permanently remove a list and all related data | user  | —     |
|   **GET**   | `/me/groups/{groups_id}` | Retrieve a list                                | user  | —     |

### Lists

| HTTP Method | Endpoint                                | Description                                               | Role  | Notes                                 |
| :---------: | --------------------------------------- | --------------------------------------------------------- | :---: | ------------------------------------- |
|   **GET**   | `/me/lists`                             | Retrieve all the ungrouped lists                          | user  | —                                     |
|  **POST**   | `/me/lists`                             | Create a new ungrouped list                               | user  | —                                     |
|   **GET**   | `/me/lists/{list_id}`                   | Retrieve a lungrouped ist                                 | user  | —                                     |
|  **PATCH**  | `/me/lists/{list_id}`                   | Partially update a ungrouped list                         | user  | —                                     |
| **DELETE**  | `/me/lists/{list_id}`                   | Permanently remove an ungrouped list and all related data | user  | Can't remove Today and Tomorrow lists |
|   **GET**   | `/me/groups/{group_id}/lists`           | Retrieve all the lists of a group                         | user  | —                                     |
|  **POST**   | `/me/groups/{group_id}/lists`           | Create a new list for a group                             | user  | —                                     |
|   **GET**   | `/me/groups/{group_id}/lists/{list_id}` | Retrieve a list of a group                                | user  | —                                     |
|  **PATCH**  | `/me/groups/{group_id}/lists/{list_id}` | Partially update a list of a group                        | user  | —                                     |
| **DELETE**  | `/me/groups/{group_id}/lists/{list_id}` | Permanently remove a list of a group and all related data | user  | —                                     |

### Tasks

| HTTP Method | Endpoint                                      | Description                                        | Role  | Notes                                             |
| :---------: | --------------------------------------------- | -------------------------------------------------- | :---: | ------------------------------------------------- |
|   **GET**   | `/me/today`                                   | Retrieve all the tasks from the Today list         | user  | —                                                 |
|   **GET**   | `/me/tomorrow`                                | Retrieve all the tasks for tomorrow                | user  | —                                                 |
|   **GET**   | `/me/tasks`                                   | Retrieve all the tasks                             | user  | —                                                 |
|  **POST**   | `/me/tasks`                                   | Create a new task and store in the Today list      | user  | —                                                 |
|   **GET**   | `/me/tasks/search`                            | Search for tasks                                   | user  | —                                                 |
|   **GET**   | `/me/tasks/completed`                         | Retrieve all the completed tasks                   | user  | —                                                 |
|   **GET**   | `/me/tasks/archived`                          | Retrieve archived tasks                            | user  | —                                                 |
|   **GET**   | `/me/tasks/trashed`                           | Retrieve trashed tasks                             | user  | —                                                 |
|   **GET**   | `/me/tasks/{task_id}`                         | Retrieve a task                                    | user  | Might come from any list                          |
|  **PATCH**  | `/me/tasks/{task_id}`                         | Partially update a task                            | user  | —                                                 |
| **DELETE**  | `/me/tasks/{task_id}`                         | Permanently remove a task and all related data     | user  | —                                                 |
|   **PUT**   | `/me/tasks/{task_id}/trash`                   | Move a task to trash                               | user  | Will get destroyed 30 days later                  |
| **DELETE**  | `/me/tasks/{task_id}/trash`                   | Recover a task from trash                          | user  | —                                                 |
|   **PUT**   | `/me/tasks/{task_id}/reorder`                 | Rearrange a task in its list                       | user  | Use the `dst` query parameter as the new position |
|   **GET**   | `/me/lists/{list_id}/tasks`                   | Retrieve all the tasks of an ungrouped list        | user  | —                                                 |
|  **POST**   | `/me/lists/{list_id}/tasks`                   | Create a task and save it in an ungrouped list     | user  | —                                                 |
|   **GET**   | `/me/groups/{group_id}/lists/{list_id}/tasks` | Retrieve all the tasks of a list in a group        | user  | —                                                 |
|  **POST**   | `/me/groups/{group_id}/lists/{list_id}/tasks` | Create a task and save it in a list within a group | user  | —                                                 |

### Steps

| HTTP Method | Endpoint                                         | Description                           | Role  | Notes                                             |
| :---------: | ------------------------------------------------ | ------------------------------------- | :---: | ------------------------------------------------- |
|   **GET**   | `/me/tasks/{task_id}/steps`                      | Retrieve the steps to achieve a task  | user  | —                                                 |
|  **POST**   | `/me/tasks/{task_id}/steps`                      | Add a new step to achive a task       | user  | —                                                 |
|  **PATCH**  | `/me/tasks/{task_id}/steps/{step_id}`            | Partially update a step               | user  | —                                                 |
|   **PUT**   | `/me/tasks/{task_id}/steps/{step_id}/accomplish` | Mark a step as acomplished            | user  | —                                                 |
| **DELETE**  | `/me/tasks/{task_id}/steps/{step_id}/accomplish` | Unmark a step as acomplished          | user  | —                                                 |
| **DELETE**  | `/me/tasks/{task_id}/steps/{step_id}`            | Permanently remove a step from a task | user  | —                                                 |
|  **POST**   | `/me/tasks/{task_id}/steps/{step_id}/reorder`    | Rearrange a step in a task            | user  | Use the `dst` query parameter as the new position |

### Tags

| HTTP Method | Endpoint            | Description                          | Role  | Notes |
| :---------: | ------------------- | ------------------------------------ | :---: | ----- |
|   **GET**   | `/me/tags`          | Retrieve the steps to achieve a task | user  | —     |
|  **POST**   | `/me/tags`          | Create a new tag                     | user  | —     |
|  **PATCH**  | `/me/tags/{tag_id}` | Partially update a tag               | user  | —     |
| **DELETE**  | `/me/tags/{tag_id}` | Permanently remove a tag             | user  | —     |

### Attachments

| HTTP Method | Endpoint                                          | Description                                     | Role  | Notes |
| :---------: | ------------------------------------------------- | ----------------------------------------------- | :---: | ----- |
|   **GET**   | `/me/tasks/{task_id}/attachments`                 | Retrieve all the attachments in a task (if any) | user  | —     |
|   **GET**   | `/me/tasks/{task_id}/attachments/{attachment_id}` | Get an attachments in a task                    | user  | —     |
| **DELETE**  | `/me/tasks/{task_id}/attachments/{attachment_id}` | Permanently remove an attachments in a task     | user  | —     |
