package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/manifoldco/promptui"
	"github.com/urfave/cli/v3"

	"github.com/alexpts/go-env-manager/internal/build"
	"github.com/alexpts/go-env-manager/internal/config"
)

func main() {
	cmd := &cli.Command{
		Version:               "0.0.1",
		EnableShellCompletion: true,

		Name:  "use",
		Usage: "активация профиля",
		Action: func(context.Context, *cli.Command) error {
			presets := config.LoadFromConfig()

			prompt := promptui.Select{
				Label: "Какой из env профилей активировать?",
				Items: presets.Keys(),
			}
			_, profileName, err := prompt.Run()
			if err != nil {
				fmt.Printf("Prompt failed %v\n", err)
			}

			preset := presets[profileName]
			config.ApplyConfig(preset)

			return nil
		},

		Commands: []*cli.Command{
			{
				Name:        "list",
				Aliases:     []string{"ls"},
				Usage:       "список профилей",
				Description: "Показывает список всех известных профилей go-env в менеджере профилей",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					presets := config.LoadFromConfig()
					if len(presets) == 0 {
						appName := filepath.Base(os.Args[0])
						fmt.Printf("Нет профилей в системе. Вы можете создать новый профиль из текущих настроек командой `%s %s`\n", appName, "create")
						return nil
					}

					fmt.Printf("Посмотреть детали профилей можно в конфиге %s\n", config.PathConfig)
					for name := range presets {
						fmt.Printf("- %s\n", name)
					}

					return nil
				},
			},
			{
				Name:  "version",
				Usage: "вывести данные о версии",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					fmt.Printf("%+v\n", build.GetWithBuildFlags())
					return nil
				},
			},
			{
				Name:        "reset",
				Usage:       "сбросить настройки",
				Description: "Сбрасывает настройки части go env переменных на стандартные",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					for _, name := range []string{"GOPRIVATE", "GOPROXY", "GONOPROXY", "GONOSUMDB", "GOSUMDB"} {
						err := (exec.Command("go", "env", "-u", name)).Run()
						if err != nil {
							fmt.Println(err)
						}
					}

					return nil
				},
			},
			{
				Name:        "create",
				Aliases:     []string{"new"},
				Usage:       "создать профиль",
				Description: "Создает профиль из текущих настроек `go env`",
				Arguments: []cli.Argument{
					&cli.StringArg{
						Name: "profile-name",
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					profileName := cmd.StringArg("profile-name")

					isValid, _ := regexp.MatchString("^[a-z][a-z0-9-_]+$", profileName)
					if !isValid {
						fmt.Printf("Ошибка валидации имени профиля. Маска имени профиля: `^[a-z][a-z0-9-_]+$`")
						return nil
					}

					profiles := config.LoadFromConfig()
					profile, isOk := profiles[profileName]

					newPreset := config.FromProcessEnv()

					if isOk {
						fmt.Printf(
							"Профиль в системе:\n%s\nСоздавайемый профиль:\n%s\n",
							profile.String(),
							newPreset.String(),
						)

						prompt := promptui.Select{
							Label: "Профиль уже существует, заместить его текущими настройками?",
							Items: []string{"yes", "no"},
						}
						_, result, err := prompt.Run()
						if err != nil {
							return err
						}

						if result == "no" {
							fmt.Printf("Вы не подтвердили замену профиля\n")
							return nil
						}
					}

					profiles[profileName] = newPreset

					err := profiles.Persist()
					if err != nil {
						fmt.Printf("Ошибка при сохранении конфигурации на диск\n")
						return err
					}

					fmt.Printf("Конфигурация была успешно сохранена в файл %s\n", config.PathConfig)

					return nil
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
