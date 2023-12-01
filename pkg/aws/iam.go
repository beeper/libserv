package aws

import (
	"net/http"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	"github.com/rs/zerolog/log"
	"sigs.k8s.io/aws-iam-authenticator/pkg/token"
)

type IAMAuthRoundTripper struct {
	clusterID string
	lock      sync.Mutex
	token     token.Token
	generator token.Generator
	stsAPI    stsiface.STSAPI
	rt        http.RoundTripper
}

func NewIAMAuthRoundTripper(clusterID string, rt http.RoundTripper) (http.RoundTripper, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}

	gen, err := token.NewGenerator(true, false)
	if err != nil {
		return nil, err
	}

	return &IAMAuthRoundTripper{
		clusterID: clusterID,
		stsAPI:    sts.New(sess),
		generator: gen,
		rt:        rt,
	}, nil
}

func (rt *IAMAuthRoundTripper) RoundTrip(req *http.Request) (r *http.Response, err error) {
	if len(req.Header.Get("Authorization")) != 0 {
		return rt.rt.RoundTrip(req)
	}

	rt.lock.Lock()
	defer rt.lock.Unlock()

	// token renewal
	if rt.token.Token == "" || time.Now().Add(4*time.Minute).After(rt.token.Expiration) {
		log := log.Ctx(req.Context())
		log.Info().Msgf("Renewing STS token for cluster %s, old one expires at %s", rt.clusterID, rt.token.Expiration.String())

		rt.token, err = rt.generator.GetWithSTS(rt.clusterID, rt.stsAPI)
		if err != nil {
			log.Err(err).Msg("Failed to renew STS token")
			return
		}

		log.Info().Msgf("New STS token for cluster %s expires at %s (in %s)", rt.clusterID, rt.token.Expiration.String(), time.Until(rt.token.Expiration).String())
	}

	req.Header.Set("Authorization", "Bearer "+rt.token.Token)
	return rt.rt.RoundTrip(req)
}
