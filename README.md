go SSTV
=======

Encoding of images into audio using the SSTV standard (and its most popular encoding modes such as
Martin, Robot and Scottie).

Usage
-----

```go
tv := sstv.NewMartin(sstv.Martin1, &audio.Format{
  SampleRate: 41000,
  NumChannels: 1,
})

// Or for the other modes:
// sstv.NewPasokon(sstv.Pasokon3, format)
// sstv.NewRobot(sstv.Robot36, format)
// sstv.NewScottie(sstv.Scottie1, format)
// sstv.NewWrasse(sstv.WrasseSC2180, format)
// For a full list of mode constants, refer to the package documentation
```

For a full list of mode constants, refer to the [package documentation](https://godoc.org/github.com/dotStart/go-sstv)

Command Line Interface
----------------------

The package includes a command line interface for testing purposes:

```sh
# Generate a Martin1 encoded image:
$ sstv-cli -m1 -sample-rate=41000 input.png output.wav

# Generate a Scottie1 encoded image:
$ sstv-cli -s1 -sample-rate=41000 input.png output.wav

# Display all modes:
$ sstv-cli -help
```

You may install the CLI command via `go get -u github.com/dotStart/go-sstv/...`

License
-------

```
Copyright [year] [name] <[email]>
and other copyright owners as documented in the project's IP log.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```
