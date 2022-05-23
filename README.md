So I had started working with the twitter api for a small project. The most annoying part was getting setup and authorizing the first request. So I thought I would writeup a simple version of it here in go using no dependencies and explain what I did and why (and maybe where the docs are a bit misleading).

You could use this like a simple library but I won't make any guarantee as to it's reliability.

https://developer.twitter.com/en/docs/authentication/guides/log-in-with-twitter#obtain-a-request-token