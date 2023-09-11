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
	"regexp"
	"strings"

	"github.com/lib/pq"
)

var (
	db      *sql.DB
	scripts string
)

type DBQueryMapper = uint32

const (
	DBAll DBQueryMapper = iota
	DBCreateExtensions
	DBDrop
	DBCreate
	DBInit
	DBRoles
	DBObjects
	DBDomains
	DBTypes
	DBEnums
	DBComposites
	DBTables
	DBViews
	DBIndexes
	DBFunctions
	DBProcedures
	DBSeeds
)

type FileType = string

const (
	File      string = "file"
	Directory string = "directory"
)

type dbEntry struct {
	longname string
	path     string
	ftype    FileType
}

type DBDirectoryMap map[DBQueryMapper]*dbEntry

var dbMap DBDirectoryMap

var verbose bool = false

func init() {
	currentProcessPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	root := findRootDirectory(currentProcessPath)
	scripts = path.Join(root, "database/scripts")
	init := path.Join(scripts, "000_init")
	roles := path.Join(scripts, "001_roles")
	objects := path.Join(scripts, "002_objects")
	seeds := path.Join(scripts, "003_seeds")
	dbMap = DBDirectoryMap{
		DBAll:              {ftype: Directory, longname: "all", path: scripts},
		DBInit:             {ftype: Directory, longname: "init", path: init},
		DBDrop:             {ftype: File, longname: "dropdb", path: path.Join(init, "000_drop_db.sql")},
		DBCreate:           {ftype: File, longname: "createdb", path: path.Join(init, "001_create_db.sql")},
		DBCreateExtensions: {ftype: File, longname: "ext", path: path.Join(init, "002_extensions.sql")},
		DBRoles:            {ftype: Directory, longname: "roles", path: roles},
		DBObjects:          {ftype: Directory, longname: "objects", path: objects},
		DBDomains:          {ftype: Directory, longname: "objects/domains", path: path.Join(objects, "000_domains")},
		DBTypes:            {ftype: Directory, longname: "objects/types", path: path.Join(objects, "001_types")},
		DBEnums:            {ftype: Directory, longname: "objects/types/enumerations", path: path.Join(objects, "001_types", "000_enumerations")},
		DBComposites:       {ftype: Directory, longname: "objects/types/composites", path: path.Join(objects, "001_types", "001_composites")},
		DBTables:           {ftype: Directory, longname: "objects/tables", path: path.Join(objects, "002_tables")},
		DBViews:            {ftype: Directory, longname: "objects/views", path: path.Join(objects, "003_views")},
		DBIndexes:          {ftype: Directory, longname: "objects/indexes", path: path.Join(objects, "004_indexes")},
		DBFunctions:        {ftype: Directory, longname: "objects/functions", path: path.Join(objects, "005_functions")},
		DBProcedures:       {ftype: Directory, longname: "objects/procedures", path: path.Join(objects, "006_procedures")},
		DBSeeds:            {ftype: Directory, longname: "seeds", path: seeds},
	}
}

func main() {
	var err error
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

	executeQueryFor(DBDrop)
	executeQueryFor(DBRoles)
	executeQueryFor(DBCreate)

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

	executeQueryFor(DBCreateExtensions)
	executeQueryFor(DBObjects)
	executeQueryFor(DBSeeds)
}

func traverseAndExecute(directory string) {
	files, err := os.ReadDir(directory)
	if err != nil {
		log.Fatal("error reading directory:", err)
	}

	for _, file := range files {
		filename := file.Name()

		if !isCorrectDatabaseFileName(filename) {
			continue
		}

		if file.IsDir() {
			traverseAndExecute(path.Join(directory, filename))
		} else {
			logNextFileToExecute(filename)
			absPath := path.Join(directory, filename)
			query := readQueryFromFile(absPath)
			executeQueryOrFatal(query)
		}
	}
}

func executeQueryFor(entiy DBQueryMapper) {
	entry := dbMap[entiy]
	switch entry.ftype {
	default:
		log.Fatal("file type is unrecognizable")
	case Directory:
		traverseAndExecute(entry.path)
	case File:
		logNextFileToExecute(path.Base(entry.path))
		query := readQueryFromFile(entry.path)
		executeQueryOrFatal(query)
	}
}

func logNextFileToExecute(filename string) {
	log.Printf("Attempting to read and execute PL/pgSQL file `\033[1;32m%s\033[0m'\n", filename)
}

func readQueryFromFile(script string) string {
	buf, err := os.ReadFile(script)
	if err != nil {
		log.Fatal(err)
	}
	return string(buf)
}

func executeQueryOrFatal(query string) {
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

func findRootDirectory(startDirectory string) string {
	for dir := startDirectory; dir != "/"; dir = filepath.Dir(dir) {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
	}
	log.Fatal("root directory not found (missed go.mod?)")
	return ""
}

func isCorrectDatabaseFileName(filename string) bool {
	var (
		dbFilePattern             = `^[0-9]{3}_[a-z_]+.sql$`
		dbDirectoryPattern        = `^[0-9]{3}_[a-z_]+$`
		fileRegex, directoryRegex *regexp.Regexp
		err                       error
	)
	if fileRegex, err = regexp.Compile(dbFilePattern); err != nil {
		log.Fatal(err)
	}
	if directoryRegex, err = regexp.Compile(dbDirectoryPattern); err != nil {
		log.Fatal(err)
	}
	if !fileRegex.MatchString(filename) && !directoryRegex.MatchString(filename) {
		return false
	}
	return true
}
