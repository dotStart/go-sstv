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

type MartinMode uint8

const (
  Martin1 MartinMode = 44
  Martin2 MartinMode = 40
)

const (
  martin1PulseLength = .4576
  martin2PulseLength = .2288

  martinLineFrequency = 1200
  martinLineLength    = 4.862

  martinSeparatorLength    = .572
  martinSeparatorFrequency = 1500
)

// provides a Martin implementation as designed by Martin Emmerson
//
// this implementation encodes RGB images into 240 lines within either 114 or 58 seconds
type martinEncoder struct {
  mode   MartinMode
  format *audio.Format
}

// creates a new Martin compatible image encoder
func NewMartin(mode MartinMode, format *audio.Format) Encoder {
  return &martinEncoder{
    mode:   mode,
    format: format,
  }
}

func (enc *martinEncoder) Vis() uint8 {
  return uint8(enc.mode)
}

func (enc *martinEncoder) Resolution() image.Rectangle {
  return image.Rect(0, 0, 320, 256)
}

func (enc *martinEncoder) Encode(img image.Image) *audio.FloatBuffer {
  var pulseLength float64
  switch enc.mode {
  case Martin1:
    pulseLength = martin1PulseLength
  case Martin2:
    pulseLength = martin2PulseLength
  default:
    panic(errors.New("illegal encoding mode"))
  }

  wr := newWriter(enc.format)
  wr.writeHeader()
  wr.writeVis(uint8(enc.mode))

  size := img.Bounds().Size()
  for y := 0; y < size.Y; y++ {
    wr.write(martinLineFrequency, martinLineLength)

    for i := 0; i < 3; i++ {
      wr.write(martinSeparatorFrequency, martinSeparatorLength)

      for x := 0; x < size.X; x++ {
        r, g, b := convertRGB(img.At(x, y))

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

      wr.write(martinSeparatorFrequency, martinSeparatorLength)
    }
  }

  return wr.buf
}
