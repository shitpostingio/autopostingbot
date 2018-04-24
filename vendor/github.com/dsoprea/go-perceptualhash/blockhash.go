package blockhash

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"math/big"
	"sort"
	"strconv"

	"github.com/dsoprea/go-logging"
)

type Blockhash struct {
	image        image.Image
	hashbits     int
	toColor      *color.Model
	hasAlpha     bool
	hexdigest    string
	isOpaqueable bool
}

// opaqueableModel automatically fulfilled by existing Go types.
type opaqueableModel interface {
	Opaque() bool
}

func NewBlockhash(image image.Image, hashbits int) *Blockhash {
	// Only images that support alpha are explicitly aware of opaqueness.
	_, isOpaqueable := image.(opaqueableModel)

	// If the bits aren't aligned, the digest won't make sense as a hex string.
	if (hashbits % 4) != 0 {
		log.Panicf("Bits must be a multiple of four: (%d)", hashbits)
	}

	return &Blockhash{
		image:        image,
		hashbits:     hashbits,
		isOpaqueable: isOpaqueable,
	}
}

func (bh *Blockhash) totalValue(p color.Color) (value uint32) {
	defer func() {
		if state := recover(); state != nil {
			log.Panic(state.(error))
		}
	}()

	// The RGBA() will return the alpha-multiplied values but the fields will
	// still be in their premultiplied state.
	if bh.image.ColorModel() != color.RGBAModel {
		p = color.RGBAModel.Convert(p)
	}

	c2 := p.(color.RGBA)

	if bh.isOpaqueable == true && c2.A == 0 {
		return 765
	}

	return uint32(c2.R) + uint32(c2.G) + uint32(c2.B)
}

func (bh *Blockhash) totalValueAt(x, y int) (value uint32) {
	defer func() {
		if state := recover(); state != nil {
			log.Panic(state.(error))
		}
	}()

	p := bh.image.At(x, y)

	return bh.totalValue(p)
}

func (bh *Blockhash) median(data []float64) float64 {
	defer func() {
		if state := recover(); state != nil {
			log.Panic(state.(error))
		}
	}()

	copied := make([]float64, len(data))
	copy(copied, data)
	sort.Float64s(copied)

	len_ := len(copied)
	if len(copied)%2 == 0 {
		v := (copied[len_/2-1] + copied[len_/2]) / 2.0

		return v
	} else {
		v := copied[len_/2]

		return v
	}
}

func (bh *Blockhash) bitsToHex(bitString []int) string {
	defer func() {
		if state := recover(); state != nil {
			log.Panic(state.(error))
		}
	}()

	s := make([]byte, len(bitString))

	for i, d := range bitString {
		if d == 0 {
			s[i] = '0'
		} else if d == 1 {
			s[i] = '1'
		} else {
			log.Panicf("invalid bit value (%d) at offset (%d)", d, i)
		}
	}

	b := new(big.Int)
	b.SetString(string(s), 2)

	width := int(math.Pow(float64(bh.hashbits), 2.0) / 4.0)
	encoded := fmt.Sprintf("%0"+strconv.Itoa(width)+"x", b)

	return encoded
}

func (bh *Blockhash) translateBlocksToBits(blocksInline []float64, pixelsPerBlock float64) (results []int) {
	defer func() {
		if state := recover(); state != nil {
			log.Panic(state.(error))
		}
	}()

	blocks := make([]int, len(blocksInline))
	halfBlockValue := pixelsPerBlock * 256.0 * 3.0 / 2.0

	bandsize := int(math.Floor(float64(len(blocksInline)) / 4.0))

	for i := 0; i < 4; i++ {
		m := bh.median(blocksInline[i*bandsize : (i+1)*bandsize])

		for j := i * bandsize; j < (i+1)*bandsize; j++ {
			v := blocksInline[j]

			// TODO(dustin): Use epsilon.
			if v > m || (math.Abs(v-m) < 1 && m > halfBlockValue) {
				blocks[j] = 1
			} else {
				blocks[j] = 0
			}
		}
	}

	return blocks
}

func (bh *Blockhash) size() (width int, height int) {
	r := bh.image.Bounds()

	width = r.Max.X
	height = r.Max.Y

	return width, height
}

func (bh *Blockhash) getBlocks() []float64 {
	width, height := bh.size()

	isEvenX := (width % bh.hashbits) == 0
	isEvenY := (height % bh.hashbits) == 0

	blocks := make([][]float64, bh.hashbits)

	for i := 0; i < bh.hashbits; i++ {
		blocks[i] = make([]float64, bh.hashbits)
	}

	blockWidth := float64(width) / float64(bh.hashbits)
	blockHeight := float64(height) / float64(bh.hashbits)

	for y := 0; y < height; y++ {
		var weightTop, weightBottom, weightLeft, weightRight float64
		var blockTop, blockBottom, blockLeft, blockRight int

		if isEvenY {
			blockTop = int(math.Floor(float64(y) / blockHeight))
			blockBottom = blockTop

			weightTop = 1.0
			weightBottom = 0.0
		} else {
			yMod := math.Mod((float64(y) + 1.0), blockHeight)
			yInt, yFrac := math.Modf(yMod)

			weightTop = (1.0 - yFrac)
			weightBottom = yFrac

			// y_int will be 0 on bottom/right borders and on block boundaries
			if yInt > 0.0 || (y+1) == height {
				blockTop = int(math.Floor(float64(y) / blockHeight))
				blockBottom = blockTop
			} else {
				blockTop = int(math.Floor(float64(y) / blockHeight))
				blockBottom = int(math.Ceil(float64(y) / blockHeight))
			}

		}

		for x := 0; x < width; x++ {
			value := bh.totalValueAt(x, y)

			if isEvenX {
				blockRight = int(math.Floor(float64(x) / blockWidth))
				blockLeft = blockRight

				weightLeft = 1.0
				weightRight = 0.0
			} else {
				xMod := math.Mod((float64(x) + 1.0), blockWidth)
				xInt, xFrac := math.Modf(xMod)

				weightLeft = (1.0 - xFrac)
				weightRight = (xFrac)

				if xInt > 0.0 || (x+1) == width {
					blockRight = int(math.Floor(float64(x) / blockWidth))
					blockLeft = blockRight
				} else {
					blockLeft = int(math.Floor(float64(x) / blockWidth))
					blockRight = int(math.Ceil(float64(x) / blockWidth))
				}
			}

			blocks[blockTop][blockLeft] += float64(value) * weightTop * weightLeft
			blocks[blockTop][blockRight] += float64(value) * weightTop * weightRight
			blocks[blockBottom][blockLeft] += float64(value) * weightBottom * weightLeft
			blocks[blockBottom][blockRight] += float64(value) * weightBottom * weightRight
		}
	}

	blocksInline := make([]float64, bh.hashbits*bh.hashbits)
	i := 0
	for y := 0; y < bh.hashbits; y++ {
		for x := 0; x < bh.hashbits; x++ {
			blocksInline[i] = blocks[y][x]
			i++
		}
	}

	return blocksInline
}

func (bh *Blockhash) process() (err error) {
	defer func() {
		if state := recover(); state != nil {
			err = log.Wrap(state.(error))
		}
	}()

	if bh.hexdigest != "" {
		return nil
	}

	blocks := bh.getBlocks()

	width, height := bh.size()

	blockWidth := float64(width) / float64(bh.hashbits)
	blockHeight := float64(height) / float64(bh.hashbits)

	digest := bh.translateBlocksToBits(blocks, blockWidth*blockHeight)
	bh.hexdigest = bh.bitsToHex(digest)

	return nil
}

func (bh *Blockhash) Hexdigest() string {
	defer func() {
		if state := recover(); state != nil {
			err := log.Wrap(state.(error))
			log.PanicIf(err)
		}
	}()

	err := bh.process()
	log.PanicIf(err)

	return bh.hexdigest
}
