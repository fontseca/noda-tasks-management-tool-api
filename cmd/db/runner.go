package main

import (
	"database/sql"
	"fmt"
	"log"
	"noda/config"
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

type Files []string

type Index struct {
	Init       Files `yaml:"init"`
	Extensions Files `yaml:"extensions"`
	Domains    Files `yaml:"domains"`
	Types      Files `yaml:"types"`
	Tables     Files `yaml:"tables"`
	Views      Files `yaml:"views"`
	Indexes    Files `yaml:"indexes"`
	Routines   Files `yaml:"routines"`
	Seeds      Files `yaml:"seeds"`
}

func main() {
	buf, err := os.ReadFile(indexPath)
	if err != nil {
		log.Fatalf("could not read %q: %v", indexPath, err)
	}

	var index Index
	err = yaml.Unmarshal(buf, &index)
	if err != nil {
		log.Fatalf("error unmarshalling %q: %v", indexPath, err)
	}

	dbRootConf := config.GetDatabaseConfigWithValues("postgres", "", "", "postgres", "postgres")
	db, err = sql.Open("postgres", dbRootConf.Conn())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	dbRootConf.LogSuccess()

	ExecuteScriptsOf("init", &index.Init)

	if err = db.Close(); err != nil {
		log.Fatal(err)
	}

	dbAdminConf := config.GetDatabaseConfigWithValues("noda", "", "", "admin", "admin")
	db, err = sql.Open("postgres", dbAdminConf.Conn())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	dbAdminConf.LogSuccess()

	ExecuteScriptsOf("extensions", &index.Extensions)
	ExecuteScriptsOf("domains", &index.Domains)
	ExecuteScriptsOf("types", &index.Types)
	ExecuteScriptsOf("tables", &index.Tables)
	ExecuteScriptsOf("views", &index.Views)
	ExecuteScriptsOf("indexes", &index.Indexes)
	ExecuteScriptsOf("routines", &index.Routines)
	ExecuteScriptsOf("seeds", &index.Seeds)
}

func ExecuteScriptsOf(directory string, filesWithPrecedence *Files) {
	length := len(*filesWithPrecedence)
	executionPrecedenceMatters := length > 0
	var alreadyExecutedFiles Files = nil
	if executionPrecedenceMatters {
		alreadyExecutedFiles = make(Files, length)
		for _, shallowFileName := range *filesWithPrecedence {
			ext := filepath.Ext(shallowFileName)
			if ext != "" {
				log.Fatalf("please do not provide any extension (%s) for %q",
					ext, shallowFileName)
			}
			absoluteScriptPath := path.Join(databasePath, directory, shallowFileName+".sql")
			alreadyExecutedFiles = append(alreadyExecutedFiles, absoluteScriptPath)
			TryExecuteScript(absoluteScriptPath)
		}
	}
	TraverseDirectory(path.Join(databasePath, directory), &alreadyExecutedFiles)
}

func TraverseDirectory(directory string, alreadyExecutedFiles *Files) {
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
			if alreadyExecutedFiles != nil &&
				!slices.Contains[[]string](*alreadyExecutedFiles, absoluteFilePath) {
				TryExecuteScript(absoluteFilePath)
			}
		}
	}
}

func LogNextFileToExecute(filename string) {
	log.Printf("Attempting to load and execute file: \033[1;32m%s\033[0m ...\n", filename)
}

func ReadQueryFromFile(script string) string {
	buf, err := os.ReadFile(script)
	if err != nil {
		log.Fatal(err)
	}
	return string(buf)
}

func TryExecuteScript(filename string) {
	LogNextFileToExecute(filepath.Base(filename))
	query := ReadQueryFromFile(filename)
	if _, err := db.Exec(query); err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			log.Fatal(failure.PQErrorToString(pqErr))
		} else {
			log.Fatal(err)
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
