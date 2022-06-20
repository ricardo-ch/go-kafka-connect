package connectors

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/jarcoal/httpmock.v1"
	"net/http"
	"testing"
)

// This test a side effect of resty: when retry on 409 the error of response is not reinitialized
func Test_Get(t *testing.T) {

	client := newBaseClient("http://randomurl")
	// mock HTTP response
	{
		typedClient := client.(*baseClient)
		httpmock.Reset()
		httpmock.ActivateNonDefault(typedClient.restClient.GetClient())
		defer httpmock.DeactivateAndReset()

		i := 0
		myresponder := func(req *http.Request) (*http.Response, error) {
			i++
			if i == 1 {
				jsonresp, _ := httpmock.NewJsonResponse(409, ErrorResponse{Message: "some random msg"})
				res := new(http.Response)
				*res = *jsonresp
				res.Request = req
				return res, nil
			}
			jsonresp, _ := httpmock.NewJsonResponse(200, ConnectorResponse{Name: "test"})
			res := new(http.Response)
			*res = *jsonresp
			res.Request = req
			return res, nil
		}

		httpmock.RegisterResponder("GET", "http://randomurl/connectors/test", myresponder)
	}

	//Act
	connector, err := client.GetConnector(ConnectorRequest{Name: "test"})

	assert.Equal(t, "test", connector.Name)
	assert.NoError(t, err)
}
