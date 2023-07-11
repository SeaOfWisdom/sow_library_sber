package pinata

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/SeaOfWisdom/sow_library/src/service/storage"
)

// API Key: bdaf0e02c3e8ac255561
//  API Secret: d8df802a359885a6e13a54114d6fa50049a4807ee6949586f312986accb13271
//  JWT: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySW5mb3JtYXRpb24iOnsiaWQiOiI3OGQzZWUwNC04YmU3LTQ2NjQtYTQwMy1hYTEwYWZkYmQzODYiLCJlbWFpbCI6Imdvcm9ob3ZkYW5paWwxOUBnbWFpbC5jb20iLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZSwicGluX3BvbGljeSI6eyJyZWdpb25zIjpbeyJpZCI6IkZSQTEiLCJkZXNpcmVkUmVwbGljYXRpb25Db3VudCI6MX0seyJpZCI6Ik5ZQzEiLCJkZXNpcmVkUmVwbGljYXRpb25Db3VudCI6MX1dLCJ2ZXJzaW9uIjoxfSwibWZhX2VuYWJsZWQiOmZhbHNlLCJzdGF0dXMiOiJBQ1RJVkUifSwiYXV0aGVudGljYXRpb25UeXBlIjoic2NvcGVkS2V5Iiwic2NvcGVkS2V5S2V5IjoiYmRhZjBlMDJjM2U4YWMyNTU1NjEiLCJzY29wZWRLZXlTZWNyZXQiOiJkOGRmODAyYTM1OTg4NWE2ZTEzYTU0MTE0ZDZmYTUwMDQ5YTQ4MDdlZTY5NDk1ODZmMzEyOTg2YWNjYjEzMjcxIiwiaWF0IjoxNjgyNzAzMDYwfQ.1koXoIVOYFjng6KI2WbmalJPPY2rrUqwtXoSoQj0Ry4

const pinataJWT = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySW5mb3JtYXRpb24iOnsiaWQiOiI3OGQzZWUwNC04YmU3LTQ2NjQtYTQwMy1hYTEwYWZkYmQzODYiLCJlbWFpbCI6Imdvcm9ob3ZkYW5paWwxOUBnbWFpbC5jb20iLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZSwicGluX3BvbGljeSI6eyJyZWdpb25zIjpbeyJpZCI6IkZSQTEiLCJkZXNpcmVkUmVwbGljYXRpb25Db3VudCI6MX0seyJpZCI6Ik5ZQzEiLCJkZXNpcmVkUmVwbGljYXRpb25Db3VudCI6MX1dLCJ2ZXJzaW9uIjoxfSwibWZhX2VuYWJsZWQiOmZhbHNlLCJzdGF0dXMiOiJBQ1RJVkUifSwiYXV0aGVudGljYXRpb25UeXBlIjoic2NvcGVkS2V5Iiwic2NvcGVkS2V5S2V5IjoiYmRhZjBlMDJjM2U4YWMyNTU1NjEiLCJzY29wZWRLZXlTZWNyZXQiOiJkOGRmODAyYTM1OTg4NWE2ZTEzYTU0MTE0ZDZmYTUwMDQ5YTQ4MDdlZTY5NDk1ODZmMzEyOTg2YWNjYjEzMjcxIiwiaWF0IjoxNjgyNzAzMDYwfQ.1koXoIVOYFjng6KI2WbmalJPPY2rrUqwtXoSoQj0Ry4"

func PublishJson(work *storage.Work) string {

	// 	payload := strings.NewReader(`{
	//     "pinataOptions": {
	//         "cidVersion": 1
	//     },
	//     "pinataMetadata": {
	//         "name": "testing",
	//         "keyvalues": {
	//             "customKey": "customValue",
	//             "customKey2": "customValue2"
	//         }
	//     },
	//     "pinataContent": {
	//         "somekey":"somevalue"
	//     }
	// }`)

	req, err := http.NewRequest(
		"POST",
		"https://api.pinata.cloud/pinning/pinJSONToIPFS",
		strings.NewReader(""),
	)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", pinataJWT))

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	fmt.Println(string(body))
	return "ID"
}
