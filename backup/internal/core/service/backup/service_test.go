package backup

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_pickBackupNamesBefore(t *testing.T) {
	type args struct {
		backups []string
		t       time.Time
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "normal case",
			args: args{
				backups: []string{
					"backups/postgresql/20240101000000-backup-foo.sql",
					"backups/postgresql/20240105000000-backup-foo.sql",
					"backups/postgresql/20240110000000-backup-foo.sql",
					"backups/postgresql/20240115000000-backup-foo.sql",
					"backups/postgresql/20240120000000-backup-foo.sql",
				},
				t: time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
			},
			want: []string{
				"backups/postgresql/20240101000000-backup-foo.sql",
				"backups/postgresql/20240105000000-backup-foo.sql",
			},
		},
		{
			name: "result should be sorted in ascending order",
			args: args{
				backups: []string{
					"backups/postgresql/20240120000000-backup-foo.sql",
					"backups/postgresql/20240115000000-backup-foo.sql",
					"backups/postgresql/20240110000000-backup-foo.sql",
					"backups/postgresql/20240105000000-backup-foo.sql",
					"backups/postgresql/20240101000000-backup-foo.sql",
				},
				t: time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
			},
			want: []string{
				"backups/postgresql/20240101000000-backup-foo.sql",
				"backups/postgresql/20240105000000-backup-foo.sql",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := pickBackupKeysBefore(tt.args.backups, tt.args.t)
			assert.Equal(t, tt.want, got)
		})
	}
}
