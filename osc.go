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

import "math"

// provides a stateful sine wave oscillator
type oscillator struct {
  sampleRate int
  amplitude  float64
  phase      float64
}

// creates a new oscillator with the indicated sample rate
func newOscillator(sampleRate int, amplitude float64) *oscillator {
  return &oscillator{
    sampleRate: sampleRate,
    amplitude:  amplitude,
    phase:      0,
  }
}

// generates the value for a single sample of the indicated frequency
func (osc *oscillator) sample(frequency float64) float64 {
  osc.phase += frequency * 2 * math.Pi / float64(osc.sampleRate)
  return math.Sin(osc.phase) * osc.amplitude
}

// generates a signal of the indicated length (in milliseconds)
func (osc *oscillator) signal(frequency float64, length float64) []float64 {
  samples := int(math.Round(length / 1000 * float64(osc.sampleRate)))
  values := make([]float64, samples)

  for i := 0; i < samples; i++ {
    values[i] = osc.sample(frequency)
  }

  return values
}
