# Secret Server Coding Task

## Notes

* Server currently holds all the secrets in memory using json file for persistence.
* Given above with the big enough amount of queries server will be severely bound by disk IO operations.
* Memory usage depends on the amount of secrets and given big enough secrets count server can run out of memory.
* Above issues can be resolved by switching to a proper DB engine.
* Although server does not have any protection against brute-force attacks, all the 404 requests are logged along with the caller IP. This can be used by fail2ban to implement some basic brute-force protection.