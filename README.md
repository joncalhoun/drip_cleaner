# Drip Cleaner

A tool to scan your Drip mailing list for disposable email addresses and delete them.

## Why?

I have a few freebies and courses that require people to subscribe to my mailing list to get access. A few examples include [Gophercises](https://gophercises.com/), Web Development with Go [course samples](https://www.usegolang.com/#chapters), and [a guide to learning Go](https://www.calhoun.io/6-tips-for-using-strings-in-go/#subscribe).

When a user signs up to my mailing list I respect their privacy. I don't sell their email address, I limit the emails I send to them to only ones I think they will find useful, and if they opt to unsubscribe at any time I never email them again or shadily add them to a new campaign. In short, I treat their email address like I hope others treat my own.

Unfortunately there are many bad actors out there who don't deserve our trust and they ruin things for everyone. As a result, many users will want a freebie, but won't be willing to part with their real email address. Instead they will use a disposable email address.

I don't block these email addresses because I want those users to still get access to useful material, but it does pose a problem - once they sign up and get their freebie I am left with a useless email in my mailing list and given that Drip's pricing model is based on the number of subscribers, this inevitably increases my costs.

This program is a remedy to that problem. Users are still welcome to sign up with a disposable email address, but every month or so I prune my mailing list of all disposable email addresses.

## How to use it

TODO(jon): Fill this in

## Disclaimer

The Drip client is not commented, tested, or moved into its own package. It also only grabs data that I care about form JSON responses. In short, **the code is ugly and only does what I need it to do**. I'm well aware, but it gets the job done so I don't intend to change much until I need more functionality.

I have tested this code by running it against my own drip account, but you might want to review it before trusting it freely. If you have questions about what anything is doing just email me - <jon@calhoun.io>

