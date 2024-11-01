package googleoauth

import (
	"encoding/json"
	"net/http"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/oauth"
)

func (g *GoogleOAuth) GetUser(
	token string,
) (*oauth.User, error) {
	url := g.BaseURL.String() + "/userinfo/v2/me"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errs.New(err)
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("authorization", "Bearer "+token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errs.New(err)
	}
	if res == nil {
		return nil, errs.New("response is nil")
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	if res.StatusCode < 200 && res.StatusCode >= 300 {
		jsonData := map[string]any{}
		if err := decoder.Decode(&jsonData); err != nil {
			return nil, errs.New(err)
		}
		return nil, errs.New(jsonData)
	}

	userInfo := oauth.User{}
	if err := decoder.Decode(&userInfo); err != nil {
		return nil, errs.New(err)
	}

	return &userInfo, nil
}
