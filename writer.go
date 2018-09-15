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
)

const BitDepth = 16

const headerFrequency = 1900
const headerLength = 300
const headerPauseFrequency = 1200
const headerPauseLength = 10

const headerVisFrequency = 1200
const bitLength = 30
const trueFrequency = 1100
const falseFrequency = 1300

const blackFrequency = 1500
const whiteFrequency = 2300

type audioWriter struct {
  gen *oscillator
  buf *audio.FloatBuffer
}

func newWriter(format *audio.Format) *audioWriter {
  return &audioWriter{
    gen: newOscillator(format.SampleRate, float64(audio.IntMaxSignedValue(BitDepth))),
    buf: &audio.FloatBuffer{
      Format: format,
      Data:   make([]float64, 0),
    },
  }
}

// appends a signal with the given frequency and length to the buffer
func (wr *audioWriter) write(freq float64, length float64) {
  wr.buf.Data = append(wr.buf.Data, wr.gen.signal(freq, length)...)
}

// writes a boolean bit to the buffer (in the VIS code format)
func (wr *audioWriter) writeBit(val bool) {
  if val {
    wr.write(trueFrequency, bitLength)
  } else {
    wr.write(falseFrequency, bitLength)
  }
}

func (wr *audioWriter) writeHeader() {
  wr.write(headerFrequency, headerLength)
  wr.write(headerPauseFrequency, headerPauseLength)
  wr.write(headerFrequency, headerLength)
}

func (wr *audioWriter) writeVis(val uint8) {
  wr.write(headerVisFrequency, bitLength)

  p := parity(val)
  for i := 0; i < 7; i++ {
    wr.writeBit(val&0x1 == 0x1)
    val >>= 1
  }
  wr.writeBit(p)

  wr.write(headerVisFrequency, bitLength)
}

func (wr *audioWriter) writeValue(val float64, length float64) {
  wr.write(val*(whiteFrequency-blackFrequency)+blackFrequency, length)
}

// computes the VIS parity
func parity(val uint8) bool {
  var p = true
  for val != 0 {
    p = !p
    val >>= 1
  }
  return p
}
