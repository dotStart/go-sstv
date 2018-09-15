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

type RobotMode uint8

const (
  Robot36 RobotMode = 8
  Robot72 RobotMode = 12
)

const (
  robotLineFrequency = 1200
  robotLineLength    = 9

  robotSyncFrequency = 1500
  robotSyncLength    = 3

  robotPorchFrequency = 1900
  robotPorchLength    = 1.5

  robotEvenSeparatorFrequency = 1500
  robotOddSeparatorFrequency  = 2300
  robotSeparatorLength        = 4.5
)

const (
  robot36YLength = 0.275
  robot36Length  = 0.1375

  robot72YLength = 0.43125
  robot72Length  = 0.215625
)

// provides a Scottie implementation as designed by Eddie Murphy
//
// this implementation encodes RGB images into 240 lines
type robotEncoder struct {
  mode   RobotMode
  format *audio.Format
}

// creates a new Martin compatible image encoder
func NewRobot(mode RobotMode, format *audio.Format) Encoder {
  return &robotEncoder{
    mode:   mode,
    format: format,
  }
}

func (enc *robotEncoder) Vis() uint8 {
  return uint8(enc.mode)
}

func (enc *robotEncoder) Resolution() image.Rectangle {
  return image.Rect(0, 0, 320, 240)
}

func (enc *robotEncoder) Encode(img image.Image) *audio.FloatBuffer {
  wr := newWriter(enc.format)
  wr.writeHeader()
  wr.writeVis(uint8(enc.mode))

  // different from other encoders, Robot provides two completely different encoding types which
  // share little to nothing and thus we 'll need two separate encoding methods
  switch enc.mode {
  case Robot36:
    enc.encode36(wr, img)
  case Robot72:
    enc.encode72(wr, img)
  default:
    panic(errors.New("illegal encoding mode"))
  }

  return wr.buf
}

func (enc *robotEncoder) encode36(wr *audioWriter, img image.Image) {
  size := img.Bounds().Size()
  ymap, umap, vmap := enc.generateYUVMap(img)
  for y := 0; y < size.Y; y++ {
    wr.write(robotLineFrequency, robotLineLength)
    wr.write(robotSyncFrequency, robotSyncLength)
    even := y%2 == 0

    for i := 0; i < 2; i++ {
      for x := 0; x < size.X; x++ {
        mapIndex := y*size.X + x

        var val float64
        l := robot36Length
        if i == 0 {
          val = float64(ymap[mapIndex]) / 255
          l = robot36YLength
        } else if even {
          val = float64(umap[mapIndex]) / 255
        } else {
          val = float64(vmap[mapIndex]) / 255
        }

        wr.writeValue(val, l)
      }

      if i == 0 {
        if even {
          wr.write(robotEvenSeparatorFrequency, robotSeparatorLength)
        } else {
          wr.write(robotOddSeparatorFrequency, robotSeparatorLength)
        }

        wr.write(robotPorchFrequency, robotPorchLength)
      }
    }
  }
}

func (enc *robotEncoder) encode72(wr *audioWriter, img image.Image) {
  size := img.Bounds().Size()
  for y := 0; y < size.Y; y++ {
    wr.write(robotLineFrequency, robotLineLength)
    wr.write(robotSyncFrequency, robotSyncLength)

    for i := 0; i < 3; i++ {
      for x := 0; x < size.X; x++ {
        y, u, v := convertYUV(img.At(x, y))

        var val float64
        l := robot72Length
        if i == 0 {
          val = float64(y) / 255
          l = robot72YLength
        } else if i == 1 {
          val = float64(u) / 255
        } else {
          val = float64(v) / 255
        }

        wr.writeValue(val, l)
      }

      if i != 2 {
        if i % 2 == 0 {
          wr.write(robotEvenSeparatorFrequency, robotSeparatorLength)
          wr.write(robotPorchFrequency, robotPorchLength)
        } else {
          wr.write(robotOddSeparatorFrequency, robotSeparatorLength)
          wr.write(robotSyncFrequency, robotPorchLength)
        }
      }
    }
  }
}

func (enc *robotEncoder) generateYUVMap(img image.Image) ([]byte, []byte, []byte) {
  size := img.Bounds().Size()
  y := make([]byte, size.X*size.Y)
  u := make([]byte, size.X*size.Y)
  v := make([]byte, size.X*size.Y)

  for ycord := 0; ycord < size.Y; ycord++ {
    for xcord := 0; xcord < size.X; xcord++ {
      yv, uv00, vv00 := convertYUV(img.At(xcord, ycord))
      _, uv01, vv01 := convertYUV(img.At(xcord, ycord+1))
      _, uv10, vv10 := convertYUV(img.At(xcord+1, ycord))
      _, uv11, vv11 := convertYUV(img.At(xcord+1, ycord+1))

      i := ycord*size.X + xcord
      y[i] = yv
      u[i] = byte((int(uv00) + int(uv01) + int(uv10) + int(uv11)) / 4)
      v[i] = byte((int(vv00) + int(vv01) + int(vv10) + int(vv11)) / 4)
    }
  }

  return y, u, v
}
