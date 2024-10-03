package backup

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os/exec"
	"strconv"

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

func (e *Executor) BackupDatabase(ctx context.Context, req model.DatabaseBackupRequest) (backup io.Reader, err error) {
	cmd := exec.CommandContext(ctx, e.pgDump)

	stdout := new(bytes.Buffer)
	setBackupVariables(cmd, req, stdout)

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("executing pg_dump: %w", err)
	}

	return stdout, nil
}

func setBackupVariables(cmd *exec.Cmd, req model.DatabaseBackupRequest, destination io.Writer) {
	cmd.Env = []string{
		fmt.Sprintf("PGPASSWORD=%s", req.Password),
	}

	cmd.Args = []string{
		"-h", req.Host,
		"-p", strconv.Itoa(req.Port),
		"-U", req.User,
		req.DatabaseName,
	}

	cmd.Stdout = destination
}
