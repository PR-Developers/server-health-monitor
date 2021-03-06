# Server Health Monitor API

## Setup

Open a command prompt, create a `certs` directory in the root of this project. Run the following in the `certs` directory:

```bash
 openssl req -x509 -out localhost.crt -keyout localhost.key   -newkey rsa:2048 -nodes -sha256   -subj '/CN=localhost' -extensions EXT -config <( \
   printf "[dn]\nCN=localhost\n[req]\ndistinguished_name = dn\n[EXT]\nsubjectAltName=DNS:localhost\nkeyUsage=digitalSignature\nextendedKeyUsage=serverAuth")
```

Connect to your mongoDB database with a command prompt, execute the following:

```bash
use admin

db.createUser(
{
    user: "admin",
    pwd: "admin",
    roles: [ "root" ]
})
```
