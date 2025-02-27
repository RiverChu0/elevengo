package elevengo

import (
	"fmt"
	"github.com/deadblue/elevengo/internal/web"
	"github.com/deadblue/elevengo/internal/webapi"
	"math/rand"
	"time"
)

// Credential contains required information which upstream server uses to
// authenticate a signed-in user.
// In detail, three cookies are required: "UID", "CID", "SEID", caller can
// find them from browser cookie storage.
type Credential struct {
	UID  string
	CID  string
	SEID string
}

// UserInfo contains the basic information of a signed-in user.
type UserInfo struct {
	Id   int
	Name string
}

// CredentialImport imports credentials into agent.
func (a *Agent) CredentialImport(cr *Credential) (err error) {
	cookies := map[string]string{
		webapi.CookieNameUid:  cr.UID,
		webapi.CookieNameCid:  cr.CID,
		webapi.CookieNameSeid: cr.SEID,
	}
	a.wc.ImportCookies(cookies, webapi.CookieDomain115, webapi.CookieDomainAnxia)
	if !a.LoginCheck() {
		err = webapi.ErrCredentialInvalid
	}
	return
}

// CredentialExport exports credentials for future-use.
func (a *Agent) CredentialExport(cr *Credential) {
	cookies := a.wc.ExportCookies(webapi.CookieUrl)
	cr.UID = cookies[webapi.CookieNameUid]
	cr.CID = cookies[webapi.CookieNameCid]
	cr.SEID = cookies[webapi.CookieNameSeid]
}

func (a *Agent) LoginCheck() bool {
	qs := web.Params{}.WithNowMilli("_")
	resp := &webapi.LoginCheckResponse{}
	if err := a.wc.CallJsonApi(webapi.ApiLoginCheck, qs, nil, resp); err != nil {
		return false
	}
	data := &webapi.LoginCheckData{}
	if err := resp.Decode(data); err == nil {
		a.uid = data.UserId
		return true
	} else {
		return false
	}
}

// UserGet retrieves user information from cloud.
func (a *Agent) UserGet(info *UserInfo) (err error) {
	cb := fmt.Sprintf("jQuery%d_%d", rand.Uint64(), time.Now().Unix())
	qs := web.Params{}.
		With("callback", cb).
		WithNow("_")
	resp := webapi.BasicResponse{}
	if err = a.wc.CallJsonpApi(webapi.ApiUserInfo, qs, &resp); err != nil {
		return
	}
	result := webapi.UserInfoData{}
	if err = resp.Decode(&result); err == nil {
		info.Id = result.UserId
		info.Name = result.UserName
	}
	return
}
