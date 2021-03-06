/*
Copyright 2017 the Heptio Ark contributors.

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

package restore

import (
	"os"
	"time"

	"github.com/spf13/cobra"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/heptio/ark/pkg/apis/ark/v1"
	"github.com/heptio/ark/pkg/client"
	"github.com/heptio/ark/pkg/cmd"
	"github.com/heptio/ark/pkg/cmd/util/downloadrequest"
	arkclient "github.com/heptio/ark/pkg/generated/clientset/versioned"
)

func NewLogsCommand(f client.Factory) *cobra.Command {
	timeout := time.Minute

	c := &cobra.Command{
		Use:   "logs RESTORE",
		Short: "Get restore logs",
		Args:  cobra.ExactArgs(1),
		Run: func(c *cobra.Command, args []string) {
			l := NewLogsOptions()
			cmd.CheckError(l.Complete(args))
			cmd.CheckError(l.Validate(f))
			arkClient, err := f.Client()
			cmd.CheckError(err)
			err = downloadrequest.Stream(arkClient.ArkV1(), f.Namespace(), args[0], v1.DownloadTargetKindRestoreLog, os.Stdout, timeout)
			cmd.CheckError(err)
		},
	}

	c.Flags().DurationVar(&timeout, "timeout", timeout, "how long to wait to receive logs")

	return c
}

// LogsOptions contains the fields required to retrieve logs of a restore
type LogsOptions struct {
	RestoreName string

	client arkclient.Interface
}

// NewLogsOptions returns a new instance of LogsOptions
func NewLogsOptions() *LogsOptions {
	return &LogsOptions{}
}

// Complete fills in LogsOptions with the given parameters, like populating the
// restore name from the input args
func (l *LogsOptions) Complete(args []string) error {
	l.RestoreName = args[0]
	return nil
}

// Validate validates the LogsOptions against the cluster, like validating if
// the given restore exists in the cluster or not
func (l *LogsOptions) Validate(f client.Factory) error {
	c, err := f.Client()
	if err != nil {
		return err
	}
	l.client = c

	_, err = l.client.ArkV1().Restores(f.Namespace()).Get(l.RestoreName, metav1.GetOptions{})
	return err
}
