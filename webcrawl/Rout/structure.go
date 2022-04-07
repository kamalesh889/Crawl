package Rout

var ReqChannel = make(chan string)

type Request struct { //To get the request body
	Urls []string
}

type Dat struct { //Collecting url and crawl data
	Url  string
	Data string
}

type Response struct {
	Result []Dat
}
