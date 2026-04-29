package index

import (
	"database/sql"
	"fmt"
)

type GroupCount struct {
	Name  string
	Count int
}

type HistoryStats struct {
	TotalEntries int
	ByShell      []GroupCount
	BySource     []GroupCount
}

func QueryHistoryStats(db *sql.DB) (HistoryStats, error) {
	if db == nil {
		return HistoryStats{}, fmt.Errorf("query history stats: database is required")
	}

	var stats HistoryStats

	if err := db.QueryRow(`SELECT COUNT(*) FROM history_entries;`).Scan(&stats.TotalEntries); err != nil {
		return HistoryStats{}, fmt.Errorf("query history stats: %w", err)
	}

	byShell, err := queryGroupedCounts(db, `
		SELECT shell, COUNT(*)
		FROM history_entries
		GROUP BY shell
		ORDER BY COUNT(*) DESC, shell ASC;
	`)
	if err != nil {
		return HistoryStats{}, fmt.Errorf("query history stats: %w", err)
	}
	stats.ByShell = byShell

	bySource, err := queryGroupedCounts(db, `
		SELECT source_file, COUNT(*)
		FROM history_entries
		GROUP BY source_file
		ORDER BY COUNT(*) DESC, source_file ASC;
	`)
	if err != nil {
		return HistoryStats{}, fmt.Errorf("query history stats: %w", err)
	}
	stats.BySource = bySource

	return stats, nil
}

func queryGroupedCounts(db *sql.DB, query string) ([]GroupCount, error) {
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var counts []GroupCount
	for rows.Next() {
		var item GroupCount
		if err := rows.Scan(&item.Name, &item.Count); err != nil {
			return nil, err
		}
		counts = append(counts, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return counts, nil
}
