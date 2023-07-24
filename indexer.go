package indexer

import (
	"context"
	"fmt"
	"os"

	"github.com/konidev20/rapi/repository"
	"github.com/rs/zerolog"
)

var log = zerolog.New(os.Stdout).With().Timestamp().Logger()

func LoadIndex(ctx context.Context, repo *repository.Repository) (IndexStats, error) {
	stats := NewStats()

	log.Debug().Msg("loading index")
	if err := repo.LoadIndex(ctx); err != nil {
		return stats, err
	}
	log.Debug().Msg("index loaded")

	if err := repo.Close(); err != nil {
		fmt.Printf("Error closing repository: %v\n", err)
	}
	return stats, nil
}

// func ListSnapshots(ctx context.Context, repo *repository.Repository) (<-chan *restic.Snapshot, error) {
// 	out := make(chan *restic.Snapshot)
// 	var errfound error
// 	go func() {
// 		defer close(out)

// 		snapshots := []*restic.Snapshot{}

// 		err := repo.List(ctx, restic.SnapshotFile, func(id restic.ID, size int64) error {
// 			sn, err := restic.LoadSnapshot(ctx, repo, id)
// 			if err != nil {
// 				log.Error().Err(err).Msgf("could not load snapshot %v: %v\n", id.Str(), err)
// 				return nil
// 			}
// 			snapshots = append(snapshots, sn)
// 			return nil
// 		})

// 		if err != nil {
// 			errfound = err
// 			return
// 		}

// 		for _, sn := range snapshots {
// 			select {
// 			case <-ctx.Done():
// 				return
// 			case out <- sn:
// 			}
// 		}
// 	}()

// 	return out, errfound
// }
