package services

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"hnews/Godeps/_workspace/src/github.com/boltdb/bolt"
	"hnews/Godeps/_workspace/src/github.com/headzoo/surf"
	"log"
	"strconv"
	"time"
)

var (
	newsdb, _ = bolt.Open("a.db", 0644, nil)
	commdb, _ = bolt.Open("b.db", 0644, nil)
)

/** Login service **/
func Login(username string, password string) bool {
	bow := surf.NewBrowser()
	err := bow.Open("https://news.ycombinator.com/login?goto=news")
	if err != nil {
		return false
	}

	fm, _ := bow.Form("form")
	fm.Input("acct", username)
	fm.Input("pw", password)
	erro := fm.Submit()
	if erro != nil {
		log.Println(erro)
		return false
	}
	log.Println(fm)

	return true
}

/** Database Service **/
// Ints need to be a defined size since binary.Read does not support just 'int'
type News struct {
	ID       int32     `json:"id"`
	Rank     int32     `json:"rank"`
	Title    string    `json:"title"`
	Link     string    `json:"link"`
	Author   string    `json:"author"`
	Points   int32     `json:"points"`
	Time     time.Time `json:"time"`
	Comments int32     `json:"comments"` // Number of comments on the News
}

func ReadNews(from int, to int) []News {
	var news []News
	newsdb.View(func(tx *bolt.Tx) error {
		for i := from; i <= to; i++ {
			b := tx.Bucket([]byte(strconv.Itoa(int(i))))
			if b == nil {
				log.Println("Bucket", i, "not found.")
				continue
			}

			title := string(b.Get([]byte("title")))
			author := string(b.Get([]byte("author")))
			link := string(b.Get([]byte("link")))

			var rank int32
			binary.Read(bytes.NewReader(b.Get([]byte("rank"))), binary.LittleEndian, &rank)

			var t int64
			binary.Read(bytes.NewReader(b.Get([]byte("time"))), binary.LittleEndian, &t)

			var points int32
			binary.Read(bytes.NewReader(b.Get([]byte("points"))), binary.LittleEndian, &points)

			var comments int32
			binary.Read(bytes.NewReader(b.Get([]byte("comments"))), binary.LittleEndian, &comments)

			var id int32
			binary.Read(bytes.NewReader(b.Get([]byte("id"))), binary.LittleEndian, &id)

			news = append(news, News{id, rank, title, link, author, points, time.Unix(t, 0), comments})
		}
		return nil
	})
	return news
}

func SaveNews(news []News) {
	newsdb.Update(func(tx *bolt.Tx) error {
		for _, aNews := range news {
			b, err := tx.CreateBucketIfNotExists([]byte(strconv.Itoa(int(aNews.Rank))))
			if err != nil {
				log.Println(err)
				return err
			}

			b.Put([]byte("title"), []byte(aNews.Title))
			b.Put([]byte("link"), []byte(aNews.Link))
			b.Put([]byte("author"), []byte(aNews.Author))

			var t bytes.Buffer
			binary.Write(&t, binary.LittleEndian, aNews.Time.Unix())
			b.Put([]byte("time"), []byte(t.Bytes()))

			var r bytes.Buffer
			binary.Write(&r, binary.LittleEndian, aNews.Rank)
			b.Put([]byte("rank"), []byte(r.Bytes()))

			var p bytes.Buffer
			binary.Write(&p, binary.LittleEndian, aNews.Points)
			b.Put([]byte("points"), []byte(p.Bytes()))

			var c bytes.Buffer
			binary.Write(&c, binary.LittleEndian, aNews.Comments)
			b.Put([]byte("comments"), []byte(c.Bytes()))

			var id bytes.Buffer
			binary.Write(&id, binary.LittleEndian, aNews.ID)
			b.Put([]byte("id"), []byte(id.Bytes()))

			if err != nil {
				log.Println("Err SaveNews:", err)
				return err
			}
		}
		return nil
	})
}

// Read all the keys from News db
func ReadNewsIds() []int32 {
	var ids []int32
	newsdb.View(func(tx *bolt.Tx) error {
		for i := 1; i < 480; i++ {
			b := tx.Bucket([]byte(strconv.Itoa(int(i))))
			if b == nil {
				log.Println("Bucket", i, "not found.")
				continue
			}
			var id int32
			binary.Read(bytes.NewReader(b.Get([]byte("id"))), binary.LittleEndian, &id)
			ids = append(ids, id)
		}
		return nil

	})
	return ids
}

type Comment struct {
	ParentID int32     `json:"parentid"` // ID of the News
	ID       int32     `json:"id"`       // The Comments unique ID
	Offset   int32     `json:"offset"`   // Level of offset for the Comment
	Time     time.Time `json:"time"`
	Author   string    `json:"author"`
	Text     string    `json:"text"`
}

// Dumps the Comments into the newsDB as JSON.
func SaveComments(comments []Comment) {
	if len(comments) == 0 {
		return
	}
	newsid := comments[0].ParentID
	json, err := json.Marshal(comments)
	if err != nil {
		log.Println("SaveComments:", err)
	}
	commdb.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("comments"))
		if err != nil {
			log.Println("SaveComments:", err)
			return err
		}
		b.Put([]byte(strconv.Itoa(int(newsid))), []byte(json))
		return nil
	})
}

// Returns the comments on the News item specified by the id.
func ReadComments(newsid int) string {
	var json string
	commdb.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("comments")) // Might not even need a bucket.
		k := strconv.Itoa(newsid)
		json = string(b.Get([]byte(k)))
		return nil
	})
	return json
}
