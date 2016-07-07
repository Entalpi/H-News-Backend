# H-News-Backend

# Introduction
The backend is up and running on a hobby dyno over at [Heroku](https://h-news.herokuapp.com).

It consists of one public API written in Go with Gin and this API wraps calls
to a other HTTP server running on the same machine at Heroku. This is because
the other API written in Ruby with Sinatra takes care of all the functionality
surrounding login and posting stuff to Hacker News.

# API documentation

## Endpoints
All endpoints are prefixed with '/v1'
### GET /top
URL params: from: Int, to: Int
Returns the top stories from H.N front page from index 'from' to index 'to'.

### GET /newest
URL params: from: Int, to: Int
Returns the newest stories from H.N front page from index 'from' to index 'to'.

### GET /show
URL params: from: Int, to: Int
Returns the show stories from H.N front page from index 'from' to index 'to'.

### GET /ask
URL params: from: Int, to: Int
Returns the ask stories from H.N front page from index 'from' to index 'to'.

### GET /comments
Each item (news story, comment) at Hacker News has a unique ID and this is used
to lookup and scrape a specific comment.

#### Example
TDA

# License
The MIT License (MIT)
Copyright (c) 2015 Alexander Lingtorp

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
