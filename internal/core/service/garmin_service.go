package service

import (
	"context"
	"event-registration/internal/common"
	"event-registration/internal/common/helper"
	"event-registration/internal/core/domain"
	"fmt"
	"io"
	"net/http"

	"go.uber.org/zap"
)

type GarminService struct {
	repo   domain.GarminRepository
	logger *zap.Logger
	config *common.Config
}

func NewGarminService(repo domain.GarminRepository, logger *zap.Logger, config *common.Config) *GarminService {
	return &GarminService{
		repo:   repo,
		logger: logger,
		config: config,
	}
}

func (s *GarminService) Refresh(ctx context.Context) (res *domain.Activity, err error) {
	url := "https://connect.garmin.com/activitylist-service/activities/search/activities?limit=10&start=0"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		s.logger.Error("error_make_new_request", zap.Error(err))
		return nil, err
	}

	// Header
	req.Header.Set("accept", "application/json, text/plain, */*")
	req.Header.Set("authorization", "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6ImRpLW9hdXRoLXNpZ25lci1wcm9kLTIwMjQtcTEifQ.eyJzY29wZSI6WyJDT01NVU5JVFlfQ09VUlNFX1JFQUQiLCJHT0xGX0FQSV9SRUFEIiwiR0hTX0hJRCIsIkFUUF9SRUFEIiwiRElWRV9BUElfV1JJVEUiLCJHSFNfU0FNRCIsIklOU0lHSFRTX1JFQUQiLCJESVZFX0FQSV9SRUFEIiwiRElWRV9BUElfSU1BR0VfUFJFVklFVyIsIkNPTU1VTklUWV9DT1VSU0VfV1JJVEUiLCJDT05ORUNUX1dSSVRFIiwiRElWRV9TSEFSRURfUkVBRCIsIkdIU19SRUdJU1RSQVRJT04iLCJEVF9DTElFTlRfQU5BTFlUSUNTX1dSSVRFIiwiR09MRl9BUElfV1JJVEUiLCJJTlNJR0hUU19XUklURSIsIlBST0RVQ1RfU0VBUkNIX1JFQUQiLCJHT0xGX1NIQVJFRF9SRUFEIiwiT01UX0NBTVBBSUdOX1JFQUQiLCJPTVRfU1VCU0NSSVBUSU9OX1JFQUQiLCJDT05ORUNUX05PTl9TT0NJQUxfU0hBUkVEX1JFQUQiLCJDT05ORUNUX1JFQUQiLCJBVFBfV1JJVEUiXSwiaXNzIjoiaHR0cHM6Ly9kaWF1dGguZ2FybWluLmNvbSIsInJldm9jYXRpb25fZWxpZ2liaWxpdHkiOlsiR0xPQkFMX1NJR05PVVQiXSwiY2xpZW50X3R5cGUiOiJVTkRFRklORUQiLCJleHAiOjE3NTMxNzM0MTgsImlhdCI6MTc1MzE2OTgxOCwiZ2FybWluX2d1aWQiOiI3OTlmYzk4Zi1jYzRiLTQwZmQtOWJmYS1hZGVmMzViMDM1ZGIiLCJqdGkiOiIxNTlmM2QxMi0zMzQzLTRjYzQtYWExZS0wNDI0ZDIwMDUxNjAiLCJjbGllbnRfaWQiOiJDT05ORUNUX1dFQiIsImZncCI6IjI1Nzk4ZmVkMTYzM2Q2MTZkYjczYjA5NDQ4OTUyMmViNTMzNjBiYzE2YTBhOGRiYTRiYzljNDAxNjAyZDFjM2UifQ.DIVDyY6LzX0p0CjkSv6TfWlbdjwniTNbjWJbTxelG53OgbJ2CxGVIZ-E0pgdeR__I0sfItXGVcr1NuIMLxH718QPQsrfjFHHCWPBlpDdRD5ZUQqasXdyiEyJ-bo9K2Fkc3P6-1cpLFuTw2WKwroxCgL5soFm01onDy1-Qv21LS4S9F44frnVW_YFYa30r2iMvnEC03Bdxl_XlbsOCKopPkQ13ctt5PM9KDKY26HVxjCkt4qj0b-CAgMdSf1O0Vq30PfO4SKZyrfZCC19G_n1eTPwJ2Fphc6umTYPoh5lFWReQO_x5x1LnXyMwuP8uzV5JTZSeNgYwj0Z8pepF4u9BA")
	req.Header.Set("di-backend", "connectapi.garmin.com")

	// Cookie
	req.Header.Set("Cookie", `GarminUserPrefs=en-US; notice_behavior=none; _pk_id.6.e3f4=8a21d2a331c95d06.1751977650.; GarminNoCache=true; GARMIN-SSO-GUID=1F77EF8AF3B791E455D9F973A27762A317B8799B; _cfuvid=xIY2Qjrr306_LYPFxbswXYbCE8ulmWfG1iQNEbnDPaI-1752582969092-0.0.1.1-604800000; ADRUM=s=1752836308384&r=https://www.garmin.com/in-ID/privacy/connect/?0; __cfruid=e9a84b4b25debd928274ade6ea570ed2d3bbd041-1752993849; GARMIN-SSO=1; GARMIN-SSO-CUST-GUID=799fc98f-cc4b-40fd-9bfa-adef35b035db; GMN_TRACKABLE=1; utag_main__sn=4; _hjSessionUser_1939392=eyJpZCI6IjAzNTkwZDdiLTdmZmItNWM4ZS04MmUzLWM5MjM0OWNhYWY2ZiIsImNyZWF0ZWQiOjE3NTMwMDEyMzA3NTYsImV4aXN0aW5nIjp0cnVlfQ==; notice_behavior=none; CONSENTMGR=consent:true|ts:1753001406730; _hjSessionUser_748865=eyJpZCI6ImU1MDdlNTEwLTQ3MzUtNWJlMi1iN2M0LTViZWU1MDEzZGQ5OCIsImNyZWF0ZWQiOjE3NTMwMDE0MDY5NjksImV4aXN0aW5nIjpmYWxzZX0=; _ga=GA1.1.230768271.1751977651; _ga_1K2WNZ9N3T=GS2.1.s1753001409$o3$g1$t1753001811$j60$l0$h0; SESSIONID=NzY1YWVlY2QtMmRkNC00ODIyLThiMWEtZmMyOGQ1NTY5ZDUy; __cflb=02DiuJLbVZHipNWxN8wwnxZhF2QbAv3GYx6H5EPspNrtU; JWT_FGP=c1f1b6f1-db0b-477b-a678-5dce84455e16; JWT_WEB=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTMxMTAwMDgsImlzcyI6ImF1dGgwIiwicm9sZXMiOlsiNCIsIjciLCI4Il19.fLONTaa8nNwAeEng1QYDGG51oQcBNl9qoWxpe0X7qmA; _pk_ref.6.e3f4=["","","1753103082","https://sso.garmin.com/"]; _ga_DHVMQZ6WGH=GS2.1.s1753101385$o13$g1$t1753103088$j60$l0$h0; ADRUM_BTa=R:44|g:4fc6621b-5fec-4e6e-9ab2-f58ea0a2d24c|n:garmin_869629ee-d273-481d-b5a4-f4b0a8c4d5a3; ADRUM_BT1=R:44|i:1376314|e:5|t:1753103080057`)

	// Kirim request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		s.logger.Error("error_do_request", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	// Baca response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.Error("error_read_request", zap.Error(err))
		return nil, err
	}

	fmt.Println(string(body))

	helper.PrettyPrint(string(body), "Garmin Response =========================================")

	return res, err
}
