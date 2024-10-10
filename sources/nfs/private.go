package nfs

import (
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/ajmandourah/tinshop/nsp"
	"github.com/ajmandourah/tinshop/repository"
	"github.com/ajmandourah/tinshop/utils"
	"github.com/vmware/go-nfs-client/nfs"
	"github.com/vmware/go-nfs-client/nfs/rpc"
	"github.com/vmware/go-nfs-client/nfs/util"
)

func getHostTarget(share string) (string, string, error) {
	shareInfos := strings.Split(share, ":")

	if len(shareInfos) != 2 {
		return "", "", errors.New("Error parsing the nfs share configuration " + share)
	}
	return shareInfos[0], shareInfos[1], nil
}

func (src *nfsSource) loadGamesNfs(share string) {
	if src.config.DebugNfs() {
		util.DefaultLogger.SetDebug(true)
	}

	host, target, err := getHostTarget(share)
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("Loading games from nfs (host=%s target=%s)\n", host, target)

	// Connection
	mount, v := nfsConnect(host, target)
	defer mount.Close()
	defer v.Close()

	nfsGames := src.lookIntoNfsDirectory(v, share, ".")

	mount.Close()
	src.gameFiles = append(src.gameFiles, nfsGames...)

	// Add all files
	if len(nfsGames) > 0 {
		src.collection.AddNewGames(nfsGames)
	}
}

func nfsConnect(host, target string) (*nfs.Mount, *nfs.Target) {
	mount, err := nfs.DialMount(host)
	if err != nil {
		log.Fatalf("unable to dial MOUNT service: %v", err)
	}

	// Mount drive
	v, err := mount.Mount(target, rpc.AuthNull)
	if err != nil {
		log.Fatalf("unable to mount volume: %v", err)
	}

	return mount, v
}

func (src *nfsSource) lookIntoNfsDirectory(v *nfs.Target, share, path string) []repository.FileDesc {
	// Retrieve all directories
	log.Printf("Retrieving all files in directory ('%s')...\n", path)

	dirs, err := v.ReadDirPlus(path)
	if err != nil {
		_ = fmt.Errorf("readdir error: %s", err.Error())
		return nil
	}

	var newGameFiles []repository.FileDesc

	for _, dir := range dirs {
		// Handle recursive directories
		if dir.IsDir() && dir.FileName != "." && dir.FileName != ".." {
			subDirGameFiles := src.lookIntoNfsDirectory(v, share, computePath(path, dir))
			newGameFiles = append(newGameFiles, subDirGameFiles...)
			continue
		}

		// Handle only NSP and NSZ files
		extension := filepath.Ext(dir.FileName)
		if extension != ".nsp" && extension != ".nsz" {
			continue
		}

		nfsRootPath := computeFullPath(share, path)
		newFile := repository.FileDesc{Size: dir.Size(), Path: nfsRootPath + "/" + dir.FileName}
		names, _ := utils.ExtractGameID(dir.FileName)

		if names.ShortID() == "" {
			// Useful to rename you file according to readme
			log.Println("Ignoring file because parsing failed", dir.FileName)
			continue
		}
		newFile.GameID = names.ShortID()
		newFile.GameInfo = names.FullID()
		newFile.HostType = repository.NFSShare
		newFile.Extension = names.Extension()

		var valid = true
		var errTicket error
		if src.config.VerifyNSP() {
			valid, errTicket = src.nspCheck(newFile)
		}
		if valid || (errTicket != nil && errTicket.Error() == "TitleDBKey for game "+newFile.GameID+" is not found") {
			newGameFiles = append(newGameFiles, newFile)
		} else {
			log.Println(errTicket)
		}
	}

	return newGameFiles
}

func computeFullPath(share, path string) string {
	nfsRootPath := share
	if path != "." {
		nfsRootPath += path
	}
	return nfsRootPath
}

func computePath(path string, dir *nfs.EntryPlus) string {
	var newPath string
	if path == "." {
		newPath = "/" + dir.FileName
	} else {
		newPath = path + "/" + dir.FileName
	}
	return newPath
}

func (src *nfsSource) nspCheck(file repository.FileDesc) (bool, error) {
	key, err := src.collection.GetKey(file.GameID)
	if err != nil {
		if src.config.DebugTicket() && err.Error() == "TitleDBKey for game "+file.GameID+" is not found" {
			log.Println(err)
		}
		return false, err
	}

	host, target, err := getHostTarget(file.Path)
	if err != nil {
		return false, err
	}

	mount, v := nfsConnect(host, filepath.Dir(target))
	defer mount.Close()
	defer v.Close()

	f, err := v.Open(filepath.Base(target))
	if err != nil {
		return false, err
	}
	defer f.Close()

	valid, err := nsp.IsTicketValid(f, key, src.config.DebugTicket())
	if err != nil {
		return false, err
	}
	if !valid {
		return false, errors.New("The ticket in '" + file.Path + "' is not valid!")
	}

	return valid, err
}
