# alb2psql

A simple utility to pull AWS Application Load Balancer (ALB) request logs, and load them into a local PostgreSQL database for exploration/investigation.

> Note: If you're doing this alot, setting up a hosted workflow to automatically load ALB Logs -> Redshift/Postgres (optionally with Lambda), will be a much better option, and I suggest you consider it.

## Workflow

 - `alb2psql init`
 - `alb2psql fetch <options>`
 - `psql alb-logs`

## Installation

**For macOS:**

> Use [Homebrew](https://brew.sh/) package manager:
>
>      brew tap nrmitchi/alb2psql https://github.com/nrmitchi/alb2psqlgit
>      brew install alb2psql


**Other platforms:**

> Clone this repo, and run `go build`. Move the binary to somewhere in your path.

