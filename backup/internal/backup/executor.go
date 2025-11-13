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

	var (
		stdout = new(bytes.Buffer)
		stderr = new(bytes.Buffer)
	)
	setBackupVariables(cmd, req, stdout, stderr)

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("executing pg_dump: %w: %s", err, stderr.String())
	}

	result, err := io.ReadAll(stdout)
	if err != nil {
		return nil, fmt.Errorf("reading pg_dump output: %w", err)
	}

	return bytes.NewReader(result), nil
}

func setBackupVariables(cmd *exec.Cmd, req model.DatabaseBackupRequest, stdout, stderr io.Writer) {
	cmd.Env = []string{
		fmt.Sprintf("PGPASSWORD=%s", req.Password),
	}

	cmd.Args = append(cmd.Args,
		"-h", req.Host,
		"-p", strconv.Itoa(req.Port),
		"-U", req.User,
		req.DatabaseName,
	)

	cmd.Stdout = stdout
	cmd.Stderr = stderr
}
