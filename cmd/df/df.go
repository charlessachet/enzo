package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"syscall"
)

const (
	pathfstab  string = "/etc/fstab"
	pathmounts string = "/proc/mounts"
)

type filesystem struct {
	Name        string
	Total       int64
	Used        int64
	UsedPercent int
	Available   int64
	MountedOn   string
}

func load(path string) ([]string, error) {
	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}
	defer file.Close()

	var parsed []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {

		line := scanner.Text()
		if isCommentedLine(line) == false {
			parsed = append(parsed, line)
		}

	}

	if err = scanner.Err(); err != nil {
		return nil, err
	}

	return parsed, nil
}

func isCommentedLine(line string) bool {
	if line[:1] == "#" {
		return true
	}
	return false
}

func findMountedFs(fstab []string, mounts []string) []*filesystem {
	filesystems := []*filesystem{}
	for _, fs := range fstab {

		toFind := parseFstabMountPointAndType(fs)

		for _, m := range mounts {

			if parseMountPointAndType(m) == toFind {
				d := new(filesystem)
				d.Name = parseFsName(m)
				d.MountedOn = parseMountOn(m)
				filesystems = append(filesystems, d)
				break
			}

		}

	}

	return filesystems
}

func parseFstabMountPointAndType(line string) string {
	fields := strings.Fields(line)[1:3]
	return strings.Join(fields, " ")

}

func parseMountPointAndType(line string) string {
	fields := strings.Fields(line)[1:3]
	return strings.Join(fields, " ")
}

func parseFsName(line string) string {
	return strings.Fields(line)[0]
}

func parseMountOn(line string) string {
	return strings.Fields(line)[1]
}

func calcFsUsage(filesystems []*filesystem) ([]*filesystem, error) {
	scall := syscall.Statfs_t{}
	for _, d := range filesystems {

		err := syscall.Statfs(d.MountedOn, &scall)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("error[%s] doing syscall to filesystem for mounted point[%s]", err, d.MountedOn))
		}

		d.Total = calcTotal(scall)
		d.Available = calcFree(scall)
		d.Used = calcUsed(d.Total, d.Available)
		d.UsedPercent = calcUsedPercent(d.Total, d.Used)

	}

	return filesystems, nil
}

func calcTotal(scall syscall.Statfs_t) int64 {
	return (scall.Frsize * int64(scall.Blocks)) / 1024
}

func calcFree(scall syscall.Statfs_t) int64 {
	return (scall.Frsize * int64(scall.Bfree)) / 1024
}

func calcUsed(total int64, free int64) int64 {
	return total - free
}

func calcUsedPercent(total int64, used int64) int {
	return int((float64(used) * 100) / float64(total))
}

func main() {
	fstab, err := load(pathfstab)
	if err != nil {
		errors.New(fmt.Sprintf("error loading '%v': %v", pathmounts, err.Error()))
	}

	mounts, err := load(pathmounts)
	if err != nil {
		errors.New(fmt.Sprintf("error loading '%v': %v", pathmounts, err.Error()))
	}

	filesystems := findMountedFs(fstab, mounts)
	filesystems, err = calcFsUsage(filesystems)
	if err != nil {
		errors.New(fmt.Sprintf("error calculating paths usage: %v", err.Error()))
	}

	for _, p := range filesystems {
		fmt.Println("Filesystem", p)
	}
}
