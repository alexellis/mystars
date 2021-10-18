# mystars

A Go function that tells you about stars on projects you maintain.

Works on:
* repositories
* organisations

Get an incoming Slack URL to the channel of your choice, then create a secret for it:

```bash
faas-cli secret create slack-stars-webhook-url \
  --from-literal https://hooks.slack.com/services/abc/xyz/def
```

Then set a secret for the webhook that GitHub will use to calculate a digest for each webhook:

```bash
faas-cli secret create slack-stars-hmac-secret \
  --from-literal "so-random"
```

Then deploy the function and set up a webhook for the "Watch" (Star) and "Fork" event using your public URL for OpenFaaS.

If you are running this on your Raspberry Pi or home server, use [inlets](https://inlets.dev) to get a public URL.

[faasd](https://github.com/openfaas/faasd) is a cheap and efficent way to run openfaas on a public cloud VM.
