package nfs

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/DblK/tinshop/config"
	"github.com/DblK/tinshop/repository"
	"github.com/vmware/go-nfs-client/nfs/util"
)

var gameFiles []repository.FileDesc

type nfsSource struct {
}

// New create a nfs source
func New() repository.Source {
	gameFiles = make([]repository.FileDesc, 0)
	return &nfsSource{}
}

func (src *nfsSource) Download(w http.ResponseWriter, r *http.Request, game, share string) {
	if config.GetConfig().DebugNfs() {
		util.DefaultLogger.SetDebug(true)
	}

	shareInfos := strings.Split(share, ":")
	if len(shareInfos) != 2 {
		log.Printf("Error parsing the nfs share configuration (%s)\n", share)
		return
	}

	// Cut the share string
	host := shareInfos[0]
	path := shareInfos[1]
	name := path[strings.LastIndex(path, "/")+1:]
	target := path[:strings.LastIndex(path, "/")]

	// Connection
	mount, v := nfsConnect(host, target)
	defer mount.Close()
	defer v.Close()

	// Open file
	rdr, err := v.Open(name)
	if err != nil {
		util.Errorf("read error: %v", err)
		return
	}
	// Stats file
	fsInfo, _, err := v.Lookup(name)
	if err != nil {
		log.Fatalf("lookup error: %s", err.Error()) //nolint:gocritic
	}

	byteRange := strings.Split(strings.Replace(strings.Join(r.Header["Range"], ""), "bytes=", "", -1), "-")
	start, _ := strconv.Atoi(byteRange[0])
	end, _ := strconv.Atoi(byteRange[1])

	if start > int(fsInfo.Size()) || end > int(fsInfo.Size()) {
		w.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
		return
	}

	// Check if partial content
	if end-start+1 == int(fsInfo.Size()) {
		// Full Content
		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Length", fmt.Sprint(fsInfo.Size()))
		_, _ = io.Copy(w, rdr)
	} else {
		// Partial Content
		rng := make([]byte, end-start+1)
		if start != 0 {
			_, _ = rdr.Seek(int64(start), 0)
		}
		_, err = rdr.Read(rng)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("Error while trying to read file from nfs", err)
			return
		}
		w.WriteHeader(http.StatusPartialContent)
		w.Header().Add("Content-Range", "bytes "+fmt.Sprint(start)+"-"+fmt.Sprint(end)+"/"+fmt.Sprint(fsInfo.Size()))
		w.Header().Add("Accept-Ranges", "bytes")
		w.Header().Add("Content-Length", fmt.Sprint(end-start+1))
		_, _ = w.Write(rng)
	}

	if err = mount.Unmount(); err != nil {
		log.Fatalf("unable to umount target: %v", err)
	}
}
func (src *nfsSource) Load(shares []string, unique bool) {
	for _, share := range shares {
		loadGamesNfs(share)
	}
}
func (src *nfsSource) Reset() {
	gameFiles = make([]repository.FileDesc, 0)
}

func (src *nfsSource) UnWatchAll() {
	// Do nothing for now until nfs watcher as been done
}

func (src *nfsSource) GetFiles() []repository.FileDesc {
	return gameFiles
}
