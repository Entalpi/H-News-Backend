package services

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"log"
	"strconv"
	"time"

	"hnews/Godeps/_workspace/src/github.com/boltdb/bolt"
)

var (
	newsdb, _ = bolt.Open("news.db", 0644, nil)
	commdb, _ = bolt.Open("comments.db", 0644, nil)
)

// PrintDBStats prints our database statistics once per second.
func PrintDBStats() {
	go func() {
		for {
			log.Println(newsdb.Stats())
			log.Println(commdb.Stats())
			time.Sleep(4 * time.Second)
		}
	}()
}

/** Database Service **/

// News represent one news story/item on Hacker News
type News struct {
	ID       int32     `json:"id"` // Ints need to be a defined size since binary.Read does not support just 'int'
	Rank     int32     `json:"rank"`
	Title    string    `json:"title"`
	Link     string    `json:"link"`
	Author   string    `json:"author"`
	Points   int32     `json:"points"`
	Time     time.Time `json:"time"`
	Comments int32     `json:"comments"` // Number of comments on the News
}

// Close closes database connections
func Close() {
	newsdb.Close()
	commdb.Close()
}

// ReadNews ...
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

// SaveNews saves the News in the DB
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

// ReadNewsIds Read all the keys from News db
func ReadNewsIds() []int32 {
	var ids []int32
	newsdb.View(func(tx *bolt.Tx) error {
		for i := 1; i < 480; i++ {
			b := tx.Bucket([]byte(strconv.Itoa(int(i))))
			if b == nil {
				// log.Println("Bucket", i, "not found.")
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

// A Comment on a News
type Comment struct {
	Num      int32     `json:"num"`      // The ith comment on the post
	ParentID int32     `json:"parentid"` // ID of the News
	ID       int32     `json:"id"`       // The Comments unique ID
	Offset   int32     `json:"offset"`   // Level of offset for the Comment
	Time     time.Time `json:"time"`
	Author   string    `json:"author"`
	Text     string    `json:"text"`
}

// SaveComments Dumps the Comments into the newsDB as JSON.
func SaveComments(comments []Comment) {
	if len(comments) == 0 {
		return
	}
	newsid := comments[0].ParentID
	commdb.Update(func(tx *bolt.Tx) error {
		k := strconv.Itoa(int(newsid))
		b, err := tx.CreateBucketIfNotExists([]byte(k))
		if err != nil {
			log.Println("SaveComments:", err)
			return err
		}
		for _, comment := range comments {
			v, err := json.Marshal(comment)
			if err != nil {
				log.Println("SaveComments:", err)
				continue
			}
			b.Put([]byte(strconv.Itoa(int(comment.Num))), []byte(v))
		}
		return nil
	})
}

// ReadComments Returns the comments on the News item specified by the id.
func ReadComments(newsid int, from int, to int) []Comment {
	var comments []Comment
	commdb.View(func(tx *bolt.Tx) error {
		k := strconv.Itoa(newsid)
		b := tx.Bucket([]byte(k))
		if b == nil {
			return nil
		}
		for i := from; i < to; i++ {
			k := []byte(strconv.Itoa(i))
			v := b.Get(k)
			str := string(v)
			if len(str) == 0 {
				continue
			}
			var comment Comment
			json.Unmarshal([]byte(str), &comment)
			comments = append(comments, comment)
		}
		return nil
	})
	return comments
}
