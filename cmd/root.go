/*
Copyright © 2020 Patrick Miziewicz <patrick.miziewicz@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import "github.com/spf13/cobra"

func RootCmd() (*cobra.Command, error) {
  cmd := &cobra.Command{
    Use:          "aws-chaos-cli",
    Short:        "CLI tool to inject failure into AWS infrastructure",
    Long:         "CLI tool created with the sole purpose of injecting failure into AWS infrastructure in a convenient manner",
  }

  cmd.AddCommand(
    failAzCmd(),
  )

  return cmd, nil
}