package safety

// #include <sys/types.h>
// #include <unistd.h>
import "C"
import (
	"fmt"
	"github.com/fredlahde/kobana/config"
	"github.com/fredlahde/kobana/errors"
	osuser "os/user"
	"strconv"
)

func DropRootPriviliges(config *config.Config) error {
	op := errors.Op("safety.DropRootPriviliges")
	userName := config.User.User
	groupName := config.User.Group
	user, err := osuser.Lookup(userName)
	if err != nil {
		return errors.E(op, errors.IO, err, errors.C("failed to get information about given user"))
	}
	uid := user.Uid

	foundGroup, err := osuser.LookupGroup(groupName)
	if err != nil {
		return errors.E(op, errors.IO, err, errors.C("failed to get informaton about given group"))
	}

	userGroups, err := user.GroupIds()
	if err != nil {
		return errors.E(op, errors.IO, err, errors.C("failed to get informaton about user groups"))
	}

	found := false
	for _, ug := range userGroups {
		if ug == foundGroup.Gid {
			found = true
			break
		}
	}

	if !found {
		return errors.E(op, errors.Security, fmt.Errorf("user %s is not in group %s", userName, groupName))
	}

	gidNumber, err := strconv.ParseInt(foundGroup.Gid, 10, 32)
	if err != nil {
		return errors.E(op, errors.Security, fmt.Errorf("user %s is not in group %s", userName, groupName))
	}

	ret := C.setgid(C.uint(gidNumber))
	if ret != 0 {
		return errors.E(op, errors.Security, fmt.Errorf("could not set gid to %d", gidNumber))
	}

	uidNumber, err := strconv.ParseInt(uid, 10, 32)
	if err != nil {
		return errors.E(op, errors.Security, err, errors.C("uid ist not numeric"))
	}

	ret = C.setuid(C.uint(uidNumber))
	if ret != 0 {
		return errors.E(op, errors.Security, fmt.Errorf("could not set uid to %d", uidNumber))
	}
	return nil
}
