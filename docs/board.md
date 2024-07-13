## Board

The `Board Entity` struct represents a board within the database. Below is a detailed breakdown of its fields and relationships:

```go
type Board struct {
    ID        uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
    Name      string         `gorm:"index"`
    Type      string
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"`
    // Relationships
    Users          []User          `gorm:"many2many:user_board_roles;constraint:OnDelete:CASCADE;"`
    Tasks          []Task          `gorm:"foreignKey:BoardID;constraint:OnDelete:CASCADE;"`
    UserBoardRoles []UserBoardRole `gorm:"foreignKey:BoardID;constraint:OnDelete:CASCADE"`
    Columns        []Column        `gorm:"foreignKey:BoardID;constraint:OnDelete:CASCADE"`
}
```

## Fields
- **ID**: A unique identifier for the board, generated automatically using UUID.

$\quad$ $\quad$**gorm**:"type:uuid;default:uuid_generate_v4();primaryKey"
- **Name**: The name of the board.

$\quad$ $\quad$**gorm**:"index"
- **Type**: The type of the board.

- **CreatedAt**: Timestamp indicating when the board was created.

- **UpdatedAt**: Timestamp indicating when the board was last updated.

- **DeletedAt**: Timestamp indicating when the board was deleted (soft delete).

$\quad$ $\quad$**gorm:"index"**
## Relationships
- **Users**: A many-to-many relationship with the User struct, managed through the `user_board_roles` table. Deleting a board will cascade delete the associated records in the user_board_roles table.

$\quad$ $\quad$**gorm**:"many2many:user_board_roles;constraint:OnDelete:CASCADE;"
- **Tasks**: A one-to-many relationship with the `Task` struct. Each task is associated with a specific board. Deleting a board will cascade delete the associated tasks.

$\quad$ $\quad$**gorm**:"foreignKey:BoardID;constraint:OnDelete:CASCADE;"
- **UserBoardRoles**: A one-to-many relationship with the `UserBoardRole` struct, representing the roles users have on the board. Deleting a board will cascade delete the associated user board roles.

$\quad$ $\quad$**gorm**:"foreignKey:BoardID;constraint:OnDelete:CASCADE;"
- **Columns**: A one-to-many relationship with the `Column` struct. Each column is associated with a specific board. Deleting a board will cascade delete the associated columns.

$\quad$ $\quad$**gorm**:"foreignKey:BoardID;constraint:OnDelete:CASCADE;"

# Storage Package
An overview on the `boardRepo` struct and its methods for managing boards in a database using GORM.

## boardRepo Struct
The boardRepo struct provides methods to interact with the board storage.

### Methods
- **GetByID**: Retrieves a board by its ID.

- **GetUserBoards**: Retrieves the boards that a user is a member of, with pagination support.

- **GetPublicBoards**: Retrieves public boards, with pagination support.

- **Insert**: Inserts a new board into the database.

- **GetFullByID**: Retrieves a full board by its ID, including related users, columns, and tasks.

- **DeleteByID**: Deletes a board by its ID, including related tasks and their dependencies.

- **deleteTaskDependencies**: Deletes task dependencies for a list of task IDs.

## Internal/board
`types.go` includes:
- **Constants**: Defines types of boards.
```go
type BoardType string

const (
    Private BoardType = "private"
    Public  BoardType = "public"
)
```

- **Errors**: Various error messages used within the project.

- **Repo Interface**: The Repo interface defines methods for interacting with the board storage.

- **Board Struct**: Represents a board with its properties and related users and columns as a Domain Object.

- **Validations**:

$\quad$ $\quad$**ValidateBoardName**: Validates the name of a board.

`ops.go`:

- **Ops Struct**: The Ops struct provides methods for operating on boards using the Repo interface.

# Board Service

The `BoardService` handles operations related to boards, including creation, deletion, retrieving boards, and managing tasks and users on the board.

- **Errors**: Various errors are defined for handling permission and role issues.

- **BoardService Struct**: The BoardService struct encapsulates operations related to boards.

- **NewBoardService**: Creates a new instance of BoardService.

## Methods
- **GetFullBoardByID**: Retrieves a full board by its ID, checking if the user has permission to view the board if it is private. the public board will be returned in any case.

- **GetUserBoards**: Retrieves list of boards that a user is a member of, with pagination support.

- **GetPublicBoards**: Retrieves list of public boards, with pagination support.

- **CreateBoard**: Creates a new board and assigns the creator as the owner:

- **InviteUser**: Invites a user to the board, ensuring the inviter has the necessary permissions and the invitee is a user.

- **DeleteBoardByID**: Deletes a board by its ID, ensuring the user has the necessary permissions.

## Board Routes
Board Related routes are registered in `api/http/setup.go` using registerBoardRoutes.

## Handlers Operations:
- If there is any body in the request, any function retrieves the desired data using presenters out of body request.
- Calls the related service
- Converts output of a service to the desired response using presenter related structs.
- Some handlers need to be transactional like inviting users.