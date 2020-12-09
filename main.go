package main

import (
	"github.com/disintegration/imaging"
	"image"
	"image/jpeg"
	_ "image/png"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

func main() {
	b, err := tb.NewBot(tb.Settings{

		Token:  "1411061889:AAHATOKI6Vi-Zt0XrfH0jEdxD3obqVnOfPw", // Insert your telegram bot token
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	// Функция принимает команду на ролл. Если команда корректная, то ролит случайное число и отправляет его.
	b.Handle("/roll", func(m *tb.Message) {
		if m.Payload[0] != 'd' {
			b.Send(m.Sender, "Ошибка. Неизвестная команда. Команды должны начинаться с d")
			return
		}
		number, err := strconv.Atoi(m.Payload[1:])
		if err != nil {
			b.Send(m.Sender, "Такого кубика не существует")
			return
		}
		if !(number == 4 || number == 20 || number == 6 || number == 12 || number == 8) {
			b.Send(m.Sender, "Такого кубика не существует")
			return
		}
		result := rand.Intn(number)
		b.Send(m.Sender, strconv.Itoa(result))
	})

	// Функция принимает фото, делает его черно-белым, размывает и возвращает
	b.Handle(tb.OnPhoto, func(m *tb.Message) {
		file, _ := b.GetFile(&m.Photo.File)
		result, _, err := image.Decode(file)

		if err != nil {
			b.Send(m.Sender, "Ошибка распознавания изображения")
			return
		}

		img := imaging.Grayscale(result)
		img = imaging.Blur(img, 2)
		f, err := os.Create("img.jpg")

		if err != nil {
			panic(err)
		}

		defer f.Close()
		jpeg.Encode(f, img, nil)
		image := &tb.Photo{File: tb.FromDisk("img.jpg")}
		b.Send(m.Sender, image)
	})

	b.Start()
}
