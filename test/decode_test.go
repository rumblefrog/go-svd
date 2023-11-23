package svd_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
	"github.com/rumblefrog/go-svd/svd"
)

func TestDanSamples(t *testing.T) {
	count := 97

	payloads := make([][]byte, 0, count)

	for i := 0; i < count; i += 1 {
		b, err := os.ReadFile(fmt.Sprintf("dan_samples/dan-data-%d.bin", i))

		if err != nil {
			t.Fatal(err)
		}

		payloads = append(payloads, b)
	}

	decoder, err := svd.NewOpusDecoder(24000, 1)

	if err != nil {
		t.Fatal(err)
	}

	o := make([]int, 0, 1024)

	for _, payload := range payloads {
		c, err := svd.DecodeChunk(payload)

		if err != nil {
			t.Fatal(err)
		}

		// Not silent frame
		if len(c.Data) > 0 {
			pcm, err := decoder.Decode(c.Data)

			if err != nil {
				t.Fatal(err)
			}

			converted := make([]int, len(pcm))
			for i, v := range pcm {
				// Float32 buffer implementation is wrong in go-audio, so we have to convert to int before encoding
				converted[i] = int(v * 2147483647)
			}

			o = append(o, converted...)
		}
	}

	outFile, err := os.Create("out.wav")

	if err != nil {
		t.Fatal(err)
	}
	defer outFile.Close()

	enc := wav.NewEncoder(outFile, 24000, 32, 1, 1)

	buf := &audio.IntBuffer{
		Data: o,
		Format: &audio.Format{
			SampleRate:  24000,
			NumChannels: 1,
		},
	}

	if err := enc.Write(buf); err != nil {
		t.Fatal(err)
	}

	enc.Close()
}
