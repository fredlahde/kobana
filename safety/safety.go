package safety

// #include <sys/types.h>
// #include <unistd.h>
import "C"
import (
	"fmt"
	"github.com/pkg/errors"
	osuser "os/user"
	"strconv"
)

func DropRootPriviliges(username, groupname string) error {
	user, err := osuser.Lookup(username)
	if err != nil {
		return errors.Wrap(err, "failed to get information about given user")
	}
	uid := user.Uid

	foundGroup, err := osuser.LookupGroup(groupname)
	if err != nil {
		return errors.Wrap(err, "failed to get informaton about given group")
	}

	userGroups, err := user.GroupIds()
	if err != nil {
		return errors.Wrap(err, "failed to get informaton about user groups")
	}

	found := false
	for _, ug := range userGroups {
		if ug == foundGroup.Gid {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("user %s is not in group %s", username, groupname)
	}

	gidNumber, err := strconv.ParseInt(foundGroup.Gid, 10, 32)
	if err != nil {
		return errors.Wrap(err, "gid ist not numeric")
	}

	ret := C.setgid(C.uint(gidNumber))
	if ret != 0 {
		return fmt.Errorf("Could not set gid to %d", gidNumber)
	}

	uidNumber, err := strconv.ParseInt(uid, 10, 32)
	if err != nil {
		return errors.Wrap(err, "uid ist not numeric")
	}

	ret = C.setuid(C.uint(uidNumber))
	if ret != 0 {
		return fmt.Errorf("Could not set uid to %d", uidNumber)
	}
	return nil
}
