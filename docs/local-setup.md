# Setup 

The only real component for a dev setup is a vault server running  
1. Simply run vault in dev mode
```
$ vault server -dev
```

2. The make file currently has some generation for mock data in order to test features. Will be updated as more features are added.
```
$ make setup-test-data
```
    * Note: This requires Vault tokent to be set in order to be able to write to Vault

3. Configure the env variables required to auth to vault or `.vaul7y.yaml` in your home directory

4. Make sure to set    
`VAULTY_LOG_FILE` env variable and point to a file, to log to a file 
`VAULTY_LOG_LEVEL` env variable - define the log level you want to use

```
❯ export VAULTY_LOG_LEVEL=debug
❯ export VAULTY_LOG_FILE=/tmp/my-vault-log.log
```