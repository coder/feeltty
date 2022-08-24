package main

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"text/tabwriter"

	"github.com/spf13/pflag"
	"go.coder.com/cli"
)

type rootCmd struct {
	iterations int64
}

func (r *rootCmd) Spec() cli.CommandSpec {
	return cli.CommandSpec{
		Name:    "feeltty",
		Usage:   "[flags] <command [args ...]>",
		Desc:    "Assess the latency of a TTY interface",
		RawArgs: true,
	}
}

func mean(fs []float64) float64 {
	var r float64
	for _, f := range fs {
		r += f
	}
	return r / float64(len(fs))
}

func stddev(fs []float64) float64 {
	m := mean(fs)

	var rs []float64
	for _, f := range fs {
		rs = append(rs, math.Pow(f-m, 2))
	}
	return math.Sqrt(mean(rs))
}

func (r *rootCmd) Run(fl *pflag.FlagSet) {
	args := fl.Args()
	if len(args) == 0 {
		fl.Usage()
		os.Exit(1)
	}
	// TODO: make this configurable.
	r.iterations = 32

	t := test(exec.Command(args[0], args[1:]...), int(r.iterations))
	wr := tabwriter.NewWriter(os.Stdout, 0, 2, 4, ' ', 0)
	fmt.Fprintf(wr, "connect\t%0.5fms\n", t.connect.took().Seconds()*1000)
	var (
		tooks []float64
	)
	for _, it := range t.input {
		// flog.Infof("%+v\n", it.took())
		tooks = append(tooks, it.took().Seconds())
	}
	fmt.Fprintf(wr, "input sample size\t%v\n", len(tooks))
	fmt.Fprintf(wr, "input mean\t%0.3fms\n", mean(tooks)*1000)
	fmt.Fprintf(wr, "input stddev\t%0.3fms\n", stddev(tooks)*1000)
	wr.Flush()
}

func main() {
	cli.RunRoot(&rootCmd{})
}
