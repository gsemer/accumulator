package persistence

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"signal/domain"
)

type BackgroundRepository struct {
	client http.Client
	URL    string
}

func NewBackgroundRepository(client http.Client, URL string) *BackgroundRepository {
	return &BackgroundRepository{client: client, URL: URL}
}

type StateResult struct {
	Res domain.State `json:"res"`
}

func (br *BackgroundRepository) FetchData() (domain.State, error) {
	request, err := http.NewRequest("GET", fmt.Sprintf("%s/state", br.URL), nil)
	if err != nil {
		return domain.State{}, err
	}
	query := request.URL.Query()
	query.Set("format", "both")
	request.URL.RawQuery = query.Encode()

	response, err := br.client.Do(request)
	if err != nil {
		return domain.State{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		return domain.State{}, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return domain.State{}, err
	}

	var result StateResult
	err = json.Unmarshal(body, &result)
	if err != nil {
		return domain.State{}, err
	}

	state := domain.State{
		Accumulator: result.Res.Accumulator,
		Values:      result.Res.Values,
	}

	return state, nil
}
