package function

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"

	pb "handler/function/generated_grpc"

	handler "github.com/openfaas/templates-sdk/go-http"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type Rerank_request struct {
	User_id         string     `json:"user_id"`
	Item_ids        [20]string `json:"item_ids"`
	Qid             string     `json:"qid"`
	Normalise_items bool       `json:"normalise_items"`
	Normalise_user  bool       `json:"normalise_user"`
	Use_case        string     `json:"use_case"`
}
type Rerank_response struct {
	User_id string             `json:"user_id"`
	Scores  map[string]float64 `json:"scores"`
	Qid     string             `json:"qid"`
	Status  int                `json:"status"`
}
type Pair struct {
	Key   string
	Value float64
}

type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PairList) Less(i, j int) bool { return p[j].Value < p[i].Value }

func Handle(req handler.Request) (handler.Response, error) {
	var request_body map[string]interface{}
	fmt.Println(string(req.QueryString))
	json.Unmarshal([]byte(string(req.Body)), &request_body)
	uid := strings.Split(req.QueryString, "=")[1]
	var conn *grpc.ClientConn
	conn, err := grpc.Dial("host.docker.internal:56937", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	c := pb.NewFunctionServiceClient(conn)

	contextRequest := pb.LiveContextRequest{
		Uid:      uid,
		ArtistId: "",
	}
	request := pb.FunctionRequest{
		Query: &contextRequest,
	}

	response, err := c.GetFunctionResponse(context.Background(), &request)
	if err != nil {
		log.Fatalf("Error when calling onloop: %s", err)
	}

	var contents [20]string
	length := len(response.Contents)
	if length > 20 {
		length = 20
	}
	for i := 0; i < length; i++ {
		contents[i] = response.Contents[i].GetContentId()
	}
	rerank_request := &Rerank_request{User_id: uid, Item_ids: contents, Qid: "25", Normalise_items: false, Normalise_user: false, Use_case: "song"}
	postBody, err := json.Marshal(rerank_request)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	responseBody := bytes.NewBuffer(postBody)
	resp, err := http.Post("http://discovery-monitoring.wynk.in/v2/reranking/rerank", "application/json", responseBody)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()
	//Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	// sb := string(body)
	// log.Printf(sb)
	rerank_response := Rerank_response{}
	if resp.StatusCode == http.StatusOK {
		json.Unmarshal(body, &rerank_response)
	}
	log.Println(rerank_response)

	p := make(PairList, len(rerank_response.Scores))

	i := 0
	for k, v := range rerank_response.Scores {
		p[i] = Pair{k, v}
		i++
	}

	var sortedSongs [20]string
	sort.Sort(p)
	for i := 0; i < len(p); i++ {
		sortedSongs[i] = p[i].Key
	}
	log.Println("-----------------")
	log.Println(sortedSongs)
	log.Println(request_body)

	request_body["content"] = sortedSongs
	log.Println(request_body)

	message, err := json.Marshal(request_body)
	if err != nil {
		log.Fatalln(err)
	}
	var testrr map[string]interface{}
	json.Unmarshal(message, &testrr)
	log.Println(testrr)
	total_response := &handler.Response{
		Body:       message,
		StatusCode: http.StatusOK,
		Header:     req.Header,
	}
	log.Println(*total_response)
	return handler.Response{
		Body:       message,
		StatusCode: http.StatusOK,
		// Header:     req.Header,
	}, err

}
