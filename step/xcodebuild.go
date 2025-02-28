package step

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
)

func CollectTests(testProductsPath, destination string) ([]string, error) {
	tmpDir, err := CreateTempFolder()
	if err != nil {
		return nil, err
	}

	testOutput := filepath.Join(tmpDir, "result.txt")
	//cmd := exec.Command("xcodebuild", "test-without-building", "-enumerate-tests", "-test-enumeration-format", "json", "-test-enumeration-style", "flat", "-test-enumeration-output-path", testOutput, "-testProductsPath", testProductsPath, "-destination", destination)
	cmd := exec.Command("xcodebuild", "test-without-building", "-enumerate-tests", "-test-enumeration-format", "json", "-test-enumeration-style", "flat", "-test-enumeration-output-path", testOutput, "-xctestrun", testProductsPath, "-destination", destination)

	if err := Execute(cmd); err != nil {
		return nil, err
	}

	bytes, err := os.ReadFile(testOutput)
	if err != nil {
		return nil, err
	}

	//outputString := string(bytes)
	//fmt.Println(outputString)

	type testData struct {
		Values []struct {
			Tests []struct {
				Identifier string `json:"identifier"`
			} `json:"enabledTests"`
		} `json:"values"`
	}

	var data testData
	if err := json.Unmarshal(bytes, &data); err != nil {
		return nil, err
	}

	var tests []string
	for _, value := range data.Values {
		for _, test := range value.Tests {
			tests = append(tests, test.Identifier)
		}
	}

	if err := os.Remove(testOutput); err != nil {
		return nil, err
	}

	return tests, nil
}
