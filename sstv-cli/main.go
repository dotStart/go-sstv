/*
 * Copyright 2018 Johannes Donath <johannesd@torchmind.com>
 * and other copyright owners as documented in the project's IP log.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package main

import (
  "flag"
  "fmt"
  "github.com/dotStart/go-sstv"
  "github.com/go-audio/audio"
  "github.com/go-audio/wav"
  "github.com/nfnt/resize"
  "image"
  _ "image/jpeg"
  _ "image/png"
  "os"
)

var audioFormat = audio.FormatMono44100

func main() {
  var flagHelp bool
  var flagSampleRate int
  var flagMartin1, flagMartin2 bool
  var flagPasokon3, flagPasokon5, flagPasokon7 bool
  var flagRobot36, flagRobot72 bool
  var flagScottie1, flagScottie2, flagScottieDx bool
  var flagWrasseSC2180 bool

  flag.BoolVar(&flagHelp, "help", false, "displays this help message")
  flag.IntVar(&flagSampleRate, "sample-rate", 44100, "specifies the sample rate (defaults to 19200 Hz)")
  flag.BoolVar(&flagMartin1, "m1", false, "uses Martin encoding in M1 mode")
  flag.BoolVar(&flagMartin2, "m2", false, "uses Martin encoding in M2 mode")
  flag.BoolVar(&flagPasokon3, "p3", false, "uses Pasokon (\"P\") in P3 mode")
  flag.BoolVar(&flagPasokon5, "p5", false, "uses Pasokon (\"P\") in P3 mode")
  flag.BoolVar(&flagPasokon7, "p7", false, "uses Pasokon (\"P\") in P3 mode")
  flag.BoolVar(&flagRobot36, "r36", false, "uses Robot encoding in 36 mode")
  flag.BoolVar(&flagRobot72, "r72", false, "uses Robot encoding in 72 mode")
  flag.BoolVar(&flagScottie1, "s1", false, "uses Scottie encoding in S1 mode")
  flag.BoolVar(&flagScottie2, "s2", false, "uses Scottie encoding in S2 mode")
  flag.BoolVar(&flagScottieDx, "sdx", false, "uses Scottie encoding in DX mode")
  flag.BoolVar(&flagWrasseSC2180, "wrsc2-180", false, "uses Wrasse encoding in SC2-180 mode")

  flag.Parse()

  if flagHelp {
    printHelp()
    return
  }

  if flag.NArg() != 2 {
    printHelp()
    os.Exit(1)
  }

  format := &audio.Format{
    NumChannels: 1,
    SampleRate:  flagSampleRate,
  }

  var tv sstv.Encoder
  if flagMartin1 || flagMartin2 {
    mode := sstv.Martin1
    if flagMartin2 {
      mode = sstv.Martin2
    }

    tv = sstv.NewMartin(mode, format)
  } else if flagPasokon3 || flagPasokon5 || flagPasokon7 {
    mode := sstv.Pasokon3
    if flagPasokon5 {
      mode = sstv.Pasokon5
    } else if flagPasokon7 {
      mode = sstv.Pasokon7
    }

    tv = sstv.NewPasokon(mode, format)
  } else if flagRobot36 || flagRobot72 {
    mode := sstv.Robot36
    if flagRobot72 {
      mode = sstv.Robot72
    }

    tv = sstv.NewRobot(mode, format)
  } else if flagScottie1 || flagScottie2 || flagScottieDx {
    mode := sstv.Scottie1
    if flagScottie2 {
      mode = sstv.Scottie2
    } else if flagScottieDx {
      mode = sstv.ScottieDx
    }

    tv = sstv.NewScottie(mode, format)
  } else if flagWrasseSC2180 {
    mode := sstv.WrasseSC2180

    tv = sstv.NewWrasse(mode, format)
  }

  fmt.Printf("==> using VIS 0x%02x\n", tv.Vis())

  var img image.Image
  fmt.Print("loading file ... ")
  if f, err := os.OpenFile(flag.Arg(0), os.O_RDONLY, os.ModePerm); err == nil {
    img, _, err = image.Decode(f)
    if err != nil {
      fmt.Printf("failed: %s", err)
      os.Exit(2)
    }
    fmt.Printf("ok (%d x %d)\n", img.Bounds().Size().X, img.Bounds().Size().Y)
  } else {
    fmt.Printf("failed: %s", err)
    os.Exit(2)
  }

  fmt.Print("resizing ... ")
  targetRes := tv.Resolution().Size()
  currentRes := img.Bounds().Size()
  if currentRes.X != targetRes.X || currentRes.Y != targetRes.Y {
    img = resize.Resize(uint(targetRes.X), uint(targetRes.Y), img, resize.Lanczos3)
    fmt.Print("ok\n")
  } else {
    fmt.Print("skipped\n")
  }

  fmt.Print("generating ... ")
  buf := tv.Encode(img)
  fmt.Print("ok\n")

  fmt.Print("encoding ... ")
  if wr, err := os.OpenFile(flag.Arg(1), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm); err == nil {
    defer wr.Close()

    enc := wav.NewEncoder(wr, flagSampleRate, sstv.BitDepth, 1, 1)
    defer enc.Close()

    if err = enc.Write(buf.AsIntBuffer()); err != nil {
      fmt.Printf("failed: %s", err)
      os.Exit(2)
    }

    fmt.Print("ok\n")
  } else {
    fmt.Printf("failed: %s\n", err)
    os.Exit(2)
  }
}

// writes the command line help to stdout
func printHelp() {
  fmt.Printf("Usage: %s [flags] <in> <out>\n\n", os.Args[0])
  flag.PrintDefaults()
}
