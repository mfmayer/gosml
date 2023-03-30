# gosml

<!--[![Build Status](https://travis-ci.org/andig/gosml.svg?branch=master)](https://travis-ci.org/andig/gosml)-->

This repository is a fork of [github.com/andig/gosml](https://github.com/andig/gosml), which is a Go port of [github.com/volkszaehler/libsml](https://github.com/volkszaehler/libsml).

I adapted it to my needs and to make it more usable in my eyes.

## Usage

To include in your code use:

````
import (
	"github.com/mfmayer/gosml"
)
````

## Example

For an example see [examples/emmon](https://github.com/mfmayer/gosml/blob/master/examples/emmon/readme.md) and the [libsml](https://github.com/volkszaehler/libsml) documentation.

## Status

The implementation of this port and fork is not complete and has not been extensively tested. It's main intension is to use it for decoding list entries with well known obis codes like for electricity meter readings.

Test binaries and SML files from real world meters can be found here: <https://github.com/devZer0/libsml-testing>