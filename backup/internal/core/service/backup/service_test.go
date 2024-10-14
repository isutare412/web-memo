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

func Test_backupKey(t *testing.T) {
	type args struct {
		t        time.Time
		prefix   string
		database string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "normal case",
			args: args{
				t:        time.Date(2024, 1, 1, 1, 1, 0, 0, time.UTC),
				prefix:   "some/path/to/destination",
				database: "foo",
			},
			want: "some/path/to/destination/20240101010100-backup-foo.sql",
		},
		{
			name: "slash ended prefix",
			args: args{
				t:        time.Date(2024, 1, 1, 1, 1, 0, 0, time.UTC),
				prefix:   "some/path/to/destination/",
				database: "foo",
			},
			want: "some/path/to/destination/20240101010100-backup-foo.sql",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := backupKey(tt.args.t, tt.args.prefix, tt.args.database)
			assert.Equal(t, tt.want, got)
		})
	}
}
