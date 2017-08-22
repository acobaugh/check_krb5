Nagios check for krb5 with ping-like avg, min, max output. Accepts a password or keytab.

## Usage
```
check_krb5 0.1
Usage: check_krb5 [--keytab KEYTAB] [--password PASSWORD] --client CLIENT --service SERVICE [--count COUNT] [--interval INTERVAL] [--warn WARN] [--crit CRIT]

Options:
  --keytab KEYTAB, -k KEYTAB
  --password PASSWORD, -P PASSWORD
  --client CLIENT, -c CLIENT
  --service SERVICE, -s SERVICE
  --count COUNT, -c COUNT [default: 1]
  --interval INTERVAL, -i INTERVAL
                         Time to wait between each check
  --warn WARN, -W WARN   Warning threshold
  --crit CRIT, -C CRIT   Critical threshold
  --help, -h             display this help and exit
  --version              display version and exit
```

## Example output
```
OK: Authenticated as testppp to krbtgt/dce.psu.edu (avg: 20.856253ms, min: 16.92726ms, max: 35.077536ms, i: 5) | t_avg=0.020856:1.000000:1.000000:0.016927:0.035078
```

## TODO
