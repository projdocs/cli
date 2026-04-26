package dkr

import "regexp"

var (
	containerNameConflictRE = regexp.MustCompile(`^Error response from daemon: Conflict\. The container name "/[^"]+" is already in use by container "[a-f0-9]{64}"\. You have to remove \(or rename\) that container to be able to reuse that name\.$`)
	imageNotFoundRe         = regexp.MustCompile(`^Error response from daemon: No such image: `)
	networkNotFoundRe       = regexp.MustCompile(`^Error response from daemon: \{"message":"network [^"]+ not found"}$`)
	networkNotFoundRe2      = regexp.MustCompile(`^Error response from daemon: network [^"]+ not found$`)
)

func isNetworkNotFoundErr(err error) bool {
	if err == nil {
		return false
	}
	return networkNotFoundRe.MatchString(err.Error()) || networkNotFoundRe2.MatchString(err.Error())
}

func isNameConflictErr(err error) bool {
	if err == nil {
		return false
	}
	return containerNameConflictRE.MatchString(err.Error())
}

func isImageNotFoundErr(err error) bool {
	if err == nil {
		return false
	}
	return imageNotFoundRe.MatchString(err.Error())
}
