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

type PasokonMode uint8

const (
  Pasokon3 PasokonMode = 113
  Pasokon5 PasokonMode = 114
  Pasokon7 PasokonMode = 115
)

const (
  pasokonLineFrequency = 1200
  pasokon3LineLength   = 5.208
  pasokon5LineLength   = 7.813
  pasokon7LineLength   = 10.417

  pasokonSyncFrequency = 1500
  pasokon3SyncLength   = 1.042
  pasokon5SyncLength   = 1.563
  pasokon7SyncLength   = 2.083

  pasokon3PulseLength = .2083
  pasokon5PulseLength = .3125
  pasokon7PulseLength = .4167
)

// provides a Pasokon ("P") implementation
//
// this implementation encodes RGB images into 496 lines
type pasokonEncoder struct {
  mode   PasokonMode
  format *audio.Format
}

// creates a new Pasokon compatible image encoder
func NewPasokon(mode PasokonMode, format *audio.Format) Encoder {
  return &pasokonEncoder{
    mode:   mode,
    format: format,
  }
}

func (enc *pasokonEncoder) Vis() uint8 {
  return uint8(enc.mode)
}

func (enc *pasokonEncoder) Resolution() image.Rectangle {
  return image.Rect(0, 0, 640, 496)
}

func (enc *pasokonEncoder) Encode(img image.Image) *audio.FloatBuffer {
  wr := newWriter(enc.format)
  wr.writeHeader()
  wr.writeVis(uint8(enc.mode))

  var lineLength, syncLength, pulseLength float64
  switch enc.mode {
  case Pasokon3:
    lineLength = pasokon3LineLength
    syncLength = pasokon3SyncLength
    pulseLength = pasokon3PulseLength
  case Pasokon5:
    lineLength = pasokon5LineLength
    syncLength = pasokon5SyncLength
    pulseLength = pasokon5PulseLength
  case Pasokon7:
    lineLength = pasokon7LineLength
    syncLength = pasokon7SyncLength
    pulseLength = pasokon7PulseLength
  }

  size := img.Bounds().Size()
  for y := 0; y < size.Y; y++ {
    wr.write(pasokonLineFrequency, lineLength)

    for i := 0; i < 3; i++ {
      wr.write(pasokonSyncFrequency, syncLength)

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
    }

    wr.write(pasokonSyncFrequency, syncLength)
  }

  return wr.buf
}
