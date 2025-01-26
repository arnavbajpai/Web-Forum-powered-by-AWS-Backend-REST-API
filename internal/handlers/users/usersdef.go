package users

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/arnavbajpai/web-forum-project/internal/api"
	"github.com/arnavbajpai/web-forum-project/internal/dataaccess"
	"github.com/arnavbajpai/web-forum-project/internal/database"
	"github.com/pkg/errors"
)

const (
	ListUsers = "users.HandleList"

	SuccessfulListUsersMessage = "Successfully listed users"
	ErrRetrieveDatabase        = "Failed to retrieve database in %s"
	ErrRetrieveUsers           = "Failed to retrieve users in %s"
	ErrEncodeView              = "Failed to retrieve users in %s"
)

func HandleUserList(w http.ResponseWriter, r *http.Request) (*api.Response, error) {

	users, err := dataaccess.FindUser(database.DBCon, "", 0)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf(ErrRetrieveUsers, ListUsers))
	}

	data, err := json.Marshal(users)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf(ErrEncodeView, ListUsers))
	}

	return &api.Response{
		Payload: api.Payload{
			Data: data,
		},
		Messages: []string{SuccessfulListUsersMessage},
	}, nil
}
