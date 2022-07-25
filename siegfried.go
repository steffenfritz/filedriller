package filedriller

import (
	//"log"
	"os"
	"os/exec"
	"strconv"
	"syscall"

	"github.com/richardlehane/siegfried"
)

func siegfriedIdent(s *siegfried.Siegfried, inFile string) (bool, string) {
	var oneFile string
	var resultBool bool

	f, err := os.Open(inFile)
	if err != nil {
		return resultBool, err.Error()
	}

	defer f.Close()

	fi, _ := f.Stat()
	if fi.Size() == 0 {
		return resultBool, "\"" + inFile + "\",,,,,,,,,"
	}

	// We read the file owner here. As operation systems handle the users differently
	// we have to distinguish it here.
	// most Linux/Unix os should work with the syscall. However, this returns only the UID
	// as we mostly work with mounted images and a lookup would fail.
	var win bool
	var UID int
	var out []byte
	//var GID int
	if stat, ok := fi.Sys().(*syscall.Stat_t); ok {
		UID = int(stat.Uid)
		//GID = int(stat.Gid)
	} else {
		// We probably fail above when not on *nix and try Windows here.
		// If we are not on Windows this will finally fail.
		win = true
		arg0 := "path"
		arg1 := "Win32_LogicalFileSecuritySetting where Path=\"C:\\windows\\winsxs\""
		arg2 := "ASSOC /RESULTROLE:Owner /ASSOCCLASS:Win32_LogicalFileOwner /RESULTCLASS:Win32_SID"
		out, err = exec.Command("wmic", arg0, arg1, arg2).Output()

		if err != nil {
			ret := inFile + " : " + err.Error()
			return resultBool, ret
		}
	}

	ids, err := s.Identify(f, "", "")
	if err != nil {
		ret := inFile + " : " + err.Error()
		return resultBool, ret
	}

	for _, id := range ids {
		values := id.Values()
		for _, value := range values {
			oneFile += "\"" + value + "\"" + ","
		}
		// Normalize the owner information
		var ownerInfo string
		if win {
			ownerInfo = string(out)
		} else {
			ownerInfo = strconv.Itoa(UID)
		}
		oneFile = "\"" + inFile + "\",\"" + strconv.Itoa(int(fi.Size())) + "\"," + ownerInfo + "\"," + oneFile[:len(oneFile)-1] // remove last comma
	}

	return true, oneFile

}
