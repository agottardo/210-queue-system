
# CPSC 210 Office Hours Queue System

This application lets students taking CPSC 210 at UBC sign up for office hours via a convenient web interface.

## Running Locally

Make sure you have [Go](http://golang.org/doc/install) and the [Heroku Toolbelt](https://toolbelt.heroku.com/) installed.

```sh
$ go get -u github.com/agottardo/210-queue-system
$ cd $GOPATH/src/github.com/agottardo/210-queue-system
$ heroku local
```

The app will be running on [localhost:5000](http://localhost:5000/).

## Deploying to Heroku

```sh
$ heroku create
$ git push heroku master
$ heroku open
```

or

[![Deploy](https://www.herokucdn.com/deploy/button.png)](https://heroku.com/deploy)
