
Refactored (
    /src/controllers/events/delete.go RemoveByIdWithParticipants function:
        -> Row affected output removed
        -> sqlx rebinder removed and replaced by switch
        -> Now with context
)

Add (
    /src/database/mysql.sql 
        -> Statements for manual construction
)

Renamed (
    /src/database/statements.sql 
        -> psql.sql
)

Mod (
    /src/controllers/persistence/persistence.go
        -> Regular expression now accepts something like 'PSQL', 'psql'
)