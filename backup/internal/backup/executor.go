package backup

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/isutare412/web-memo/backup/internal/core/model"
)

type Executor struct {
	pgDump string
}

func NewExecutor() *Executor {
	return &Executor{
		pgDump: "pg_dump",
	}
}

func (e *Executor) BackupDatabase(ctx context.Context, req model.DatabaseBackupRequest) error {
	file, err := os.Create(req.BackupFilePath)
	if err != nil {
		return fmt.Errorf("creating backup file: %w", err)
	}
	defer file.Close()

	cmd := exec.CommandContext(ctx, e.pgDump)
	setBackupVariables(cmd, req, file)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("executing pg_dump: %w", err)
	}

	return nil
}

func (e *Executor) ReadFile(ctx context.Context, fileName string) (io.ReadCloser, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("opening file %s: %w", fileName, err)
	}
	return file, nil
}

func setBackupVariables(cmd *exec.Cmd, req model.DatabaseBackupRequest, destination io.Writer) {
	cmd.Env = []string{
		fmt.Sprintf("PGPASSWORD=%s", req.Password),
	}

	cmd.Args = []string{
		"-U", req.User,
		req.DatabaseName,
	}

	cmd.Stdout = destination
}
