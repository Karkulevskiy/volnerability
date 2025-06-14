package curl

import (
	"bytes"
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
	outBytes, err := cmd.CombinedOutput()
	if err != nil {
		return domain.Response{}, fmt.Errorf("failed to run cmd due err: %w", err)
	}
	//   % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
	//                                  Dload  Upload   Total   Spent    Left  Speed
	// 100    46  100    46    0     0  69802      0 --:--:-- --:--:-- --:--:-- 46000
	// {"status":"200. StatusOK","isCompleted":true}
	fmt.Printf("curl response: %s\n", string(outBytes))
	ind := bytes.LastIndexByte(outBytes, '{')
	resp := domain.Response{}
	if err := json.Unmarshal(outBytes[ind:], &resp); err != nil {
		return domain.Response{}, err
	}
	return resp, nil
}

func NewTask(db *db.Storage, levelId int, input string) func(context.Context) (domain.Response, error) {
	return func(context.Context) (domain.Response, error) {
		const op = "curl.NewTask"
		resp, err := runCmd(input)
		if err != nil {
			return resp, err
		}
		if resp.CurlLevelId != levelId {
			return domain.Response{}, fmt.Errorf("%s: requested levelId: %d, response levelId: %d", op, levelId, resp.CurlLevelId)
		}
		return resp, nil
	}
}
