# tracker-cli [![Build Status](https://travis-ci.org/kkallday/tracker-cli.svg?branch=master)](https://travis-ci.org/kkallday/tracker-cli)

A simple command-line tool to display stories from your Pivotal Tracker project in 
the terminal. Less "Command+Tab"-ing to-and-from your browser.

tracker-cli works with public and private Tracker projects.

## Requirements

- Go (v1.6+)

## Setup

`tracker-cli` needs the project ID of the Tracker Project you would like to pull 
stories from. This is provided via the `PROJECT_ID` environment variable. If your
project is private, you will need to provide an access token via the `TOKEN` 
environment variable. Information on how to obtain an access token can be 
[found here] (https://www.pivotaltracker.com/help/api#Getting_Started).

It is recommended that you add the `PROJECT_ID` and `TOKEN` variables to your 
`~/.bash_profile`

## Usage

To display all in-flight stories run:

```
PROJECT_ID=12345 TOKEN="XYZ" tracker-cli
```

NOTE: The environment variable `TOKEN` is only required for non-public tracker projects.


## Testing

- To run tests you will need [ginkgo] (https://onsi.github.io/ginkgo) and [gomega] (https://onsi.github.io/gomega)
- Run `ginkgo -r` or `./bin/test` to run all tests
