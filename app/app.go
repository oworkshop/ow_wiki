package app

import (
	"flag"
	"os"
	"strings"

	"github.com/spf13/pflag"
	"k8s.io/klog/v2"

	"github.com/oworkshop/ow_wiki/app/cmd"
)

func Run() error {
	klog.InitFlags(nil)
	pflag.CommandLine.SetNormalizeFunc(WordSepNormalizeFunc)
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)

	pflag.Set("logtostderr", "true")
	// We do not want these flags to show up in --help
	// These MarkHidden calls must be after the lines above
	pflag.CommandLine.MarkHidden("add-dir-header")
	pflag.CommandLine.MarkHidden("log-file")
	pflag.CommandLine.MarkHidden("log-file-max-size")
	pflag.CommandLine.MarkHidden("one-output")
	pflag.CommandLine.MarkHidden("version")
	pflag.CommandLine.MarkHidden("log-flush-frequency")
	pflag.CommandLine.MarkHidden("alsologtostderr")
	pflag.CommandLine.MarkHidden("log-backtrace-at")
	pflag.CommandLine.MarkHidden("log-dir")
	pflag.CommandLine.MarkHidden("logtostderr")
	pflag.CommandLine.MarkHidden("stderrthreshold")
	pflag.CommandLine.MarkHidden("vmodule")
	pflag.CommandLine.MarkHidden("skip-headers")
	pflag.CommandLine.MarkHidden("skip-log-headers")
	pflag.CommandLine.MarkHidden("v")

	cmd := cmd.NewCommand(os.Stdin, os.Stdout, os.Stderr)
	return cmd.Execute()
}

// WordSepNormalizeFunc changes all flags that contain "_" separators
func WordSepNormalizeFunc(f *pflag.FlagSet, name string) pflag.NormalizedName {
	if strings.Contains(name, "_") {
		return pflag.NormalizedName(strings.Replace(name, "_", "-", -1))
	}
	return pflag.NormalizedName(name)
}
