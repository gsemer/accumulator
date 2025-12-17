package persistence

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type BackgroundRepository struct {
	client http.Client
	URL    string
}

func NewBackgroundRepository(client http.Client, URL string) *BackgroundRepository {
	return &BackgroundRepository{client: client, URL: URL}
}

func (br *BackgroundRepository) FetchAccumulator() (int64, error) {
	request, err := http.NewRequest("GET", fmt.Sprintf("%s/state", br.URL), nil)
	if err != nil {
		return 0, err
	}
	query := request.URL.Query()
	query.Set("format", "sum")
	request.URL.RawQuery = query.Encode()

	response, err := br.client.Do(request)
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		return 0, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return 0, err
	}

	var JSONResponse struct {
		Res int64 `json:"res"`
	}
	err = json.Unmarshal(body, &JSONResponse)
	if err != nil {
		return 0, err
	}

	accumulator := JSONResponse.Res
	return accumulator, nil
}

func (br *BackgroundRepository) FetchValues() ([]int64, error) {
	request, err := http.NewRequest("GET", fmt.Sprintf("%s/state", br.URL), nil)
	if err != nil {
		return nil, err
	}
	query := request.URL.Query()
	query.Set("format", "list")
	request.URL.RawQuery = query.Encode()

	response, err := br.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var JSONResponse struct {
		Res []int64 `json:"res"`
	}
	err = json.Unmarshal(body, &JSONResponse)
	if err != nil {
		return nil, err
	}

	values := JSONResponse.Res
	return values, nil
}
