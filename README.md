So I had started working with the twitter api for a project. The most annoying part was getting setup and authorizing the first request. So I thought I would writeup a simple version of it here in go using few dependencies and explain what I did and why (and maybe where the docs are a bit misleading).

You could use this like a simple library but I won't make any guarantee as to it's reliability.

1- Generate Keys on twitter.
   - Go here: https://developer.twitter.com/en/portal/dashboard
   - Create and then go to your app page, there will be two tabs "settings" & "keys and tokens"
   - Enable Oauth on "settings" tab, probably with read write, make sure to set callback as localhost:8080/callback
   - Go to "keys and tokens" and generate a set. you need to save these somwhere, maybe a password manager

2- Set the key values as envorinment variables.
   - export TW_TOKEN="GENERATEDTOKEN" TW_TOKENSECRET="GENERATEDTOKENSECRET" TW_CONSUMERKEY=yourconsumerkey TW_CONSUMERSECRET=yourconsumersecret

3- Run `go run main.go` in that shell environment.

4- Go to: localhost:8080/auth to start the oauth flow.

5- View results in the shell/terminal log.

Sample output (replaced sensitive stuff with XXX):

step 1 oauthflow: 
request header: OAuth oauth_callback="http%3A%2F%2Flocalhost%3A8080%2Fcallback",oauth_consumer_key="XXX",oauth_nonce="fe79d0a19a4af2f3f2f7de9c10541efe70f51XXX",oauth_signature="n868PkpJIF180bL%2FDhf2yWXXX",oauth_signature_method="HMAC-SHA1",oauth_timestamp="1653340748",oauth_token="XXX",oauth_version="1.0"
code: 200 body: oauth_token=XXX&oauth_token_secret=XXX&oauth_callback_confirmed=true
redirecting to https://api.twitter.com/oauth/authenticate?oauth_token=XXX
starting step 3
final result: 
200 oauth_token=XXX&oauth_token_secret=XXX&user_id=XXX&screen_name=YOURTWITTERNAME

Note: In some places (including the twitter docs) it says to leave the token secret off the end of your key and just have a hanging '&' character there for the initial authentication requests. Do not do that, it will not work.

Resources:
https://developer.twitter.com/en/docs/authentication/guides/log-in-with-twitter#obtain-a-request-token