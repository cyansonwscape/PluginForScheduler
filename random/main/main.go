// main.go 插件调用入口
package main

import (
	"os"

	"github.com/cyansonwscape/PluginForScheduler/random/plugin"
	"k8s.io/kubernetes/cmd/kube-scheduler/app"
)

func main() {
	command := app.NewSchedulerCommand(
		app.WithPlugin(plugin.Name, plugin.New),
	)

	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
