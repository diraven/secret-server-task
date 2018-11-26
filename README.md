# Secret Server Coding Task

## Notes

* Server currently holds all the secrets in memory using json file for persistence.
* Given above with the big enough amount of queries server will be severely bound by disk IO operations.
* Memory usage depends on the amount of secrets and given big enough secrets count server can run out of memory.
* Above issues can be resolved by switching to a proper DB engine.
* Although server does not have any protection against brute-force attacks, all the 404 requests are logged along with the caller IP. This can be used by fail2ban to implement some basic brute-force protection.
* There are a couple of hardcoded limitations to the server to make sure it does not crash the server due to the memory usage - number of secrets is limited to 1000 and maximum message size is 1024 symbols. These can be lifted given proper limitations are designed and implemented.