# diffenvfiles

A simple command line program that compares two .env files and prints
out what's different between them. As a developer who's had to work on some
large projects, the "well it works fine with my .env" problem has gotten
very annoying.

This parses the two env files, sorts them, prints which variables
are different, and prints which variables are in only one file or the other.
Duplicate env vars in a file are treated as errors, because they're more trouble than they're
worth and usually a mistake.

# Usage

`diffenvfiles path/to/file1 path/to/file2`

# Installation

You can download a precompiled binary from github, and run `chmod u+x <filename>` on it, and
move it to somewhere like `/usr/local/bin`.

If you have go, you can run `go run main.go <file1> <file2>`

# Caveats

I just started working on this and wanted something within a couple hours, so
the code is ugly and there aren't many tests or automations just yet.
But so far it does the job!