package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/nozomi-nishinohara/jwt_validation/domain/model"
	"github.com/nozomi-nishinohara/jwt_validation/domain/repository"
)

func savePemCert(repo repository.Cache) error {
	var res *http.Response
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
	for _, oauth := range model.GetSetting().Oauths {
		res, err := http.Get(oauth.JwkSetUri)
		if err != nil {
			if strings.Contains(oauth.Domain, "localhost") {
				return nil
			}
			return err
		}
		jwks := model.JWKS{}
		err = json.NewDecoder(res.Body).Decode(&jwks)
		if err != nil {
			return err
		}
		if res != nil {
			res.Body.Close()
		}
		ctx := context.TODO()
		for _, jwk := range jwks.Keys {
			if err := repo.Save(ctx, jwk.Kid, jwk); err != nil {
				return err
			}
		}
	}
	return nil
}
