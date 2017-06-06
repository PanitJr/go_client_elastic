package go_client_elastic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

// DateLayout is a date formular to serarch in @timestamp field
const DateLayout = "2006-01-02"

//Shards is struct for keep the _shards part of response from elasticsearch
type Shards struct {
	Total      int `json:"total"`
	Successful int `json:"successful"`
	Failed     int `json:"failed"`
}
type Beat struct {
	Hostname string `json:"hostname,omitempty"`
	Name     string `json:"name,omitempty"`
	Version  string `json:"version,omitempty"`
}

//Source is struct for keep the _source part of response from elasticsearch
type Source struct {
	DestNodeResultCode string     `json:"DestNodeResultCode,omitempty"`
	ResTimeStamp       string     `json:"ResTimeStamp,omitempty"`
	CmdName            string     `json:"CmdName,omitempty"`
	Instance           string     `json:"Instance,omitempty"`
	DestNodeCmd        string     `json:"DestNodeCmd,omitempty"`
	Message            string     `json:"message,omitempty"`
	TimeStamp          string     `json:"TimeStamp,omitempty"`
	AppName            string     `json:"AppName,omitempty"`
	App                string     `json:"App,omitempty"`
	ArrayofDR          string     `json:"ArrayofDR,omitempty"`
	Timestamp          *time.Time `json:"@timestamp,omitempty"`
	InitInvoke         string     `json:"InitInvoke,omitempty"`
	Version            string     `json:"@version,omitempty"`
	Host               string     `json:"host,omitempty"`
	UsageTime          string     `json:"UsageTime,omitempty"`
	ResultDesc         string     `json:"ResultDesc,omitempty"`
	ReqTimeStamp       string     `json:"ReqTimeStamp,omitempty"`
	Identity           string     `json:"Identity,omitempty"`
	DestNodeResultDesc string     `json:"DestNodeResultDesc,omitempty"`
	DestNodeName       string     `json:"DestNodeName,omitempty"`
	HostName           string     `json:"Hostname,omitempty"`
	ResultCode         string     `json:"ResultCode,omitempty"`
	Session            string     `json:"Session,omitempty"`
	HostIP             string     `json:"HostIP,omitempty"`
	Port               int        `json:"Port,omitempty"`
	Process            string     `json:"Process,omitempty"`
	Source             string     `json:"source,omitempty"`
	Type               string     `json:"type,omitempty"`
	DateTime           string     `json:"DateTime,omitempty"`
	BaseLineQueue      int        `json:"Base_line_queue,omitempty"`
	Tps                int        `json:"Tps,omitempty"`
	MemoryUsed         int        `json:"Memory_Used,omitempty"`
	*Beat              `json:"beat,omitempty"`
	Protocol           string   `json:"Protocol,omitempty"`
	Offset             int      `json:"offset,omitempty"`
	CPUUsed            int      `json:"CPU_Used,omitempty"`
	DestIP             string   `json:"DestIP,omitempty"`
	InputType          string   `json:"input_type,omitempty"`
	Service            string   `json:"Service,omitempty"`
	BaseLineTps        int      `json:"Base_line_tps,omitempty"`
	QueueCount         int      `json:"Queue_count,omitempty"`
	Tags               []string `json:"tags,omitempty"`
	SessionCount       int      `json:"Session_count,omitempty"`
	BaseLineSession    int      `json:"Base_line_session,omitempty"`
}

//Hit is struct for keep the hits(elsment) part of response from elasticsearch
type Hit []struct {
	Index  string      `json:"_index"`
	Type   string      `json:"_type"`
	ID     string      `json:"_id"`
	Score  interface{} `json:"_score"`
	Source `json:"_source"`
}

//Hits is struct for keep the hits(array) part of response from elasticsearch
type Hits struct {
	Total    int         `json:"total"`
	MaxScore interface{} `json:"max_score"`
	Hit      `json:"hits"`
}
type Response struct {
	Response []Result `json:"responses"`
}

//Result is a big struct for keep the all parts of response from elasticsearch
type Result struct {
	Took     int  `json:"took"`
	TimedOut bool `json:"timed_out"`
	Shards   `json:"_shards"`
	Hits     `json:"hits"`
	Status   int `json:"status"`
}

//Result is a big struct for keep the all parts of response from elasticsearch

//Timestamp is the struct for build request for range query type
//to serarch in field @timestamp only
type Timestamp struct {
	Gte    string `json:"gte"`
	Lte    string `json:"lte"`
	Format string `json:"format"`
}

//DateRange is the range type for build request for elasticsearch query
type DateRange struct {
	Timestamp `json:"@timestamp"`
}

//Must is the struct for build must condition for elasticsearch query
//TermSource condition mean the field of source must have this term
//DateRange for query range of field @timestamp must in this duration
//MatchSource condition mean the field of source must have this text
type Must []struct {
	TermSource  *Source `json:"term,omitempty"`
	*DateRange  `json:"range,omitempty"`
	MatchSource *Source `json:"match,omitempty"`
}

//MustNot is the struct for build must_not condition for elasticsearch query
//TermSource condition mean the field of source must_not have this term
//DateRange for query range of field @timestamp, must_not in this duration
//MatchSource condition mean the field of source must_not have this text
type MustNot []struct {
	TermSource  *Source `json:"term,omitempty"`
	*DateRange  `json:"range,omitempty"`
	MatchSource *Source `json:"match,omitempty"`
}

//Bool is the struct for build Bool structure for elasticsearch query,
//which is contain Must and MustNot struct
type Bool struct {
	*Must    `json:"must,omitempty"`
	*MustNot `json:"must_not,omitempty"`
}

//Query is the struct for build Query structure for elasticsearch query,
//which is contain bool struct
type Query struct {
	Bool `json:"bool"`
}

//SearchReqBody is the struct for contain request Body for elasticsearch API
type SearchReqBody struct {
	Query `json:"query"`
}

//GCERequestBuilder is the patttern that user can use this package
type GCERequestBuilder struct {
	Host    string   `json:"host"`
	Port    string   `json:"port"`
	Index   string   `json:"index"`
	Type    string   `json:"type"`
	Fields  []string `json:"fields"`
	Must    `json:"must"`
	MustNot `json:"mustnot"`
}
type AppHostReqstBuilder struct {
	Host         string `json:"host" binding:"required"`
	Port         string `json:"port" binding:"required"`
	Index        string `json:"index" binding:"required"`
	Type         string `json:"type" binding:"required"`
	AppHostParam `json:"AppHostParam" binding:"required"`
}
type AppHostParam struct {
	Param []AppHost `json:"Param" binding:"required"`
}
type AppHost struct {
	App  string `json:"app" binding:"required"`
	Host string `json:"host" binding:"required"`
}

func GetList(r GCERequestBuilder) (Result, string) {
	getListError := ""
	var reqBody []byte
	if r.MustNot != nil || r.Must != nil {
		searchReq := &SearchReqBody{
			Query: Query{
				Bool: Bool{
					Must:    &r.Must,
					MustNot: &r.MustNot,
				},
			},
		}
		if r.MustNot == nil {
			searchReq.Query.Bool.MustNot = nil
		}
		if r.Must == nil {
			searchReq.Query.Bool.Must = nil
		}
		reqBody, err := json.Marshal(searchReq)
		fmt.Println("Json elastic request body:", string(reqBody))

		if err != nil {
			fmt.Println("Marshal error:", err)
			getListError += fmt.Sprintf("Marshal error : %v", err)
		}
	}
	reqURL := ""
	if r.Fields == nil {
		reqURL = fmt.Sprintf("http://%s:%s/%s/%s/_search?size=1&sort=@timestamp:desc",
			url.QueryEscape(r.Host),
			url.QueryEscape(r.Port),
			url.QueryEscape(r.Index),
			url.QueryEscape(r.Type))
	} else {
		reqURL = fmt.Sprintf("http://%s:%s/%s/%s/_search?size=1&sort=@timestamp:desc&_source=*.id,%s",
			url.QueryEscape(r.Host),
			url.QueryEscape(r.Port),
			url.QueryEscape(r.Index),
			url.QueryEscape(r.Type),
			url.QueryEscape(stringField(r.Fields)))
	}

	req, err := http.NewRequest("GET", reqURL, bytes.NewBuffer(reqBody))
	if err != nil {
		log.Fatal("NewRequest: ", err)
		getListError += fmt.Sprintf("NewRequest: %v", err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("client.Do: ", err)
		getListError += fmt.Sprintf("client.Do: %v", err)
	}
	defer resp.Body.Close()
	var record Result
	if err := json.NewDecoder(resp.Body).Decode(&record); err != nil {
		log.Fatal(err)
		getListError += fmt.Sprintf("NewDecoder: %v", err)
	}

	return record, getListError
}
func stringField(f []string) string {
	res := ""
	for _, val := range f {
		res += val
		res += ","
	}
	return res
}
func dateRangeBuild(from, to time.Time) *DateRange {
	dateRange := &DateRange{
		Timestamp: Timestamp{
			Gte:    from.Format(DateLayout),
			Lte:    to.Format(DateLayout),
			Format: "yyyy-MM-dd||yyyy-MM-dd",
		},
	}
	return dateRange
}
func GetByAppHost(r AppHostReqstBuilder) interface{} {
	reqBody := ``
	output := make(map[string]Result)

	for _, appHost := range r.AppHostParam.Param {
		reqBody += fmt.Sprintf(`{"index":"%s","type":"%s"}`, r.Index, r.Type)
		reqBody += fmt.Sprintf("\n")
		reqBody += fmt.Sprintf(`{"query":{"bool":{"must":[{"match":{"App":"%s"}},{"match":{"Hostname":"%s"}}]}},"size":1,"sort": {"@timestamp": "desc"}}`, appHost.App, appHost.Host)
		reqBody += fmt.Sprintf("\n")
	}

	reqURL := fmt.Sprintf("http://%s:%s/_msearch",
		url.QueryEscape("10.138.32.97"),
		url.QueryEscape("8080"))

	fmt.Printf("URL = %+v \n", reqURL)
	fmt.Printf("body = %+v \n", reqBody)
	req, err := http.NewRequest("POST", reqURL, bytes.NewBuffer([]byte(reqBody)))
	if err != nil {
		log.Fatal("NewRequest: ", err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("client.Do: ", err)
	}
	defer resp.Body.Close()
	var records Response
	if err := json.NewDecoder(resp.Body).Decode(&records); err != nil {
		log.Fatal(err)
	}
	for _, record := range records.Response {
		output[record.Hits.Hit[0].Source.HostName+"-"+record.Hits.Hit[0].Source.App] = record
	}
	return output
}

/*func main() {
	requestMap := AppHostParam{Param: []AppHost{
		AppHost{App: "dOCF", Host: "DINIAMOA004G"},
		AppHost{App: "alpha", Host: "DINIAMOA004G"},
		AppHost{App: "VOCF", Host: "DINIAMOA004G"},
	}}
	requestMap.Param = append(requestMap.Param, AppHost{App: "dOCF", Host: "DINIAMOA002G"})
	getByAppHost(AppHostReqstBuilder{
		Host:         "10.138.32.97",
		Port:         "8080",
		Index:        "interstellar-crawler",
		Type:         "log",
		AppHostParam: requestMap,
	})
}*/
