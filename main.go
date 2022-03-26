package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

var ctx = context.Background()

// Redis işlemlerini yapabilmek için global olarak tanımlanan redis client structı
var rdb *redis.Client

// Redis üzerinde işlem yapacağımız kanalın ismi
var channelName = "report"

func main() {
	// Redis client structımızı initialize ettik
	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	// Redis bağlantının sağlıklı bir şekilde yapılıp yapılmadığını kontrol ettik
	err := rdb.Ping(ctx).Err()
	if err != nil {
		panic("Redise bağlanılamadı => " + err.Error())
	}

	// Subscriber(abone): kanalı dinleme işlemlerini yapan fonksiyonu asenkron olarak çağırdık
	go subscriber(1)
	go subscriber(2)
	// Publisher(yayıncı): kanala yayın yapacak olan fonksiyonumuzu asenkron olarak çağırdık
	go publisher()

	// Main metodu tamamlanıp proje kapanmasın diye boş bir channel oluşturduk ve channeli dinlemeye başladık
	c := make(chan int)
	<-c
}

// Her 3 saniyede bir kanala mesaj publish eden fonksiyon
// Fonksiyon her 3 saniyede 1 kez o anki saati kanala publish etmektedir
func publisher() {
	for range time.Tick(time.Second * 3) {
		// Zaman değişkenini oluşturduk
		t := time.Now().Format("15:04:05")
		fmt.Println("*************************")
		fmt.Println("Kanala gönderilen => " + t)
		// Zaman değişkenini kanala publish ettik
		rdb.Publish(ctx, channelName, t)
	}
}

// Kanala gönderilen mesajları dinleyen fonksiyon
// SubscriberNumber değişkeni birden fazla kanalı dinleyen subscriber(abone)'yi simule edilebilmek için yapıldı.
// Fonksiyon yukarıda 2 defa ve farklı numaralar ile çağırıldı
func subscriber(subscriberNumber int) {
	// kanala abone oluyoruz
	subs := rdb.Subscribe(ctx, channelName)
	// kanaldan gelen mesajları dinlemeye başlıyoruz
	for msg := range subs.Channel() {
		// Kanala gelen mesajın içeriğini(Payload ile) konsola yazdırdık
		fmt.Println(fmt.Sprintf("Subscribe %d için kanaldan okunan => %s", subscriberNumber, msg.Payload))
	}
}
