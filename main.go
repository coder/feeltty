package main

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"text/tabwriter"
	"time"

	"github.com/spf13/pflag"
	"go.coder.com/cli"
)

type rootCmd struct {
}

func (r *rootCmd) Spec() cli.CommandSpec {
	return cli.CommandSpec{
		Name:  "feeltty",
		Usage: "<command [args ...]>",
		Desc:  "Assess the latency of a TTY interface",
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
	t := test(exec.Command(args[0], args[1:]...))
	wr := tabwriter.NewWriter(os.Stdout, 0, 2, 4, ' ', 0)
	fmt.Fprintf(wr, "connect\t%v\n", t.connect.took())
	var (
		tookMillis []float64
	)
	for _, it := range t.input {
		tookMillis = append(tookMillis, float64(it.took().Milliseconds()))
	}
	fmt.Fprintf(wr, "input sample size\t%v\n", len(tookMillis))
	fmt.Fprintf(wr, "input mean\t%v\n", time.Millisecond * time.Duration(mean(tookMillis)))
	fmt.Fprintf(wr, "input stddev\t%v\n", time.Millisecond * time.Duration(stddev(tookMillis)))
	wr.Flush()

}

func main() {
	cli.RunRoot(&rootCmd{})
}
