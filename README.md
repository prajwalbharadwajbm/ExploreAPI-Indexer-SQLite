# Granular Recovery Optimization with SQLite

This project aims to optimize the process of granular recovery by efficiently managing the indexing and storage of files from various snapshots. The current implementation utilizes a lightweight and fast database system, SQLite, to handle these tasks, resulting in improved speed and reduced resource overhead.

## Problem

The frequent triggering of pod deployment for granular recovery, aimed at accessing different versions of files from various snapshots, led to significant resource consumption. To optimize this process, a more efficient approach was proposed by using SQLite as the database system for handling indexing and storage.

## Implementation

The implementation is based on two main components: `main.go` and `indexer.go`. The `main.go` file contains the core logic for traversing each snapshot, gathering essential data about each file, and inserting this data into the SQLite database. The `indexer.go` file provides the functionality for loading the index from the repository.

The `walker function` from the rapi library is used to traverse each snapshot efficiently. It operates within the scope of the "forAllSnapshot" function, which is imported from rapi. During the traversal, essential file information such as `fileName`, `path`, `ctime` (creation time), `mtime` (modification time), and `size` is collected. To uniquely identify each file version, a `fileID` is calculated using the concatenation of `bhash` (blob hash), `ctime`, `mtime`, and the file's `path`, which is then hashed with SHA256.

To store this data efficiently, it is inserted into an SQLite database. To handle duplicate entries when encountering the same `fileID` from different snapshots, the "*upsert*" operation (insert or update on conflict) is used to ensure seamless handling.

## How to Use

1. Clone the repository to your local machine.
2. Make sure you have Go (Golang) installed on your system.
3. Install the required dependencies by running the following command in the terminal:

   ```bash
   go mod tidy  
    ```
4. Modify the main.go file as needed to specify the repository location, password, and other options.
5. Run the application using the following command:
   ```bash
   go run main.go
   ```
## Observation
Based on the provided observations, the implementation shows promising results in terms of optimization. For a hundred snapshots, it took approximately 35 seconds to load and record the data into the SQLite database. Additionally, the solution was able to correctly identify and store different versions of a file from various backups, allowing for efficient granular recovery.
