package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"noda/failure"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/lib/pq"
	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v3"
)

var (
	db           *sql.DB
	databasePath string
	indexPath    string
)

var verbose bool = false

func init() {
	root := findRootDirectory()
	databasePath = path.Join(root, "database")
	indexPath = path.Join(databasePath, "index.yaml")
}

type files []string

type index struct {
	init       files `yaml:"init"`
	extensions files `yaml:"extensions"`
	domains    files `yaml:"domains"`
	types      files `yaml:"types"`
	tables     files `yaml:"tables"`
	views      files `yaml:"views"`
	indexes    files `yaml:"indexes"`
	routines   files `yaml:"routines"`
	seeds      files `yaml:"seeds"`
}

// mustGetEnv tries to get an env var or exists.
func mustGetEnv(key string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if "" == value {
		log.Fatalf("could not load env var: %s", key)
	}
	return value
}

func main() {
	buf, err := os.ReadFile(indexPath)
	if err != nil {
		log.Fatalf("could not read %q: %v", indexPath, err)
	}

	var index index
	err = yaml.Unmarshal(buf, &index)
	if err != nil {
		log.Fatalf("error unmarshalling %q: %v", indexPath, err)
	}

	var (
		dbUser     = mustGetEnv("RUNNER_ROOT_DB_USER")
		dbPassword = mustGetEnv("RUNNER_ROOT_DB_PASSWORD")
		dbHost     = mustGetEnv("RUNNER_ROOT_DB_HOST")
		dbPort     = mustGetEnv("RUNNER_ROOT_DB_PORT")
		dbName     = mustGetEnv("RUNNER_ROOT_DB_NAME")
	)

	dbRootConf := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	db, err = sql.Open("postgres", dbRootConf)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	executeScriptsOf("init", index.init)

	if err = db.Close(); err != nil {
		log.Fatal(err)
	}

	dbUser = mustGetEnv("NODA_ROOT_DB_USER")
	dbPassword = mustGetEnv("NODA_ROOT_DB_PASSWORD")
	dbHost = mustGetEnv("NODA_ROOT_DB_HOST")
	dbPort = mustGetEnv("NODA_ROOT_DB_PORT")
	dbName = mustGetEnv("NODA_ROOT_DB_NAME")

	dbAdminConf := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		dbUser, dbPassword, dbHost, dbPort, dbName)
	db, err = sql.Open("postgres", dbAdminConf)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	executeScriptsOf("extensions", index.extensions)
	executeScriptsOf("domains", index.domains)
	executeScriptsOf("types", index.types)
	executeScriptsOf("tables", index.tables)
	executeScriptsOf("views", index.views)
	executeScriptsOf("indexes", index.indexes)
	executeScriptsOf("routines", index.routines)
	executeScriptsOf("seeds", index.seeds)
}

func executeScriptsOf(directory string, filesWithPrecedence files) {
	length := len(filesWithPrecedence)
	executionPrecedenceMatters := length > 0
	var alreadyExecutedFiles files = nil
	if executionPrecedenceMatters {
		alreadyExecutedFiles = make(files, length)
		for _, shallowFileName := range filesWithPrecedence {
			ext := filepath.Ext(shallowFileName)
			if ext != "" {
				log.Fatalf("please do not provide any extension (%s) for %q",
					ext, shallowFileName)
			}
			absoluteScriptPath := path.Join(databasePath, directory, shallowFileName+".sql")
			alreadyExecutedFiles = append(alreadyExecutedFiles, absoluteScriptPath)
			tryExecuteScript(absoluteScriptPath)
		}
	}
	TraverseDirectory(path.Join(databasePath, directory), alreadyExecutedFiles)
}

func TraverseDirectory(directory string, alreadyExecutedFiles files) {
	files, err := os.ReadDir(directory)
	if err != nil {
		log.Fatalf("could not open directory %q: %v",
			directory, err)
	}
	for _, file := range files {
		absoluteFilePath := path.Join(directory, file.Name())
		if file.IsDir() {
			TraverseDirectory(absoluteFilePath, alreadyExecutedFiles)
		} else {
			if filepath.Ext(absoluteFilePath) != ".sql" {
				continue
			}
			if !slices.Contains[[]string](alreadyExecutedFiles, absoluteFilePath) {
				tryExecuteScript(absoluteFilePath)
			}
		}
	}
}

func logNextFileToExecute(filename string) {
	fmt.Printf("Attempting to load and execute file: \033[1;32m%s\033[0m ...\n", filename)
}

func readQueryFromFile(script string) string {
	buf, err := os.ReadFile(script)
	if err != nil {
		log.Fatal(err)
	}
	return string(buf)
}

func tryExecuteScript(filename string) {
	logNextFileToExecute(filepath.Base(filename))
	query := readQueryFromFile(filename)
	if _, err := db.Exec(query); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			log.Fatal(failure.PQErrorToString(pqErr))
		}
	}

	lines := strings.Split(query, "\n")
	nonEmptyLines := make([]string, 0, len(lines))
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			nonEmptyLines = append(nonEmptyLines, fmt.Sprintf("  \033[0;32m│\033[0m    %s", line))
		}
	}
	if verbose {
		cleanedQuery := strings.Join(nonEmptyLines, "\n")
		fmt.Printf("\033[1;32m%s\033[0m\n", cleanedQuery)
	}
}

func findRootDirectory() string {
	currentProcessPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	for dir := currentProcessPath; dir != "/"; dir = filepath.Dir(dir) {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
	}
	log.Fatal("root directory not found (missed go.mod?)")
	return ""
}
