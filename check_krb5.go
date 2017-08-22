// krb5perf is a tool to perform performance benchmarking and stress testing of Kerberos v5 KDC AS_REQ functions
package main

import (
	"errors"
	"fmt"
	"github.com/alexflint/go-arg"
	"github.com/cobaugh/krb5-go"
	"github.com/newrelic/go_nagios"
	"os"
	"time"
)

// arguments
type Args struct {
	Keytab   string `arg:"env:KTNAME,-k"`
	Password string `arg:"-P"`
	Client   string `arg:"-c,required"`
	Service  string `arg:"-s,required"`
	Count    int    `arg:"-c:Number of checks to perform"`
	Interval string `arg:"-i,help:Time to wait between each check"`
	Warn     string `arg:"-W,help:Warning threshold"`
	Crit     string `arg:"-C,help:Critical threshold"`
}

func (Args) Version() string {
	return os.Args[0] + " 0.1"
}

func main() {
	var args Args
	args.Count = 1
	arg.MustParse(&args)

	// set up defaults for interval, warning, critical
	var interval time.Duration
	var warn time.Duration
	var crit time.Duration
	if args.Interval == "" {
		args.Interval = "1s"
	}
	if args.Warn == "" {
		args.Warn = "1s"
	}
	if args.Crit == "" {
		args.Crit = "5s"
	}

	// parse interval, warn, crit into time.Duration
	interval, err := time.ParseDuration(args.Interval)
	if err != nil {
		nagios.Unknown(fmt.Sprintf("%s", err))
	}
	warn, err = time.ParseDuration(args.Warn)
	if err != nil {
		nagios.Unknown(fmt.Sprintf("%s", err))
	}
	crit, err = time.ParseDuration(args.Crit)
	if err != nil {
		nagios.Unknown(fmt.Sprintf("%s", err))
	}

	// sanity-check warn/crit
	if (crit < warn) {
		nagios.Unknown("Critical threshold must be less than Warning")
	}

	// create a shared context
	ctx, err := krb5.NewContext()
	if err != nil {
		nagios.Unknown(fmt.Sprintf("NewContext(): %s", err))
	}
	defer ctx.Free()

	// check for keytab or password
	var keytab *krb5.KeyTab
	if args.Keytab != "" {
		keytab, err := ctx.OpenKeyTab(args.Keytab)
		if err != nil {
			nagios.Unknown(fmt.Sprintf("OpenKeyTab(): %s", err))
		}
		defer keytab.Close()
	} else if args.Password == "" {
		nagios.Unknown("One of either --password or --keytab must be specified")
	}

	// set up our client and service
	client, err := ctx.ParseName(args.Client)
	if err != nil {
		nagios.Unknown(fmt.Sprintf("ParseName(): %s", err))
	}

	service, err := ctx.ParseName(args.Service)
	if err != nil {
		nagios.Unknown(fmt.Sprintf("ParseName(): %s", err))
	}

	var results []time.Duration
	for i := 0; i < args.Count; i++ {
		start := time.Now()
		if args.Password != "" {
			_, err = ctx.GetInitialCredentialWithPassword(args.Password, client, service)
		} else {
			_, err = ctx.GetInitialCredentialWithKeyTab(keytab, client, service)
		}
		elapsed := time.Since(start)

		if err != nil {
			nagios.Critical(err)
		}
		results = append(results, elapsed)
		time.Sleep(interval)
	}

	// average, min, max
	sum := time.Duration(0)
	min := results[0]
	max := time.Duration(0)
	for i := range results {
		sum += results[i]
		if results[i] > max {
			max = results[i]
		}
		if results[i] < min {
			min = results[i]
		}
	}
	avg := sum / time.Duration((args.Count))

	perfdata := fmt.Sprintf("t_avg=%f:%f:%f:%f:%f", avg.Seconds(), warn.Seconds(), crit.Seconds(), min.Seconds(), max.Seconds())
	status := fmt.Sprintf("Authenticated as %s to %s (avg: %v, min: %v, max: %v, i: %d) | %s\n",
		args.Client, args.Service, avg, min, max, args.Count, perfdata)

	if avg > crit {
		nagios.Critical(errors.New(status))
	} else if avg > warn {
		nagios.Warning(status)
	} else {
		nagios.Ok(status)
	}
}
