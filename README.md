# qt-validator-client-demo
A primitive qt client that can exercise the Nym validators

**DEPRECATED: this code works with the old Go-based Nym validators. We are currently re-writing the validators in Rust. We leave this code up in case anyone wants to see Coconut client code in Go.**

## Building

Unfortunately at this point of time, building the qt demo is not as straightforward we would want it to be. It will be improved in the future. The current procedure is as follows:

1. Install the Qt bindings by following the instructions here on [https://github.com/therecipe/qt/wiki/Installation-on-Linux](https://github.com/therecipe/qt/wiki/Installation-on-Linux)
2. Make sure the bindings are actually in your `$GOPATH`. If not, run `go get github.com/therecipe/qt`
3. Clone this repository into your `$GOPATH`
4. Run `dep ensure` inside the repository
5. IMPORTANT: Remove `./vendor/github.com/therecipe/qt`. We want to be using this dependency from our `$GOPATH` instead. If it were in the vendor directory, it wouldn't work. Refer to [https://github.com/therecipe/qt/issues/615](https://github.com/therecipe/qt/issues/615) for 'more' details
6. Run `make build_gui` or make `build_release_gui` depending on your intent. 
