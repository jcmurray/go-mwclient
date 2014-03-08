package mwclient

import (
	"bytes"
	"fmt"
	"net/url"
	"sort"
)

// GetPageID gets the pageid of a page specified by its name.
func (w *Client) GetPageID(pageName string) (string, error) {
	params := url.Values{
		"action": {"query"},
		"prop":   {"info"},
		"titles": {pageName},
	}

	resp, err := w.Get(params)
	if err != nil {
		return "", err
	}

	pageMap, err := resp.GetPath("query", "pages").Map()
	if err != nil {
		return "", err
	}

	var id string
	for k := range pageMap {
		// There should only be one item in the map.
		id = k
	}

	if id == "-1" {
		return "", fmt.Errorf("page '%s' not found", pageName)
	}
	return id, nil
}

// URLEncode is a slightly modified version of Values.Encode() from net/url.
// It encodes url.Values into URL encoded form, sorted by key, with the exception
// of the key "token", which will be appended to the end instead of being subject
// to regular sorting. This is done because that's what the MediaWiki API wants.
func urlEncode(v url.Values) string {
	if v == nil {
		return ""
	}
	var buf bytes.Buffer
	keys := make([]string, 0, len(v))
	for k := range v {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	token := false
	for _, k := range keys {
		if k == "token" {
			token = true
			continue
		}
		vs := v[k]
		prefix := url.QueryEscape(k) + "="
		for _, v := range vs {
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}
			buf.WriteString(prefix)
			buf.WriteString(url.QueryEscape(v))
		}
	}
	if token {
		buf.WriteString("&" + url.QueryEscape("token") + "=" + url.QueryEscape(v["token"][0]))
	}
	return buf.String()
}
