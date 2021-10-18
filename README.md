# stargazers

This is a function written in Go that sends GitHub Star and Fork events to a Slack channel.

You'll be able to monitor which of your projects are popular and guess when a new Pull Request may be on its way by watching the events on Slack.

It'll work on multiple:

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

## Further resources

If you are running this on your Raspberry Pi or home server, use [inlets](https://inlets.dev) to get a public URL.

[faasd](https://github.com/openfaas/faasd) is a cheap and efficent way to run openfaas on a public cloud VM.

Learn Go with: ["Everyday Go" is the fast way to learn tools, techniques and patterns from real tools used in production.](https://openfaas.gumroad.com/l/everyday-golang) 

Learn to build similar functions using OpenFaaS and faasd: [Serverless for Everyone Else](https://gumroad.com/l/serverless-for-everyone-else)
