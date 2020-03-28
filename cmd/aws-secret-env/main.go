package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/ucarion/aws-secret-env/internal/secrets"
	"golang.org/x/sys/unix"
)

var (
	rootCmd = &cobra.Command{
		Use:   "aws-secret-env",
		Short: "A CLI tool that helps you inject secrets from AWS Secrets Manager into env vars",
	}

	execCmd = &cobra.Command{
		Use:   "exec [secret] -- [cmd...]",
		Short: "Run a command with secrets injected as env vars",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			do(func() error {
				secrets, err := secrets.Get(cmd.Context(), args[0])
				if err != nil {
					return err
				}

				env := os.Environ()
				for k, v := range secrets {
					env = append(env, fmt.Sprintf("%s=%s", k, v))
				}

				argv0, err := exec.LookPath(args[1])
				if err != nil {
					return err
				}

				return unix.Exec(argv0, args[2:], env)
			})
		},
	}

	showCmd = &cobra.Command{
		Use:   "show [secret]",
		Short: "Output the values in a secret as JSON",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			do(func() error {
				secret, err := secrets.Get(cmd.Context(), args[0])
				if err != nil {
					return err
				}

				encoder := json.NewEncoder(os.Stdout)
				if err := encoder.Encode(secret); err != nil {
					return err
				}

				return nil
			})
		},
	}

	createCmd = &cobra.Command{
		Use:   "create [secret]",
		Short: "Create a new, empty secret",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			do(func() error {
				return secrets.Create(cmd.Context(), args[0])
			})
		},
	}

	setCmd = &cobra.Command{
		Use:   "set [secret] [key] [value]",
		Short: "Set or update a key/value pair in a secret",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			do(func() error {
				return secrets.Set(cmd.Context(), args[0], args[1], args[2])
			})
		},
	}

	unsetCmd = &cobra.Command{
		Use:   "unset [secret] [key]",
		Short: "Remove a key/value pair from a secret",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			do(func() error {
				return secrets.Unset(cmd.Context(), args[0], args[1])
			})
		},
	}
)

func init() {
	rootCmd.AddCommand(execCmd, showCmd, createCmd, setCmd, unsetCmd)
}

func main() {
	do(rootCmd.Execute)
}

func do(f func() error) {
	if err := f(); err != nil {
		fmt.Fprintln(os.Stdout, err)
		os.Exit(1)
	}
}
