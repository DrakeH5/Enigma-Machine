package main

import (
	"fmt"
	"image/color"
	_ "image/png"
	"log"
	"math"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	//"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text"

	"github.com/hajimehoshi/ebiten/v2/inpututil"
	//"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth  = 640
	screenHeight = 480
)

var keys = [26]string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}

var (
	mplusNormalFont font.Face
	mplusBigFont    font.Face
)

var rotorImg *ebiten.Image
var topbgImg *ebiten.Image
var reflectortopImg *ebiten.Image
var emptyRotorSlotImg *ebiten.Image

func init() {
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 72
	mplusNormalFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    24,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	mplusBigFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    32,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}

	rotorImg, _, err = ebitenutil.NewImageFromFile("rotors.png")
	topbgImg, _, err = ebitenutil.NewImageFromFile("topbg.png")
	reflectortopImg, _, err = ebitenutil.NewImageFromFile("reflectortop.png")
	emptyRotorSlotImg, _, err = ebitenutil.NewImageFromFile("emptyrotorSlot.png")
	for i := 0; i < 5; i++ {
		rotorOptions[i].GeoM.Scale(-0.40-(math.Floor(float64(i/3))*-0.15), 0.40-(math.Floor(float64(i/3))*0.15))
		rotorOptions[i].GeoM.Translate(float64((i*160)+150)-(math.Floor(float64(i/4))*160), float64(math.Floor(float64(i/4))*100))
	}
}

type Game struct {
	keys []ebiten.Key
}

func (g *Game) Update() error {
	g.keys = inpututil.AppendPressedKeys(g.keys[:0])
	if len(g.keys) == 0 {
		keyReleased = true
	}
	return nil
}

var keyReleased bool

var plugBoardLetters []string

var movingRotor bool

var rotorOptions = [5]*ebiten.DrawImageOptions{&ebiten.DrawImageOptions{}, &ebiten.DrawImageOptions{}, &ebiten.DrawImageOptions{}, &ebiten.DrawImageOptions{}, &ebiten.DrawImageOptions{}}

var oldMouseX int
var oldMouseY int

var rotorInMotion int

var rotorNbms = [5]string{"1", "2", "3", "4", "5"}

var selectedRotor string
var err error

var rotorsRotationAmounts = [5]int{0, 0, 0, 0, 0}

var plugBoard = map[interface{}]interface{}{
	"a": "",
	"b": " ",
	"c": " ",
	"d": " ",
	"e": " ",
	"f": " ",
	"g": " ",
	"h": " ",
	"i": " ",
	"j": " ",
	"k": " ",
	"l": " ",
	"m": " ",
	"n": " ",
	"o": " ",
	"p": " ",
	"q": " ",
	"r": " ",
	"s": " ",
	"t": " ",
	"u": " ",
	"v": " ",
	"w": " ",
	"x": " ",
	"y": " ",
	"z": " ",
}

func (g *Game) Draw(screen *ebiten.Image) {
	{
		for i := 0; i < 26; i++ {
			for _, j := range g.keys {
				if keyReleased == true {
					//SEND LETTER TO ENIGMA
					keyReleased = false
				}
				encryptedKey := encrypt(strings.ToLower(ebiten.Key.String(j)))
				if encryptedKey == keys[i] {
					//vector.DrawFilledCircle(screen, 400, 400, 100, color.RGBA{0x80, 0x00, 0x80, 0x80})
					text.Draw(screen, keys[i], mplusBigFont, i*25, 310, color.RGBA{255, 255, 0, 0xff})
				} else {
					text.Draw(screen, keys[i], mplusNormalFont, i*25, 310, color.Gray16{0xffff})
				}
			}
			if len(g.keys) == 0 {
				text.Draw(screen, keys[i], mplusNormalFont, i*25, 310, color.Gray16{0xffff})
			}
		}
		for i := 0; i < 26; i++ {
			var yPos int = i / 9
			var x int = i - int(math.Floor(float64(i/9)))*9
			var xPos int = x * 75
			text.Draw(screen, keys[i], mplusNormalFont, xPos, int(400+40*math.Floor(float64(yPos))), color.Gray16{0xffff})
			for k := 0; k < len(plugBoardLetters); k++ {
				if plugBoardLetters[k] == keys[i] {
					q := int(math.Floor(float64(k / 2)))
					r := uint8((q * 90) + 100)
					g := uint8((q * 70) + 10)
					b := uint8((q * 55))
					text.Draw(screen, keys[i], mplusNormalFont, xPos, int(400+40*math.Floor(float64(yPos))), color.RGBA{r, g, b, 0xff})
				}
			}
		}
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton(ebiten.MouseButtonLeft)) == true {
			xPos, yPos := ebiten.CursorPosition()
			keyXCord := (int(xPos*9) / screenWidth)
			keyYCord := (yPos / 40) - 9
			keyPosInArray := (keyYCord * 9) + keyXCord
			if keyPosInArray > 0 && keyPosInArray < len(keys) {
				plugBoardLetters = append(plugBoardLetters, keys[keyPosInArray])
				if len(plugBoardLetters)%2 == 0 {
					plugBoard[plugBoardLetters[len(plugBoardLetters)-1]] = plugBoardLetters[len(plugBoardLetters)-2]
					plugBoard[plugBoardLetters[len(plugBoardLetters)-2]] = plugBoardLetters[len(plugBoardLetters)-1]
					fmt.Println(plugBoard)
				}
			}
		}

		screen.DrawImage(reflectortopImg, nil)
		for i := 0; i < 3; i++ {
			options := &ebiten.DrawImageOptions{}
			options.GeoM.Translate(float64((i*160)+50), 1)
			screen.DrawImage(topbgImg, options)
			text.Draw(screen, strconv.Itoa(rotorsRotationAmounts[i]), mplusNormalFont, (i*160)+130, 230, color.White)
		}
		option := &ebiten.DrawImageOptions{}
		option.GeoM.Scale(1, 0.5)
		option.GeoM.Translate(530, 1)
		screen.DrawImage(emptyRotorSlotImg, option)
		option.GeoM.Translate(0, 110)
		screen.DrawImage(emptyRotorSlotImg, option)

		if movingRotor == true {
			mouseX, mouseY := ebiten.CursorPosition()
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton(ebiten.MouseButtonLeft)) == true {
				if mouseY < 200 {
					var clickedSlot int
					if mouseX < (3*160)+25 {
						clickedSlot = int(math.Floor(float64(mouseX-50) / 160))
					} else {
						clickedSlot = int(math.Floor(float64(mouseY/100))) + 3
					}
					if mouseX < (3*160)+25 {
						rotorOptions[rotorInMotion] = &ebiten.DrawImageOptions{}
						rotorOptions[rotorInMotion].GeoM.Scale(-0.40, 0.40)
						rotorOptions[rotorInMotion].GeoM.Translate(float64((math.Floor(float64((mouseX-50)/160)*160) + 155)), float64(10))
						oldMouseX, oldMouseY = int((math.Floor(float64((mouseX-50)/160)*160) + 155)), 10
					} else {
						rotorOptions[rotorInMotion] = &ebiten.DrawImageOptions{}
						rotorOptions[rotorInMotion].GeoM.Scale(-0.25, 0.25)
						rotorOptions[rotorInMotion].GeoM.Translate(float64((math.Floor(float64((mouseX-50)/160)*160) + 155)), math.Floor(float64(mouseY/100))*100)
						oldMouseX, oldMouseY = int((math.Floor(float64((mouseX-50)/160)*160) + 155)), int(math.Floor(float64(mouseY/100))*100)
					}
					if clickedSlot == rotorInMotion {
						movingRotor = false
						rotorNbms[clickedSlot] = selectedRotor
					} else {
						rotorOptions[clickedSlot] = rotorOptions[rotorInMotion]
						shortTermSelctedRotor := rotorNbms[clickedSlot]
						rotorNbms[clickedSlot] = selectedRotor
						selectedRotor = shortTermSelctedRotor
						rotorOptions[rotorInMotion] = &ebiten.DrawImageOptions{}
						rotorOptions[rotorInMotion].GeoM.Scale(-0.50, 0.50)
						rotorOptions[rotorInMotion].GeoM.Translate(float64(mouseX), float64(mouseY))
					}
				}
			} else {
				rotorOptions[rotorInMotion].GeoM.Translate(float64(mouseX-oldMouseX), float64(mouseY-oldMouseY))
				oldMouseX = mouseX
				oldMouseY = mouseY
			}
		} else {
			mouseX, mouseY := ebiten.CursorPosition()
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton(ebiten.MouseButtonLeft)) == true && mouseY < 200 && mouseX < screenWidth {
				movingRotor = true
				if mouseX < (3*160)+25 {
					rotorInMotion = int(math.Floor(float64(mouseX-50) / 160))
				} else {
					rotorInMotion = int(math.Floor(float64(mouseY/100))) + 3
				}
				rotorOptions[rotorInMotion] = &ebiten.DrawImageOptions{}
				rotorOptions[rotorInMotion].GeoM.Scale(-0.50, 0.50)
				rotorOptions[rotorInMotion].GeoM.Translate(float64(mouseX), float64(mouseY))
				oldMouseX = mouseX
				oldMouseY = mouseY
				selectedRotor = rotorNbms[rotorInMotion]
				rotorNbms[rotorInMotion] = ""
			}
		}
		screen.DrawImage(rotorImg, rotorOptions[0])
		screen.DrawImage(rotorImg, rotorOptions[1])
		screen.DrawImage(rotorImg, rotorOptions[2])
		screen.DrawImage(rotorImg, rotorOptions[3])
		screen.DrawImage(rotorImg, rotorOptions[4])

		for i := 0; i < 5; i++ {
			if i < 3 {
				text.Draw(screen, rotorNbms[i], mplusNormalFont, (i*160)+150, 200, color.RGBA{10, 100, 10, 0xff})
			} else {
				text.Draw(screen, rotorNbms[i], mplusNormalFont, 550, (int(math.Floor(float64(i/4)))*100)+70, color.RGBA{100, 100, 10, 0xff})
			}
		}

	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Enigma Machine")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}

func encrypt(key string) string {
	var output string
	output = plugBoardFunc(key)
	return output
}

func plugBoardFunc(inputedLetter string) string {
	scrambles := plugBoard
	var output string = inputedLetter
	if scrambles[inputedLetter] != " " {
		output = scrambles[inputedLetter].(string)
	}
	return output
}
