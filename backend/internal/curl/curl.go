package curl

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"volnerability-game/internal/db"
	"volnerability-game/internal/domain"
)

// curl "http://hosting.com/about"

const base = "http://hosting.com/about"

// base GET
// curl "http://example.com/api/user?id=1"

// ----------------

// base POST
// curl -X POST "http://example.com/api/login" \
//     -d "username=admin" \
//     -d "password=secret123"

func runCmd(input string) (domain.Response, error) {
	cmd := exec.Command("sh", "-c", input)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return domain.Response{}, fmt.Errorf("failed to run cmd due err: %w", err)
	}
	resp := domain.Response{}
	if err := json.Unmarshal(out, &resp); err != nil {
		return domain.Response{}, err
	}
	return resp, nil
}

func NewTask(db *db.Storage, levelId int, input string) func(context.Context) (domain.Response, error) {
	return func(context.Context) (domain.Response, error) {
		return runCmd(input)
	}
}
