# aws-secret-env

`aws-secret-env` is a CLI tool that helps you inject secrets from [AWS Secrets
Manager](https://aws.amazon.com/secrets-manager/) into the environment variables
of a program. It also helps you create, view, add, and remove secrets in Secrets
Manager.

Install it by running:

```bash
go install github.com/ucarion/aws-secret-env
```

## What is this for?

`aws-secret-env` is meant to be used as the entrypoint in your production
services. Instead of having your service code fetch and decode secrets from AWS,
just have your service read secrets from env vars. Then wrap your service's
entrypoint with `aws-secret-env exec` to inject those secrets.

One great benefit of this approach is that your service will work as-is in local
dev, without any special casing. In prod, you can read from Secrets Manager; in
local dev, you can inject the testing passwords / API keys as environment
variables yourself.

So if normally your production service runs like this:

```bash
java -jar MyCoolService.jar
```

Instead, run it like this:

```bash
aws-secret-env exec my-cool-service -- java -jar MyCoolService.jar
```

That will decode a secret JSON from `my-cool-service` in secrets manager, and
invoke `java -jar MyCoolService.jar` with the values in that secret injected as
env vars visible to your program only.

## Injecting secrets

If you already have a secret formatted in a JSON object like this one in AWS
Secrets Manager:

```bash
aws secretsmanager get-secret-value --secret-id demo | jq .SecretString -r
```

```json
{"baz":"quux","foo":"bar"}
```

Then you can use `aws-secret-env exec YOUR_SECRET_NAME -- ...` to inject those
secrets into your environment variables:

```bash
# Before
env
```

```text
HOME=/root
TERM=xterm
PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
PWD=/
```

```bash
# After
aws-secret-env exec demo -- env
```

```text
HOME=/root
TERM=xterm
PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
PWD=/
baz=quux
foo=bar
```

## Managing secrets

In order to use `aws-secret-env` with a particular secret, you need to ensure
two things:

1. The secret's value must be in `SecretString`, not `SecretBinary`.
2. The secret's value must be a JSON object. The values of that object must all
   be strings.

To help you implement this convention, `aws-secret-env` comes with a few basic
commands for managing secrets you can use with `aws-secret-env exec`.

### Creating a secret

You can create a new secret by running:

```bash
aws-secret-env create example
```

This is basically just a wrapper around:

```bash
aws secretsmanager create-secret --name example --secret-string '{}'
```

### Viewing a secret

You can view the values in a secret by running:

```bash
aws-secret-env show example
```

This will output JSON, which can be useful if you're writing scripts that want
to work with JSON, instead of env vars.

The `show` command is basically just a wrapper around:

```bash
aws secretsmanager get-secret-value --secret-id example | jq .SecretString -r
```

### Adding (or replacing) a value to a secret

If you already have an existing secret in the format `aws-secret-env` requires,
you can add a new value to the secret by running:

```bash
# In the secret called "example", set "db-password" to "letmein".
aws-secret-env set example db-password letmein
```

Under the hood, this pulls down the value of the secret, upserts the given
key/value pair into the value, and then updates the secret with the new value.
For that reason, you'll need to have permission to both read and write the
secret in order for this command to work.

### Removing a value from a secret

To remove an existing value from a secret, you can run:

```bash
# In the secret called "example", remove "db-password" if it exists.
aws-secret-env unset example db-password
```

The `unset` command works very similarly to `set` under the hood. You'll need to
have permission to both read and write the secret in order for this command to
work.
