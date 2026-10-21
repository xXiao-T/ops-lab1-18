package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthz(t *testing.T) {
	rec := httptest.NewRecorder()
	newMux().ServeHTTP(rec, httptest.NewRequest("GET", "/healthz", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("/healthz 返回 %d，期望 200", rec.Code)
	}
}

func TestVersionEndpoint(t *testing.T) {
	rec := httptest.NewRecorder()
	newMux().ServeHTTP(rec, httptest.NewRequest("GET", "/version", nil))
	var got versionInfo
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("/version 不是合法 JSON：%v", err)
	}
	if got.Version == "" || got.Commit == "" || got.BuildNo == "" {
		t.Errorf("/version 缺少字段：%+v", got)
	}
}

func TestConflictEndpoint(t *testing.T) {
	post := func(body string) *httptest.ResponseRecorder {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/api/conflict", strings.NewReader(body))
		newMux().ServeHTTP(rec, req)
		return rec
	}

	// 单周与双周在同一时段上课，不冲突。
	rec := post(`{"a":{"id":"A","weekday":1,"start":1,"end":2,"weeks":"odd"},"b":{"id":"B","weekday":1,"start":1,"end":2,"weeks":"even"}}`)
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), `"conflict":false`) {
		t.Errorf("单双周同一时段：状态 %d，响应 %q", rec.Code, rec.Body.String())
	}

	// 每周上课的两个班同一时段，冲突。
	rec = post(`{"a":{"id":"A","weekday":1,"start":1,"end":2,"weeks":"all"},"b":{"id":"B","weekday":1,"start":2,"end":3,"weeks":"all"}}`)
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), `"conflict":true`) {
		t.Errorf("节次重叠：状态 %d，响应 %q", rec.Code, rec.Body.String())
	}

	if rec := post(`{"a":`); rec.Code != http.StatusBadRequest {
		t.Errorf("坏 JSON 应返回 400，实际 %d", rec.Code)
	}
	if rec := post(`{"a":{"weekday":9,"start":1,"end":2,"weeks":"all"},"b":{"weekday":1,"start":1,"end":2,"weeks":"all"}}`); rec.Code != http.StatusBadRequest {
		t.Errorf("不合法时段应返回 400，实际 %d", rec.Code)
	}
}

// 两个教学班都必须校验；边界合法值也应能通过 HTTP 接口。
func TestConflictEndpointValidation(t *testing.T) {
	valid := Section{Weekday: 7, Start: 1, End: 1, Weeks: AllWeeks}
	type testCase struct {
		name string
		a, b Section
		want int
	}
	cases := []testCase{
		{"周日单节课合法且冲突", valid, valid, http.StatusOK},
	}
	invalid := []struct {
		name    string
		section Section
	}{
		{"星期小于下界", Section{Weekday: 0, Start: 1, End: 1, Weeks: AllWeeks}},
		{"星期大于上界", Section{Weekday: 8, Start: 1, End: 1, Weeks: AllWeeks}},
		{"起始节次为零", Section{Weekday: 7, Start: 0, End: 1, Weeks: AllWeeks}},
		{"结束节次早于开始", Section{Weekday: 7, Start: 2, End: 1, Weeks: AllWeeks}},
		{"未知周次", Section{Weekday: 7, Start: 1, End: 1, Weeks: "weekly"}},
		{"空周次", Section{Weekday: 7, Start: 1, End: 1}},
	}
	for _, c := range invalid {
		cases = append(cases,
			testCase{"A/" + c.name, c.section, valid, http.StatusBadRequest},
			testCase{"B/" + c.name, valid, c.section, http.StatusBadRequest},
		)
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			body, err := json.Marshal(map[string]Section{"a": c.a, "b": c.b})
			if err != nil {
				t.Fatal(err)
			}
			rec := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/conflict", bytes.NewReader(body))
			newMux().ServeHTTP(rec, req)
			if rec.Code != c.want {
				t.Fatalf("状态 = %d，期望 %d；响应 %q", rec.Code, c.want, rec.Body.String())
			}
			if c.want == http.StatusOK {
				var got struct {
					Conflict bool `json:"conflict"`
				}
				if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
					t.Fatalf("响应不是合法 JSON：%v", err)
				}
				if !got.Conflict {
					t.Error("同一周日的同一节课应当冲突")
				}
			}
		})
	}
}
