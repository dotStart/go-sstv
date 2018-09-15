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
  "github.com/go-audio/audio"
  "image"
)

type WrasseMode uint8

const (
  WrasseSC2180 WrasseMode = 55
)

const wrassePulseLength = .7344

const (
  wrasseLineFrequency = 1200
  wrasseLineLength    = 5.5225

  wrasseSyncFrequency = 1500
  wrasseSyncLength    = 0.5
)

// provides a Wrasse implementation
//
// this implementation encodes RGB images into 256 lines
type wrasseEncoder struct {
  mode   WrasseMode
  format *audio.Format
}

// creates a new Wrasse compatible image encoder
func NewWrasse(mode WrasseMode, format *audio.Format) Encoder {
  return &wrasseEncoder{
    mode:   mode,
    format: format,
  }
}

func (enc *wrasseEncoder) Vis() uint8 {
  return uint8(enc.mode)
}

func (enc *wrasseEncoder) Resolution() image.Rectangle {
  return image.Rect(0, 0, 320, 256)
}

func (enc *wrasseEncoder) Encode(img image.Image) *audio.FloatBuffer {
  wr := newWriter(enc.format)
  wr.writeHeader()
  wr.writeVis(uint8(enc.mode))

  size := img.Bounds().Size()
  for y := 0; y < size.Y; y++ {
    wr.write(wrasseLineFrequency, wrasseLineLength)
    wr.write(wrasseSyncFrequency, wrasseSyncLength)

    for i := 0; i < 3; i++ {
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
        wr.writeValue(val, wrassePulseLength)
      }
    }
  }

  return wr.buf
}
