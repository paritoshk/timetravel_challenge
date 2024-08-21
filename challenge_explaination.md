# Challenge Explanation

I approached this challenge by focusing on two main objectives:

1. Switching the storage backend to SQLite
2. Adding time travel functionality to the API

## Objective 1: Switch To SQLite

I modified the existing in-memory storage to use SQLite as the persistent backend. This involved:

- Creating a new SQLite database file
- Implementing CRUD operations using SQL queries
- Ensuring data persistence across server restarts

The main changes were made in the `service/record.go` file, where I implemented the `SQLiteRecordService` struct and its methods.

## Objective 2: Add Time Travel

To add the time travel functionality, I:

1. Created new `/api/v2` endpoints
2. Implemented versioning for records
3. Added methods to retrieve specific versions of records
4. Ensured backward compatibility with `/api/v1` endpoints

The main changes for this objective were spread across multiple files, including `api/api.go`, `api/get_records.go`, and `api/post_records.go`.

## Testing

I expanded the existing test suite in `main_test.go` to cover both v1 and v2 API endpoints. The tests ensure that:

- Records can be created, retrieved, and updated
- Versioning works correctly for v2 endpoints
- Time travel functionality (retrieving specific versions) works as expected

The test results showed that all implemented features are working correctly, with all tests passing successfully.

## Objective 1: Switch To SQLite

I replaced the in-memory storage with SQLite persistence. The implementation can be found in:

go:service/record.go
startLine: 1
endLine: 205


To test this objective:

1. Run the server
2. Create and modify records
3. Restart the server
4. Verify that the data persists after restart

## Objective 2: Add Time Travel

I implemented time travel functionality in the v2 API. Key changes include:

1. New endpoints in `api/api.go`
2. Versioning logic in `api/post_records.go`
3. Version retrieval in `api/get_records.go`

To test this objective:

1. Use the `/api/v2/records/{id}` endpoint to create and update records
2. Retrieve specific versions using `/api/v2/records/{id}?version={version}`
3. Get all versions of a record using `/api/v2/records/{id}/versions`

## Testing
The `main_test.go` file contains comprehensive tests for both objectives. To run the tests:
```bash
go test ./...
```

All tests are passing, which confirms that:

1. SQLite persistence is working correctly
2. Time travel functionality is implemented and working as expected
3. Both v1 and v2 APIs are functioning properly

The test file covers various scenarios, including:

- Creating records (v1 and v2)
- Retrieving records (v1 and v2)
- Updating records (v1 and v2)
- Deleting fields (v1)
- Retrieving specific versions (v2)
- Getting all versions of a record (v2)

The passing tests indicate that the implementation meets all the requirements of the challenge.

# Appendix (AI Written)
The tests cover both v1 and v2 API endpoints to ensure backward compatibility and new functionality.
We test record creation, retrieval, and updating to verify basic CRUD operations.
Version-specific tests (v2) ensure that the time travel functionality works correctly.
The tests use a separate test database to avoid interfering with the main application data.
We clean up the test database after running all tests to maintain a clean state for future test runs.
The tests pass because:
The SQLite implementation correctly persists data and handles all required operations.
The versioning system in v2 correctly creates new versions on updates and allows retrieval of specific versions.
Error handling is implemented correctly, returning appropriate status codes and error messages.
The v1 API continues to function as before, maintaining backward compatibility.
This comprehensive test suite ensures that both objectives of the challenge are met and functioning correctly.