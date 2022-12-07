package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"
)

var inputFile = flag.String("f", "input", "Relative file path to use as input.")

type File struct {
	Name string
	Size int
	Dir  *Dir
}

func (f File) FullPath() string {
	result := f.Name

	d := f.Dir
	for d.parentDir != nil {

		result = filepath.Join(d.Name, d.Name)
		d = d.parentDir
	}
	return result

}

type Files []File
type Dir struct {
	dirs      map[string]*Dir
	parentDir *Dir
	Name      string
	files     Files
	size      int
}

func (d Dir) String() string {
	files := ""
	for _, f := range d.files {
		files += fmt.Sprintf("%s (%d)|", f.Name, f.Size)
	}
	dirs := ""
	for _, d := range d.dirs {
		dirs += fmt.Sprintf("%s", d)
	}
	return fmt.Sprintf("Name: %s\nFiles: %s\n Dirs: %s\n", d.Name, files, dirs)
}

func (d *Dir) Size() int {
	if d.size != 0 {
		return d.size
	}
	r := 0
	for _, f := range d.files {
		r += f.Size
	}
	for _, d := range d.dirs {
		r += d.Size()
	}
	d.size = r
	return r
}

func NewDir(Name string) *Dir {
	return &Dir{
		Name:  Name,
		files: make(Files, 0, 0),
		dirs:  make(map[string]*Dir, 0),
	}
}
func parseLines(lines []string) (*Dir, []*Dir) {

	index := 0
	var rootDir *Dir
	currentDir := new(Dir)
	rootDir = currentDir
	dirList := make([]*Dir, 0, 10)
	for index < len(lines) {
		l := lines[index]
		switch string(l[0]) {
		// Command
		case "$":
			// We assume that you can't cd without a previous ls
			// ls commands are ignored as we read the output afterwards
			if string(l[0:4]) == "$ cd" {
				var dirName string
				fmt.Sscanf(l, "$ cd %s", &dirName)
				switch dirName {
				case "..":
					currentDir = currentDir.parentDir
				case "/":
					*currentDir = *NewDir(dirName)
				default:
					currentDir = currentDir.dirs[dirName]
				}
			}
		// Directory
		case "d":
			var dirName string
			fmt.Sscanf(l, "dir %s", &dirName)
			currentDir.dirs[dirName] = NewDir(dirName)
			currentDir.dirs[dirName].parentDir = currentDir
			dirList = append(dirList, currentDir.dirs[dirName])
		// File
		default:
			var fileName string
			var fileSize int
			fmt.Sscanf(l, "%d %s", &fileSize, &fileName)
			currentDir.files = append(currentDir.files, File{
				Name: fileName,
				Size: fileSize,
				Dir:  currentDir,
			})
		}
		index++
	}

	return rootDir, dirList
}
func main() {
	flag.Parse()
	input, _ := ioutil.ReadFile(*inputFile)
	trimmedInput := strings.Split(strings.TrimSpace(string(input)), "\n")
	rootDir, dirList := parseLines(trimmedInput)
	sizeLimit := 100_000
	part1 := 0
	for _, d := range dirList {
		dirSize := d.Size()
		if dirSize < sizeLimit {
			part1 += dirSize
		}
	}
	fmt.Println("Part1", part1)
	totalSize := 70000000
	freeSpace := totalSize - rootDir.Size()
	missingSpace := 30000000 - freeSpace
	fmt.Println("need to free", missingSpace)
	sort.Slice(dirList, func(i, j int) bool {
		return dirList[i].Size() < dirList[j].Size()
	})
	for _, v := range dirList {
		if v.Size() >= missingSpace {
			fmt.Println("Part2", v.Size())
			break
		}
	}

}
