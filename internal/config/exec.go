package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"sync"
)

var (
	goEnv     map[string]string
	goEnvOnce sync.Once
)

func loadGoEnv() {
	cmd := exec.Command("go", "env", "-json")
	out, err := cmd.Output()

	if err != nil {
		log.Fatal("cat`t exec go env")
	}

	err = json.Unmarshal(out, &goEnv)
	if err != nil {
		log.Fatal("cat`t parse go env")
	}
}

func GetGoEnv(key string) string {
	goEnvOnce.Do(loadGoEnv)
	return goEnv[key]
}

func ApplyConfig(config Config) {
	params := []struct {
		Key string
		Val string
	}{
		{"GOPRIVATE", config.GoPrivate},
		{"GOPROXY", config.GoProxy},
		{"GOSUMDB", config.GoSumDB},
		{"GONOSUMDB", config.GoSumDB},
		{"GOSUMDB", config.GoSumDB},
	}

	for _, param := range params {
		cmd := exec.Command("go", "env", "-w", fmt.Sprintf("%s=%s", param.Key, param.Val))
		err := cmd.Run()
		if err != nil {
			log.Fatalf("cat`t set config %s=%s", param.Key, param.Val)
		}
	}
}
