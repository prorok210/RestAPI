package server

import (
	"bytes"
	"errors"
	"strconv"
	"testing"
	"time"
)

/*
http.go testing
*/
func TestParseRequest(t *testing.T) {
	testCases := []struct {
		input          []byte
		expectedMethod string
		expectedUrl    string
		expectedError  bool
	}{
		{
			[]byte("GET / HTTP/1.1\r\nHost: localhost:8080\r\n\r\n"),
			"GET", "/", false,
		},
		{
			[]byte("POST / HTTP/1.1\r\nHost: localhost:8080\r\nContent-Type: application/json\r\nContent-Length: 15\r\n\r\n{\"key\": \"value\"}"),
			"POST", "/", false,
		},
		{
			[]byte("PUT / HTTP/1.1\r\nHost: localhost:8080\r\nContent-Type: application/json\r\nContent-Length: 15\r\n\r\n"),
			"PUT", "/", false,
		},
		{
			[]byte("Host: localhost:8080\r\nContent-Type: application/json\r\nContent-Length: 15\r\n\r\n"),
			"", "", true,
		},
		{
			[]byte(""),
			"", "", true,
		},
	}

	for i, testCase := range testCases {
		testRequest := new(HttpRequest)
		err := testRequest.ParseRequest(testCase.input)
		if err != nil && !testCase.expectedError {
			t.Errorf("Unexpected error in %d test case: %s", i, err)
		}
		if testRequest.Method != testCase.expectedMethod {
			t.Errorf("Unexpected method in %d test case: %s != %s", i, testRequest.Method, testCase.expectedMethod)
		}
		if testRequest.Url != testCase.expectedUrl {
			t.Errorf("Unexpected url in %d test case: %s != %s", i, testRequest.Url, testCase.expectedUrl)
		}
	}
}

func TestSetHeader(t *testing.T) {
	testCases := []struct {
		key   string
		value string
	}{
		{
			"Content-Type", "application/json",
		},
		{
			"Content-Length", "15",
		},
		{
			"", "",
		},
	}

	for i, testCase := range testCases {
		testResponse := new(HttpResponse)
		testResponse.SetHeader(testCase.key, testCase.value)
		if testResponse.Headers[testCase.key] != testCase.value {
			t.Errorf("Unexpected value in %d test case: %s != %s", i, testResponse.Headers[testCase.key], testCase.value)
		}
	}
}

func TestParseFormData(t *testing.T) {
	boundary := "--------------------------477699956037780681100607"
	type FileInfo struct{
		FileName string
		FileData []byte
	}

	testCases := []struct {
		name           string
		request        HttpRequest
		expectedError  error
		expectedFields map[string]string
		expectedFiles  map[string][]FileInfo
	}{
		{
			name: "Valid form-data with fields and file",
			request: HttpRequest{
				Headers: map[string]string{
					"Content-Type": "multipart/form-data; boundary=" + boundary,
				},
				Body: "--" + boundary + "\r\n" +
					"Content-Disposition: form-data; name=\"field1\"\r\n\r\n" +
					"value1\r\n" +
					"--" + boundary + "\r\n" +
					"Content-Disposition: form-data; name=\"field2\"\r\n\r\n" +
					"value2\r\n" +
					"--" + boundary + "\r\n" +
					"Content-Disposition: form-data; name=\"file1\"; filename=\"test.txt\"\r\n" +
					"Content-Type: text/plain\r\n\r\n" +
					"file content\r\n" +
					"--" + boundary + "--\r\n",
			},
			expectedError: nil,
			expectedFields: map[string]string{
				"field1": "value1",
				"field2": "value2",
			},
			expectedFiles: map[string][]FileInfo{
				"file1": {
					{"test.txt", []byte("file content")},
				},
			},
		},
		{
			name: "Invalid Content-Type",
			request: HttpRequest{
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
				Body: `{"key": "value"}`,
			},
			expectedError: errors.New("invalid Content-Type"),
		},
		{
			name: "Missing boundary",
			request: HttpRequest{
				Headers: map[string]string{
					"Content-Type": "multipart/form-data",
				},
				Body: "some body",
			},
			expectedError: errors.New("boundary not found in Content-Type"),
		},
		{
			name: "Empty form-data",
			request: HttpRequest{
				Headers: map[string]string{
					"Content-Type": "multipart/form-data; boundary=" + boundary,
				},
				Body: "--" + boundary + "--\r\n",
			},
			expectedError:  nil,
			expectedFields: map[string]string{},
			expectedFiles:  map[string][]FileInfo{},
		},
		{
			name: "Form-data with special characters",
			request: HttpRequest{
				Headers: map[string]string{
					"Content-Type": "multipart/form-data; boundary=" + boundary,
				},
				Body: "--" + boundary + "\r\n" +
					"Content-Disposition: form-data; name=\"special_field\"\r\n\r\n" +
					"value with\r\nnew line and €\r\n" +
					"--" + boundary + "\r\n" +
					"Content-Disposition: form-data; name=\"file2\"; filename=\"специальный файл.txt\"\r\n" +
					"Content-Type: text/plain\r\n\r\n" +
					"содержимое файла\r\n" +
					"--" + boundary + "--\r\n",
			},
			expectedError: nil,
			expectedFields: map[string]string{
				"special_field": "value with\r\nnew line and €",
			},
			expectedFiles: map[string][]FileInfo{
				"file2": {
					{
						FileName: "специальный файл.txt",
						FileData: []byte("содержимое файла"),
					},
				},
			},
		},
		{
			name: "Invalid form-data structure",
			request: HttpRequest{
				Headers: map[string]string{
					"Content-Type": "multipart/form-data; boundary=" + boundary,
				},
				Body: "--" + boundary + "\r\n" +
					"Invalid-Header: some value\r\n\r\n" +
					"some content\r\n" +
					"--" + boundary + "--\r\n",
			},
			expectedError: errors.New("invalid Content-Disposition: name not found"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.request.ParseFormData()

			if err != nil && tc.expectedError == nil {
				t.Errorf("Unexpected error: %v", err)
			} else if tc.expectedError != nil {
				if err == nil {
					t.Errorf("Expected error %v, but got nil", tc.expectedError)
				} else if err.Error() != tc.expectedError.Error() {
					t.Errorf("Expected error %v, but got %v", tc.expectedError, err)
				}
			}

			if tc.expectedError != nil {
				return
			}

			if len(tc.request.FormData.Fields) != len(tc.expectedFields) {
				t.Errorf("Expected %d fields, but got %d", len(tc.expectedFields), len(tc.request.FormData.Fields))
			}
			for key, expectedValue := range tc.expectedFields {
				if value, ok := tc.request.FormData.Fields[key]; !ok {
					t.Errorf("Expected field %s not found", key)
				} else if value != expectedValue {
					t.Errorf("Field %s: expected %s, got %s", key, expectedValue, value)
				}
			}

			if len(tc.request.FormData.Files) != len(tc.expectedFiles) {
				t.Errorf("Expected %d files, but got %d", len(tc.expectedFiles), len(tc.request.FormData.Files))
			}
			for key, expectedFiles := range tc.expectedFiles {
				actualFiles, ok := tc.request.FormData.Files[key]
				if !ok {
					t.Errorf("No files found for key %s", key)
					continue
				}
			
				// Проверяем, что количество файлов совпадает
				if len(actualFiles) != len(expectedFiles) {
					t.Errorf("File %s: expected %d files, got %d", key, len(expectedFiles), len(actualFiles))
					continue
				}
			
				// Сравниваем каждый файл в массиве
				for i := range expectedFiles {
					expectedFile := expectedFiles[i]
					file := actualFiles[i]
			
					// Проверяем имя файла
					if file.FileName != expectedFile.FileName {
						t.Errorf("File %s [%d]: expected name %s, got %s", key, i, expectedFile.FileName, file.FileName)
					}
			
					// Проверяем содержимое файла
					if !bytes.Equal(file.FileData, expectedFile.FileData) {
						t.Errorf("File %s [%d]: content mismatch", key, i)
					}
				}
			}
		})
	}
}


func TestSerialize(t *testing.T) {
	testCases := []struct {
		data           interface{}
		expectedResult string
		expectedError  error
	}{
		{
			"test", `"test"`, nil,
		},
		{
			123, `123`, nil,
		},
		{
			map[string]interface{}{
				"key":  "value",
				"key2": "value2",
			},
			`{"key":"value","key2":"value2"}`,
			nil,
		},
		{
			nil, "", errors.New("Data is nil"),
		},
	}

	for i, testCase := range testCases {
		testResponse := new(HttpResponse)
		err := testResponse.Serialize(testCase.data)

		if err != nil {
			if testCase.expectedError == nil {
				t.Errorf("Test case %d: Unexpected error: %v", i, err)
			} else if err.Error() != testCase.expectedError.Error() {
				t.Errorf("Test case %d: Expected error '%v', but got '%v'", i, testCase.expectedError, err)
			}
		} else if testCase.expectedError != nil {
			t.Errorf("Test case %d: Expected error but got nil", i)
		}

		if err == nil && testResponse.Body != testCase.expectedResult {
			t.Errorf("Test case %d: Unexpected result: got '%s', want '%s'", i, testResponse.Body, testCase.expectedResult)
		}
	}
}

/*
server.go testing
*/
func TestCreateServer(t *testing.T) {
	testCases := []struct {
		mainApplication RequestHandler
		expectedError   bool
	}{
		{
			nil, true,
		},
		{
			func(request *HttpRequest) ([]byte, error) {
				return []byte(`{"Status": "OK"}`), nil
			}, false,
		},
	}

	for i, testCase := range testCases {
		_, err := CreateServer(testCase.mainApplication)
		if err != nil && !testCase.expectedError {
			t.Errorf("Unexpected error in %d test case: %s", i, err)
		}
	}
}

func TestStartServer(t *testing.T) {
	testCases := []struct {
		handleApp     RequestHandler
		expectedError bool
	}{
		{
			nil, true,
		},
		{
			func(request *HttpRequest) ([]byte, error) {
				return []byte(`{"Status": "OK"}`), nil
			}, false,
		},
	}

	for i, testCase := range testCases {
		server, _ := CreateServer(testCase.handleApp)
		err := server.Start()
		if err != nil && !testCase.expectedError {
			t.Errorf("Unexpected error in %d test case: %s", i, err)
		}
		server.Stop()
	}
}

/*
middlewares.go testing
*/
func TestIsAllowedHostMiddleware(t *testing.T) {
	if !IS_ALLOWED_HOSTS {
		t.Skip("Allowed hosts are not set")
	}
	testCases := []struct {
		clientAddr string
		expected   bool
	}{
		{
			"localhost:8080", true,
		},
		{
			"Это не адрес", false,
		},
		{
			"", false,
		},
		{
			"127.0.0.1:8080", true,
		},
		{
			"77.232.37.23:8080", true,
		},
		{
			"[::1]:8888", true,
		},
	}

	for i, testCase := range testCases {
		result := isAllowedHostMiddleware(testCase.clientAddr)
		if result != testCase.expected {
			t.Errorf("Unexpected result in %d test case: %t != %t", i, result, testCase.expected)
		}
	}
}

func TestReqMiddleware(t *testing.T) {
	if !REQ_MIDDLEWARE {
		t.Skip("Request middleware is disabled")
	}
	connMock := &ConnMock{
		WriteFunc: func(b []byte) (n int, err error) {
			return len(b), nil
		},
		CloseFunc: func() error {
			return nil
		},
	}

	testCases := []struct {
		request       HttpRequest
		clientConn    *ConnMock
		expectedError error
	}{
		// Test case 0: Valid GET request
		{
			HttpRequest{
				Method: "GET",
			}, connMock, nil,
		},
		// Test case 1: Valid POST request
		{
			HttpRequest{
				Method: "POST",
				Headers: map[string]string{
					"Content-Type":   "application/json",
					"Content-Length": strconv.Itoa(len(`{"key": "value"}`)),
				},
				Body: `{"key": "value"}`,
			}, connMock, nil,
		},
		// Test case 2: Unsupported method
		{
			HttpRequest{
				Method: "AdsD",
			}, connMock, errors.New("Method not allowed"),
		},
		// Test case 3: Invalid Content-Length
		{
			HttpRequest{
				Method: "POST",
				Headers: map[string]string{
					"Content-Type":   "application/json",
					"Content-Length": strconv.Itoa(len(`{"key": "value"}`)),
				},
				Body: `{""}`,
			}, connMock, errors.New("Content-Length does not match body length"),
		},
		// Test case 4: Unsupported media type
		{
			HttpRequest{
				Method: "POST",
				Headers: map[string]string{
					"Content-Type": "applicatidasdas",
				},
				Body: `{"key": "value"}`,
			}, connMock, errors.New("Unsupported media type"),
		},
	}

	for i, testCase := range testCases {
		err := reqMiddleware(&testCase.request, testCase.clientConn)

		if testCase.expectedError == nil && err != nil {
			t.Errorf("Test case %d: Expected no error, but got: %s", i, err)
		} else if testCase.expectedError != nil {
			if err == nil {
				t.Errorf("Test case %d: Expected error, but got nil", i)
			} else if err.Error() != testCase.expectedError.Error() {
				t.Errorf("Test case %d: Expected error %s, but got %s", i, testCase.expectedError, err)
			}
		}
		
		if err := connMock.Close(); err != nil {
			t.Errorf("Test case %d: Failed to close connection: %s", i, err)
		}
	}
}

func TestKeepAliveMiddleware(t *testing.T) {
	if !KEEP_ALIVE {
		t.Skip("Keep-alive is disabled")
	}
	connMock := &ConnMock{
		SetDeadlineFunc: func(t time.Time) error {
			return nil
		},
		CloseFunc: func() error {
			return nil
		},
	}

	testCases := []struct {
		request       HttpRequest
		clientConn    *ConnMock
		expectedError error
	}{
		// Test case 0: Valid request
		{
			HttpRequest{
				Headers: map[string]string{
					"Connection": "keep-alive",
				},
			}, connMock, nil,
		},
		// Test case 1: Connection: close
		{
			HttpRequest{
				Headers: map[string]string{
					"Connection": "close",
				},
			}, connMock, errors.New("Connection: close"),
		},
		// Test case 2: No Connection header
		{
			HttpRequest{
				Headers: map[string]string{},
			}, connMock, nil,
		},
		// Test case 3: Invalid Connection header
		{
			HttpRequest{}, connMock, errors.New("Connection: close"),
		},
	}

	for i, testCase := range testCases {
		err := keepAliveMiddleware(&testCase.request, testCase.clientConn)

		if testCase.expectedError == nil && err != nil {
			t.Errorf("Test case %d: Expected no error, but got: %s", i, err)
		} else if testCase.expectedError != nil {
			if err == nil {
				t.Errorf("Test case %d: Expected error, but got nil", i)
			} else if err.Error() != testCase.expectedError.Error() {
				t.Errorf("Test case %d: Expected error %s, but got %s", i, testCase.expectedError, err)
			}
		}

		if err := testCase.clientConn.Close(); err != nil {
			t.Errorf("Test case %d: Failed to close connection: %s", i, err)
		}
	}
}
