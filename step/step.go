package step

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/bitrise-io/go-steputils/v2/export"
	"github.com/bitrise-io/go-steputils/v2/stepconf"
	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/bitrise-io/go-utils/v2/log"
)

const (
	testShardsDirKey = "BITRISE_TEST_SHARDS_PATH"
)

type Input struct {
	ProductPath string `env:"product_path,required"`
	ShardCount  int    `env:"shard_count,required"`
	Destination string `env:"destination,required"`
	Verbose     bool   `env:"verbose,opt[true,false]"`
}

type Config struct {
	ProductPath string
	ShardCount  int
	Destination string
}

type Result struct {
	TestShardsDir string
}

type Step struct {
	envRepository env.Repository
	inputParser   stepconf.InputParser
	exporter      export.Exporter
	logger        log.Logger
}

func NewStep(
	envRepository env.Repository,
	inputParser stepconf.InputParser,
	exporter export.Exporter,
	logger log.Logger,
) Step {
	return Step{
		envRepository: envRepository,
		inputParser:   inputParser,
		exporter:      exporter,
		logger:        logger,
	}
}

func (s *Step) ProcessConfig() (*Config, error) {
	var input Input
	err := s.inputParser.Parse(&input)
	if err != nil {
		return &Config{}, err
	}

	stepconf.Print(input)
	s.logger.EnableDebugLog(input.Verbose)

	return &Config{
		ProductPath: input.ProductPath,
		ShardCount:  input.ShardCount,
		Destination: input.Destination,
	}, nil

	//return &Config{
	//	ProductPath: "/Users/szabi/Developer/misc/ManyTests/test-products.xctestproducts",
	//	ShardCount:  5,
	//	Destination: "platform=iOS Simulator,name=iPhone 16 Pro Max,OS=latest",
	//}, nil
}

func (s *Step) Run(config Config) (Result, error) {
	tests, err := CollectTests(config.ProductPath, config.Destination)
	if err != nil {
		return Result{}, err
	}

	shards := shardAlphabetically(tests, config.ShardCount)
	if len(shards) == 0 {
		return Result{}, fmt.Errorf("no tests found in %s", config.ProductPath)
	}

	shardFolder := "/Users/szabi/Developer/misc/ManyTests/test-output"
	//shardFolder, err := CreateTempFolder()
	//if err != nil {
	//	return Result{}, err
	//}

	for i, shard := range shards {
		shardPath := filepath.Join(shardFolder, fmt.Sprintf("%d", i))

		content := ""
		for _, test := range shard {
			content += fmt.Sprintf("%s\n", test)
		}

		if err := os.WriteFile(shardPath, []byte(content), 0644); err != nil {
			return Result{}, err
		}
	}

	return Result{
		TestShardsDir: shardFolder,
	}, nil
}

func (s *Step) Export(result Result) error {
	return s.exporter.ExportOutput(testShardsDirKey, result.TestShardsDir)
}

func shardAlphabetically(tests []string, shards int) [][]string {
	slices.Sort(tests)

	buckets := make([][]string, shards)
	bucketSize := (len(tests) + shards - 1) / shards

	for i, test := range tests {
		bucketIndex := i / bucketSize
		buckets[bucketIndex] = append(buckets[bucketIndex], test)
	}

	return buckets
}
