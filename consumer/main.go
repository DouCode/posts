//package main
//
//import (
//	"log"
//	"os"
//
//	"github.com/streadway/amqp"
//)
//
//func main() {
//	amqpConnection, err := amqp.Dial(os.Getenv("RABBITMQ_URI"))
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer amqpConnection.Close()
//	channelAmqp, _ := amqpConnection.Channel()
//	defer channelAmqp.Close()
//	forever := make(chan bool)
//	msgs, err := channelAmqp.Consume(
//		os.Getenv("RABBITMQ_QUEUE"),
//		"",
//		true,
//		false,
//		false,
//		false,
//		nil,
//	)
//	go func() {
//		for d := range msgs {
//			log.Printf("Received a message: %s", d.Body)
//		}
//	}()
//	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
//	<-forever
//}

//package main
//
//import (
//	"context"
//	"encoding/json"
//	"encoding/xml"
//	"fmt"
//	"go.mongodb.org/mongo-driver/bson"
//	"io/ioutil"
//	"log"
//	"net/http"
//	"os"
//
//	"github.com/streadway/amqp"
//	"go.mongodb.org/mongo-driver/mongo"
//	"go.mongodb.org/mongo-driver/mongo/options"
//)
//
//type Request struct {
//	URL string `json:"url"`
//}
//
//type Feed struct {
//	Entries []Entry `xml:"entry"`
//}
//
//type Entry struct {
//	Link struct {
//		Href string `xml:"href,attr"`
//	} `xml:"link"`
//	Thumbnail struct {
//		URL string `xml:"url,attr"`
//	} `xml:"thumbnail"`
//	Title string `xml:"title"`
//}
//
//func GetFeedEntries(url string) ([]Entry, error) {
//	httpClient := &http.Client{}
//	req, err := http.NewRequest("GET", url, nil)
//	if err != nil {
//		return nil, err
//	}
//	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.110 Safari/537.36")
//
//	resp, err := httpClient.Do(req)
//	if err != nil {
//		return nil, err
//	}
//	defer resp.Body.Close()
//
//	byteValue, _ := ioutil.ReadAll(resp.Body)
//	var feed Feed
//	xml.Unmarshal(byteValue, &feed)
//
//	return feed.Entries, nil
//}
//
//func main() {
//	ctx := context.Background()
//	mongoClient, _ := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
//	defer mongoClient.Disconnect(ctx)
//
//	amqpConnection, err := amqp.Dial(os.Getenv("RABBITMQ_URI"))
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer amqpConnection.Close()
//
//	channelAmqp, _ := amqpConnection.Channel()
//	defer channelAmqp.Close()
//
//	forever := make(chan bool)
//
//	msgs, err := channelAmqp.Consume(
//		os.Getenv("RABBITMQ_QUEUE"),
//		"",
//		true,
//		false,
//		false,
//		false,
//		nil,
//	)
//
//	go func() {
//		for d := range msgs {
//			log.Printf("Received a message: %s", d.Body)
//
//			var request Request
//			json.Unmarshal(d.Body, &request)
//
//			log.Println("RSS URL:", request.URL)
//
//			entries, _ := GetFeedEntries(request.URL)
//
//			fmt.Println(entries)
//
//			collection := mongoClient.Database(os.Getenv("MONGO_DATABASE")).Collection("recipes")
//			fmt.Println(len(entries))
//			for _, entry := range entries {
//				collection.InsertOne(ctx, bson.M{
//					"title":     entry.Title,
//					"thumbnail": entry.Thumbnail.URL,
//					"url":       entry.Link.Href,
//				})
//			}
//		}
//	}()
//
//	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
//	<-forever
//}

package main

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/anaskhan96/soup"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Request struct {
	URL string `json:"url"`
}

type Link struct {
	Text string
	Href string
}

func GetLinks(url string) ([]Link, error) {
	// 发送Get请求
	rsp, err := http.Get(url)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	content := string(body)
	defer rsp.Body.Close()

	slice := make([]Link, 0)

	// 下面主要是解析标签
	doc := soup.HTMLParse(content)
	links := doc.Find("div", "class", "cc-dc-c").FindAll("a")
	//log.Println(doc)
	for _, link := range links {
		//fmt.Println(link.Text(), "| Link :", link.Attrs()["href"])
		t := Link{Text: link.Text(), Href: "https://tophub.today" + link.Attrs()["href"]}
		slice = append(slice, t)
	}

	return slice, nil

}

func main() {
	ctx := context.Background()
	mongoClient, _ := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	defer mongoClient.Disconnect(ctx)

	amqpConnection, err := amqp.Dial(os.Getenv("RABBITMQ_URI"))
	if err != nil {
		log.Fatal(err)
	}
	defer amqpConnection.Close()

	channelAmqp, _ := amqpConnection.Channel()
	defer channelAmqp.Close()

	forever := make(chan bool)

	msgs, err := channelAmqp.Consume(
		os.Getenv("RABBITMQ_QUEUE"),
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)

			var request Request
			json.Unmarshal(d.Body, &request)

			log.Println("RSS URL:", request.URL)

			links, _ := GetLinks(request.URL)

			collection := mongoClient.Database(os.Getenv("MONGO_DATABASE")).Collection("links")
			log.Println("collection:", collection.Database().Name())
			log.Println("collection:", collection.Name())
			log.Println(len(links))
			for _, link := range links {
				collection.InsertOne(ctx, bson.M{
					"text": link.Text,
					"href": link.Href,
				})
			}
			log.Println(len(links))
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
