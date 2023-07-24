package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"time"

	"database/sql"

	"github.com/konidev20/rapi"
	"github.com/konidev20/rapi/backend"
	"github.com/konidev20/rapi/restic"
	"github.com/konidev20/rapi/walker"
	"github.com/rindex/indexer"

	_ "modernc.org/sqlite"
)

var ropts = rapi.DefaultOptions
var ctx = context.Background()

func main() {
	startTime := time.Now()
	db, err := openDatabase("D:/Rindex/Explorer API - Indexer/db/testdb.db")
	if err != nil {
		log.Fatal(err)
		os.Exit(5)
	}
	if err := OpenRepository(db); err != nil {
		log.Fatal(err)
		os.Exit(4)
	}
	log.Println("success")
	endTime := time.Now()
	elapsedTime := endTime.Sub(startTime)
	log.Println("Process completed in:", elapsedTime)
	defer db.Close()
}

func OpenRepository(db *sql.DB) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	ropts.Repo = "local:D:/[backup] - dirgenFiles"
	ropts.Password = "0000"
	repo, err := rapi.OpenRepository(ctx, ropts)
	if err != nil {
		return err
	}
	fmt.Printf("Opened repository: %v\n", repo)
	indexer.LoadIndex(ctx, repo)
	snapshotLister, err := backend.MemorizeList(ctx, repo.Backend(), restic.SnapshotFile)
	if err != nil {
		return err
	}

	err = restic.ForAllSnapshots(ctx, snapshotLister, repo, nil, func(id restic.ID, snap *restic.Snapshot, err error) error {

		if snap.Tree == nil {
			return fmt.Errorf("snapshot %v has no tree", snap.ID().Str())
		}

		err = walker.Walk(ctx, repo, *snap.Tree, nil, func(parentTreeID restic.ID, nodepath string, node *restic.Node, err error) (bool, error) {
			if err != nil {
				return false, walker.ErrSkipNode
			}
			if nodepath == "/" {
				return false, nil
			}
			_, fileID := hash(node, nodepath)
			fileName := node.Name
			path := nodepath
			ctime := node.ChangeTime
			mtime := node.ModTime
			size := node.Size
			err = insertFileInformation(db, fileID, fileName, path, ctime, mtime, size)
			if err != nil {
				fmt.Println("Error inserting file information:", err)
			}
			return true, nil
		})
		if err != nil {
			return err
		}

		return nil
	})
	fmt.Println(err)
	return err
}
func hash(node *restic.Node, path string) (bhash string, fileID string) {
	var bb []byte
	for _, c := range node.Content {
		bb = append(bb, []byte(c[:])...)
	}

	bh := sha256.Sum256(bb)
	bhash = hex.EncodeToString(bh[:])

	bb = append(bb, []byte(path)...)

	changeTimeBytes := []byte(node.ChangeTime.Format(time.RFC3339Nano))
	modTimeBytes := []byte(node.ModTime.Format(time.RFC3339Nano))
	pathBytes := []byte(path)

	bb = append(bb, changeTimeBytes...)
	bb = append(bb, modTimeBytes...)
	bb = append(bb, pathBytes...)

	fi := sha256.Sum256(bb)
	fileID = hex.EncodeToString(fi[:])

	return bhash, fileID
}

func createTables(db *sql.DB) error {
	query := `
        CREATE TABLE IF NOT EXISTS files (
            file_id TEXT PRIMARY KEY,
            name TEXT,
            path TEXT,
            ctime TEXT,
            mtime TEXT,
            size INTEGER
        )
    `
	_, err := db.Exec(query)
	return err
}

func openDatabase(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	if err := createTables(db); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
func insertFileInformation(db *sql.DB, fileID string, fileName string, path string, ctime time.Time, mtime time.Time, size uint64) error {

	query := `
		INSERT INTO files (file_id, name, path, ctime, mtime, size)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT (file_id)
		DO 
		UPDATE SET name =?, path =?, ctime =?, mtime =?, size =?;
	`

	_, err := db.Exec(query, fileID, fileName, path, ctime, mtime, size, fileName, path, ctime, mtime, size)
	if err != nil {
		return err
	}
	return nil
}
