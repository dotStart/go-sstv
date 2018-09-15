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
package sstv

import (
  "errors"
  "github.com/go-audio/audio"
  "image"
)

type ScottieMode uint8

const (
  Scottie1  ScottieMode = 60
  Scottie2  ScottieMode = 56
  ScottieDx ScottieMode = 76
)

const scottie1PulseLength = .4320
const scottie2PulseLength = .2752
const scottieDxPulseLength = 1.08

const scottieSyncLength = 9
const scottieSyncFrequency = 1200

const scottySeparatorLength = 1.5
const scottySeparatorFrequency = 1500

// provides a Scottie implementation as designed by Eddie Murphy
//
// this implementation encodes RGB images into 240 lines
type scottieEncoder struct {
  mode   ScottieMode
  format *audio.Format
}

// creates a new Scottie compatible image encoder
func NewScottie(mode ScottieMode, format *audio.Format) Encoder {
  return &scottieEncoder{
    mode:   mode,
    format: format,
  }
}

func (enc *scottieEncoder) Vis() uint8 {
  return uint8(enc.mode)
}

func (enc *scottieEncoder) Resolution() image.Rectangle {
  return image.Rect(0, 0, 320, 256)
}

func (enc *scottieEncoder) Encode(image image.Image) *audio.FloatBuffer {
  var pulseLength float64
  switch enc.mode {
  case Scottie1:
    pulseLength = scottie1PulseLength
  case Scottie2:
    pulseLength = scottie2PulseLength
  case ScottieDx:
    pulseLength = scottieDxPulseLength
  default:
    panic(errors.New("illegal encoding mode"))
  }

  wr := newWriter(enc.format)
  wr.writeHeader()
  wr.writeVis(uint8(enc.mode))

  size := image.Bounds().Size()
  for y := 0; y < size.Y; y++ {
    if y == 0 {
      wr.write(scottieSyncFrequency, scottieSyncLength)
    }

    for i := 0; i < 3; i++ {
      wr.write(scottySeparatorFrequency, scottySeparatorLength)

      for x := 0; x < size.X; x++ {
        r, g, b := convertRGB(image.At(x, y))

        var val float64
        switch i {
        case 0:
          val = g
        case 1:
          val = b
        case 2:
          val = r
        }
        wr.writeValue(val, pulseLength)
      }

      if i == 1 {
        wr.write(scottieSyncFrequency, scottieSyncLength)
      }
    }
  }

  return wr.buf
}
