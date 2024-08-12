## ghibp password

Lookup the given passwords for breaches in the HaveIBeenPwned database

### Synopsis

Query the HaveIBeenPwned password database for each of the given
passwords and return the number of hits. Passwords are not transmitted in cleartext
but only the first 5 digits of its sha1 hash are sent to the servers and the rest of
the lookup is done locally.

```
ghibp password PASSWORD... [flags]
```

### Options

```
  -h, --help   help for password
```

### SEE ALSO

* [ghibp](ghibp.md)	 - HaveIBeenPwned public API util

