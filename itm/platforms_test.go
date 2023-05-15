package itm

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestErrorIssuingPostOnCreatePlatform(t *testing.T) {
	category := map[string]interface{}{"id": float64(1)}
	radar := map[string]interface{}{"usePublicData": true}
	sonar := map[string]interface{}{"enabled": false}
	fakeClient := newFakeHTTPClient(
		fakeRoundTripper{
			resp: nil,
			err: &someError{
				errorString: "foo",
			},
		})
	Client, _ := NewClient(HTTPClient(fakeClient))
	createOps := PlatformOpts{
		Name:        "foo",
		DisplayName: "foo",
		Category:    category,
		RadarOpts:   radar,
		SonarOpts:   sonar,
		Description: "foo description",
	}
	platform, err := Client.Platform.Create(&createOps)
	if platform != nil {
		t.Error("Expected nil result")
	}
	expectedError := `Post "https://itm.cloud.com:443/api/v2/config/platforms.json": foo`
	if expectedError != err.Error() {
		t.Errorf("Unexpected error.\nExpected: %s.\nGot: %s", expectedError, err.Error())
	}
}

func TestErrorIssuingPutOnUpdatePlatform(t *testing.T) {
	category := map[string]interface{}{"id": float64(1)}
	radar := map[string]interface{}{"usePublicData": true}
	sonar := map[string]interface{}{"enabled": false}
	fakeClient := newFakeHTTPClient(
		fakeRoundTripper{
			resp: nil,
			err: &someError{
				errorString: "foo",
			},
		})
	Client, _ := NewClient(HTTPClient(fakeClient))
	updateOpts := PlatformOpts{
		Name:        "foo",
		DisplayName: "foo",
		Category:    category,
		RadarOpts:   radar,
		SonarOpts:   sonar,
		Description: "foo description",
	}
	platform, err := Client.Platform.Update(123, &updateOpts)
	if platform != nil {
		t.Error("Expected nil result")
	}
	expectedError := `Put "https://itm.cloud.com:443/api/v2/config/platforms.json/123": foo`
	if expectedError != err.Error() {
		t.Errorf("Unexpected error.\nExpected: %s.\nGot: %s", expectedError, err.Error())
	}
}

func TestErrorIssuingGetPlatform(t *testing.T) {
	fakeClient := newFakeHTTPClient(
		fakeRoundTripper{
			resp: nil,
			err: &someError{
				errorString: "foo",
			},
		})
	testClient, _ := NewClient(HTTPClient(fakeClient))
	platform, err := testClient.Platform.Get(123)
	if platform != nil {
		t.Error("Expected nil result")
	}
	expectedError := `Get "https://itm.cloud.com:443/api/v2/config/platforms.json/123": foo`
	if expectedError != err.Error() {
		t.Errorf("Unexpected error.\nExpected: %s.\nGot: %s", expectedError, err.Error())
	}
}

func TestPlatformCreate(t *testing.T) {
	teardown := setup()
	defer teardown()
	category := map[string]interface{}{"id": float64(1)}
	radar := map[string]interface{}{"usePublicData": true}
	sonar := map[string]interface{}{"enabled": false}
	mux.HandleFunc("/v2/config/platforms.json", func(w http.ResponseWriter, r *http.Request) {
		var parsedBody map[string]interface{}
		expectedRequestData := map[string]interface{}{
			"name":        "foo",
			"displayName": "foo",
			"category":    category,
			"radarConfig": radar,
			"sonarConfig": sonar,
			"intendedUse": "foo description",
		}
		responseBodyObj := Platform{
			Id:          123,
			Name:        "foo",
			DisplayName: "foo",
			Category:    category,
			RadarOpts:   radar,
			SonarOpts:   sonar,
			Description: "foo description",
		}
		err := json.NewDecoder(r.Body).Decode(&parsedBody)
		if err != nil {
			t.Fatalf("JSON decoding error: %v", err)
		}
		if !compareMaps(expectedRequestData, parsedBody) {
			t.Error(unexpectedValueString("Request body", expectedRequestData, parsedBody))
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		responseBody, _ := json.Marshal(responseBodyObj)
		fmt.Fprint(w, string(responseBody))
	})
	createOps := PlatformOpts{
		Name:        "foo",
		DisplayName: "foo",
		Category:    category,
		RadarOpts:   radar,
		SonarOpts:   sonar,
		Description: "foo description",
	}
	platform, err := client.Platform.Create(&createOps)
	if err != nil {
		t.Error(err)
	}
	if err := testValues("id", 123, platform.Id); err != nil {
		t.Error(err)
	}
	if err := testValues("platform alias", "foo", platform.Name); err != nil {
		t.Error(err)
	}
	if err := testValues("platform display name", "foo", platform.DisplayName); err != nil {
		t.Error(err)
	}
	if op := reflect.DeepEqual(category, platform.Category); !op {
		t.Error(unexpectedValueString("platform category", category, platform.Category))
	}
	if op := reflect.DeepEqual(radar, platform.RadarOpts); !op {
		t.Error(unexpectedValueString("platform radar options", radar, platform.RadarOpts))
	}
	if op := reflect.DeepEqual(sonar, platform.SonarOpts); !op {
		t.Error(unexpectedValueString("platform sonar options", sonar, platform.SonarOpts))
	}
	if err := testValues("description", "foo description", platform.Description); err != nil {
		t.Error(err)
	}
}

func TestPlatformUpdate(t *testing.T) {
	teardown := setup()
	defer teardown()
	category := map[string]interface{}{"id": float64(1)}
	radar := map[string]interface{}{"usePublicData": true}
	sonar := map[string]interface{}{"enabled": false}
	mux.HandleFunc("/v2/config/platforms.json/123", func(w http.ResponseWriter, r *http.Request) {
		var parsedBody map[string]interface{}
		expectedRequestData := map[string]interface{}{
			"name":        "updated_foo_name",
			"displayName": "updated_foo_name",
			"category":    category,
			"radarConfig": radar,
			"sonarConfig": sonar,
			"intendedUse": "updated foo description",
		}
		responseBodyObj := Platform{
			Id:          123,
			DisplayName: "updated_foo_name",
			Name:        "updated_foo_name",
			Category:    category,
			RadarOpts:   radar,
			SonarOpts:   sonar,
			Description: "updated foo description",
		}
		err := json.NewDecoder(r.Body).Decode(&parsedBody)
		if err != nil {
			t.Fatalf("JSON decoding error: %v", err)
		}
		if !compareMaps(expectedRequestData, parsedBody) {
			t.Error(unexpectedValueString("Request body", expectedRequestData, parsedBody))
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		responseBody, _ := json.Marshal(responseBodyObj)
		fmt.Fprint(w, string(responseBody))
	})
	updateOps := PlatformOpts{
		Name:        "updated_foo_name",
		DisplayName: "updated_foo_name",
		Category:    category,
		RadarOpts:   radar,
		SonarOpts:   sonar,
		Description: "updated foo description",
	}
	platform, err := client.Platform.Update(123, &updateOps)
	if err != nil {
		t.Error(err)
	}
	if err := testValues("id", 123, platform.Id); err != nil {
		t.Error(err)
	}
	if err := testValues("platform alias", "updated_foo_name", platform.Name); err != nil {
		t.Error(err)
	}
	if err := testValues("platform display name", "updated_foo_name", platform.DisplayName); err != nil {
		t.Error(err)
	}
	if op := reflect.DeepEqual(category, platform.Category); !op {
		t.Error(unexpectedValueString("platform category", category, platform.Category))
	}
	if op := reflect.DeepEqual(radar, platform.RadarOpts); !op {
		t.Error(unexpectedValueString("platform radar options", radar, platform.RadarOpts))
	}
	if op := reflect.DeepEqual(sonar, platform.SonarOpts); !op {
		t.Error(unexpectedValueString("platform sonar options", sonar, platform.SonarOpts))
	}
	if err := testValues("description", "updated foo description", platform.Description); err != nil {
		t.Error(err)
	}
}

func TestPlatformGet(t *testing.T) {
	teardown := setup()
	defer teardown()
	category := map[string]interface{}{"id": float64(1)}
	radar := map[string]interface{}{"usePublicData": true}
	sonar := map[string]interface{}{"enabled": false}
	mux.HandleFunc("/v2/config/platforms.json/123", func(w http.ResponseWriter, r *http.Request) {
		responseBodyObj := Platform{
			Id:          123,
			DisplayName: "foo",
			Name:        "foo",
			Category:    category,
			RadarOpts:   radar,
			SonarOpts:   sonar,
			Description: "foo description",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		responseBody, _ := json.Marshal(responseBodyObj)
		fmt.Fprint(w, string(responseBody))
	})
	platform, err := client.Platform.Get(123)
	if err != nil {
		t.Error(err)
	}
	if err := testValues("id", 123, platform.Id); err != nil {
		t.Error(err)
	}
	if err := testValues("platform alias", "foo", platform.Name); err != nil {
		t.Error(err)
	}
	if err := testValues("platform display name", "foo", platform.DisplayName); err != nil {
		t.Error(err)
	}
	if op := reflect.DeepEqual(category, platform.Category); !op {
		t.Error(unexpectedValueString("platform category", category, platform.Category))
	}
	if op := reflect.DeepEqual(radar, platform.RadarOpts); !op {
		t.Error(unexpectedValueString("platform radar options", radar, platform.RadarOpts))
	}
	if op := reflect.DeepEqual(sonar, platform.SonarOpts); !op {
		t.Error(unexpectedValueString("platform sonar options", sonar, platform.SonarOpts))
	}
	if err := testValues("description", "foo description", platform.Description); err != nil {
		t.Error(err)
	}
}

func TestPlatformDelete(t *testing.T) {
	teardown := setup()
	defer teardown()
	mux.HandleFunc("/v2/config/platforms.json/123", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	err := client.Platform.Delete(123)
	if err != nil {
		t.Error(err)
	}
}

func TestPlatformList(t *testing.T) {
	teardown := setup()
	defer teardown()
	category := map[string]interface{}{"id": float64(1)}
	radar := map[string]interface{}{"usePublicData": true}
	sonar := map[string]interface{}{"enabled": false}
	var plaforms []Platform
	platform1 := Platform{
		Id:          123,
		Name:        "foo",
		DisplayName: "foo",
		Category:    category,
		RadarOpts:   radar,
		SonarOpts:   sonar,
		Description: "foo description",
	}
	platform2 := Platform{
		Id:          456,
		Name:        "bar",
		DisplayName: "bar",
		Category:    category,
		RadarOpts:   radar,
		SonarOpts:   sonar,
		Description: "bar description",
	}

	plaforms = append(plaforms, platform1, platform2)
	mux.HandleFunc("/v2/config/platforms.json/", func(w http.ResponseWriter, r *http.Request) {
		responseBodyObj := plaforms
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		responseBody, _ := json.Marshal(responseBodyObj)
		fmt.Fprint(w, string(responseBody))
	})

	platformlist, err := client.Platform.List()
	if err != nil {
		t.Error(err)
	}
	for index, platform := range platformlist {
		if !reflect.DeepEqual(platform, plaforms[index]) {
			t.Error(unexpectedValueString("plaforms parameter", plaforms, plaforms[index]))
		}
	}
}
